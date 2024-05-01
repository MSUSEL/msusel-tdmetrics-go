package abstractor

import (
	"fmt"
	"go/types"
	"slices"

	"github.com/Snow-Gremlin/goToolbox/collections"
	"github.com/Snow-Gremlin/goToolbox/collections/enumerator"
	"github.com/Snow-Gremlin/goToolbox/collections/set"

	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs/typeDesc"
)

func convertList[T, U any](n int, getter func(i int) T, convert func(value T) *U) []*U {
	list := make([]*U, 0, n)
	for i := range n {
		if p := convert(getter(i)); p != nil {
			list = append(list, p)
		}
	}
	if len(list) <= 0 {
		return nil
	}
	return slices.Clip(list)
}

// uniqueName returns a unique name that isn't in the set.
// The new unique name will be added to the set.
// This is for naming anonymous fields and unnamed return values.
func uniqueName(names collections.Set[string]) string {
	const (
		attempts = 10_000
		pattern  = `$value%d`
	)
	for offset := 1; offset < attempts; offset++ {
		name := fmt.Sprintf(pattern, offset)
		if !names.Contains(name) {
			names.Add(name)
			return name
		}
	}
	panic(fmt.Errorf(`unable to find unique name in %d attempts`, attempts))
}

func (ab *abstractor) convertType(t types.Type) typeDesc.TypeDesc {
	switch t2 := t.(type) {
	case *types.Array:
		return ab.convertArray(t2)
	case *types.Basic:
		return ab.convertBasic(t2)
	case *types.Chan:
		return ab.convertChan(t2)
	case *types.Interface:
		return ab.convertInterface(t2)
	case *types.Map:
		return ab.convertMap(t2)
	case *types.Named:
		return ab.convertNamed(t2)
	case *types.Pointer:
		return ab.convertPointer(t2)
	case *types.Signature:
		return ab.convertSignature(t2)
	case *types.Slice:
		return ab.convertSlice(t2)
	case *types.Struct:
		return ab.convertStruct(t2)
	case *types.TypeParam:
		return ab.convertTypeParam(t2)
	case *types.Union:
		return ab.convertUnion(t2)
	default:
		panic(fmt.Errorf(`unhandled type, %[1]T: %[1]s`, t))
	}
}

func (ab *abstractor) convertArray(t *types.Array) *typeDesc.Interface {
	elem := ab.convertType(t.Elem())
	return &typeDesc.Struct{
		Length: int(t.Len()),
		Elem:   elem,
	}
}

func (ab *abstractor) convertBasic(t *types.Basic) *typeDesc.Basic {
	if t.Info() == types.IsComplex {
		panic(fmt.Errorf(`complex numbers currently isn't supported: %[1]v (%[1]T)`, t))
	}
	return &typeDesc.Basic{
		Name: t.Name(),
	}
}

func (ab *abstractor) convertChan(t *types.Chan) *typeDesc.Interface {
	return &typeDesc.Chan{
		Elem: ab.convertType(t.Elem()),
	}
}

func (ab *abstractor) convertInterface(t *types.Interface) *typeDesc.Interface {
	t = t.Complete()
	methods := convertList(t.NumMethods(), t.Method, ab.convertFunc)
	it := &typeDesc.Interface{
		Methods: methods,
	}

	if t.IsImplicit() {
		for i := range t.NumEmbeddeds() {
			emb := ab.convertType(t.EmbeddedType(i))
			fmt.Printf(">> %[1]s (%[1]T)\n", emb) // TODO: Finish
		}
	}
	return ab.registerInterface(it)
}

func (ab *abstractor) convertMap(t *types.Map) *typeDesc.Interface {
	return &typeDesc.Map{
		Key:   ab.convertType(t.Key()),
		Value: ab.convertType(t.Elem()),
	}
}

func (ab *abstractor) convertNamed(t *types.Named) *typeDesc.Named {
	return &typeDesc.Named{
		Name: t.String(),
	}
}

