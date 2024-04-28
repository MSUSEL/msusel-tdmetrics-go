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

func uniqueName(names collections.Set[string]) string {
	for offset := 1; offset < 10000; offset++ {
		name := fmt.Sprintf(`value%d`, offset)
		if !names.Contains(name) {
			names.Add(name)
			return name
		}
	}
	return `_`
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

func (ab *abstractor) convertArray(t *types.Array) *typeDesc.Array {
	return &typeDesc.Array{
		Length: int(t.Len()),
		Elem:   ab.convertType(t.Elem()),
	}
}

func (ab *abstractor) convertBasic(t *types.Basic) *typeDesc.Ref {
	return &typeDesc.Ref{
		Ref: t.Name(),
	}
}

func (ab *abstractor) convertChan(t *types.Chan) *typeDesc.Chan {
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
	it.SortMethods()

	if t.IsImplicit() {
		for i := range t.NumEmbeddeds() {
			emb := ab.convertType(t.EmbeddedType(i))
			fmt.Printf(">> %[1]s (%[1]T)\n", emb) // TODO: Finish
		}
	}
	return ab.registerInterface(it)
}

func (ab *abstractor) convertMap(t *types.Map) *typeDesc.Map {
	return &typeDesc.Map{
		Key:   ab.convertType(t.Key()),
		Value: ab.convertType(t.Elem()),
	}
}

func (ab *abstractor) convertNamed(t *types.Named) *typeDesc.Ref {
	return &typeDesc.Ref{
		Ref: t.String(),
	}
}

func (ab *abstractor) convertFunc(t *types.Func) *typeDesc.Func {
	return &typeDesc.Func{
		Name:      t.Name(),
		Signature: ab.convertSignature(t.Type().(*types.Signature)),
	}
}

func (ab *abstractor) convertPointer(t *types.Pointer) *typeDesc.Pointer {
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

func (ab *abstractor) convertSlice(t *types.Slice) *typeDesc.Slice {
	return &typeDesc.Slice{
		Elem: ab.convertType(t.Elem()),
	}
}

func (ab *abstractor) convertStruct(t *types.Struct) *typeDesc.Struct {
	return ab.registerStruct(&typeDesc.Struct{
		Fields: convertList(t.NumFields(), t.Field, ab.convertField),
	})
}

func (ab *abstractor) createReturn(returns []*typeDesc.Field) typeDesc.TypeDesc {
	switch len(returns) {
	case 0:
		return nil
	case 1:
		return returns[0].Type
	default:
		names := set.From(enumerator.Select(enumerator.Enumerate(returns...),
			func(f *typeDesc.Field) string { return f.Name }).NotZero())
		for _, f := range returns {
			f.Anonymous = false
			if len(f.Name) <= 0 || f.Name == `_` {
				f.Name = uniqueName(names)
			}
		}
		return ab.registerStruct(&typeDesc.Struct{
			Fields: returns,
		})
	}
}

func (ab *abstractor) convertParamTuple(t *types.Tuple) []*typeDesc.Param {
	return convertList(t.Len(), t.At, ab.convertParam)
}

func (ab *abstractor) convertFieldTuple(t *types.Tuple) []*typeDesc.Field {
	return convertList(t.Len(), t.At, ab.convertField)
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

func (ab *abstractor) convertTypeParam(t *types.TypeParam) *typeDesc.TypeParam {

	t2 := t.Obj().Type().Underlying()
	fmt.Printf("%+v\n", t2)

	// TODO: FIX
	return ab.registerTypeParam(&typeDesc.TypeParam{
		Index:      t.Index(),
		Constraint: ab.convertType(t.Constraint()),
		Type:       ab.convertType(t2),
	})
}

func (ab *abstractor) convertTypeParamList(t *types.TypeParamList) []*typeDesc.TypeParam {
	return convertList(t.Len(), t.At, ab.convertTypeParam)
}

func (ab *abstractor) convertParam(t *types.Var) *typeDesc.Param {
	return &typeDesc.Param{
		Name: t.Name(),
		Type: ab.convertType(t.Type()),
	}
}

func (ab *abstractor) convertField(t *types.Var) *typeDesc.Field {
	return &typeDesc.Field{
		Anonymous: t.Anonymous(),
		Name:      t.Name(),
		Type:      ab.convertType(t.Type()),
	}
}