func (ab *abstractor) convertFunc(t *types.Func) *typeDesc.Func {
	return &typeDesc.Func{
		Name:      t.Name(),
		Signature: ab.convertSignature(t.Type().(*types.Signature)),
	}
}

func (ab *abstractor) convertPointer(t *types.Pointer) *typeDesc.Interface {
	return &typeDesc.Pointer{
		Elem: ab.convertType(t.Elem()),
	}
}

func (ab *abstractor) convertSignature(t *types.Signature) *typeDesc.Signature {
	// Don't output receiver or receiver type here.
	return ab.registerSignature(&typeDesc.Signature{
		Variadic:   t.Variadic(),
		Params:     ab.convertParamTuple(t.Params()),
		Return:     ab.createReturn(ab.convertFieldTuple(t.Results())),
		TypeParams: ab.convertTypeParamList(t.TypeParams()),
	})
}

func (ab *abstractor) convertSlice(t *types.Slice) *typeDesc.Interface {
	return &typeDesc.Slice{
		Elem: ab.convertType(t.Elem()),
	}
}

func (ab *abstractor) convertStruct(t *types.Struct) *typeDesc.Struct {
	ts := typeDesc.NewStruct()
	for i := range t.NumFields() {
		f := t.Field(i)
		field := typeDesc.NewNamed(f.Name(), ab.convertType(f.Type()))
		ts.Fields = append(ts.Fields, field)
		if f.Anonymous() {
			ts.Anonymous = append(ts.Anonymous, field)
		}
	}
	return ab.registerStruct(ts)
}

func (ab *abstractor) createReturn(returns []*typeDesc.Named) typeDesc.TypeDesc {
	// TODO: Need to handle adding type parameters in struct
	//       or returning a solid type if single return has type parameters.
	switch len(returns) {
	case 0:
		return nil
	case 1:
		return returns[0].Type
	default:
		names := set.From(enumerator.Select(enumerator.Enumerate(returns...),
			func(f *typeDesc.Named) string { return f.Name }).NotZero())
		for _, f := range returns {
			if len(f.Name) <= 0 || f.Name == `_` {
				f.Name = uniqueName(names)
			}
		}
		return ab.registerStruct(&typeDesc.Struct{
			Fields: returns,
		})
	}
}

func (ab *abstractor) convertTuple(t *types.Tuple) []*typeDesc.Named {
	return convertList(t.Len(), t.At, ab.convertName)
}

func (ab *abstractor) convertName(t *types.Var) *typeDesc.Named {
	return &typeDesc.Named{
		Name: t.Name(),
		Type: ab.convertType(t.Type()),
	}
}

func (ab *abstractor) convertTerm(t *types.Term) *typeDesc.Interface {
	// TODO: add `getData() T` for t.Tilde()
	//t2 := ab.convertType(t.Type())

	// TODO: FINISH

	return nil
}

func (ab *abstractor) convertUnion(t *types.Union) *typeDesc.Interface {
	union := &typeDesc.Interface{}
	for i := range t.Len() {
		it := ab.convertTerm(t.Term(i))

		// TODO: FINISH
		panic(fmt.Errorf(`union not implemented: %[1]v (%[1]T)`, it))

	}
	return union
}

func (ab *abstractor) convertTypeParam(t *types.TypeParam) *typeDesc.Interface {

	t2 := t.Obj().Type().Underlying()
	fmt.Printf("%+v\n", t2)

	// TODO: FIX
	return ab.registerTypeParam(&typeDesc.TypeParam{
		Index:      t.Index(),
		Constraint: ab.convertType(t.Constraint()),
		Type:       ab.convertType(t2),
	})
}

func (ab *abstractor) convertTypeParamList(t *types.TypeParamList) []*typeDesc.Interface {
	return convertList(t.Len(), t.At, ab.convertTypeParam)
}

func (ab *abstractor) convertField(t *types.Var) *typeDesc.Field {
	return &typeDesc.Field{
		Anonymous: t.Anonymous(),
		Name:      t.Name(),
		Type:      ab.convertType(t.Type()),
	}
}
