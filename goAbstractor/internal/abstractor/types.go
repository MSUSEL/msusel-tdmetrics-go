package abstractor

import (
	"fmt"
	"go/types"

	"github.com/Snow-Gremlin/goToolbox/collections"
	"github.com/Snow-Gremlin/goToolbox/collections/enumerator"
	"github.com/Snow-Gremlin/goToolbox/collections/set"

	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs/typeDesc"
)

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

func (ab *abstractor) convertArray(t *types.Array) typeDesc.TypeDesc {
	elem := ab.convertType(t.Elem())
	return typeDesc.NewSolid(ab.bakeList(), elem)
}

func (ab *abstractor) convertBasic(t *types.Basic) typeDesc.TypeDesc {
	switch t.Kind() {
	case types.Complex64:
		return ab.bakeComplex64()
	case types.Complex128:
		return ab.bakeComplex128()
	default:
		return typeDesc.NewBasic(t.Name())
	}
}

func (ab *abstractor) convertChan(t *types.Chan) typeDesc.TypeDesc {
	elem := ab.convertType(t.Elem())
	return typeDesc.NewSolid(ab.bakeChan(), elem)
}

func (ab *abstractor) convertInterface(t *types.Interface) *typeDesc.Interface {
	t = t.Complete()

	it := typeDesc.NewInterface()

	for i := range t.NumMethods() {
		f := t.Method(i)
		sig := ab.convertSignature(f.Type().(*types.Signature))
		it.AddFunc(f.Name(), sig)
	}

	if t.IsImplicit() {
		for i := range t.NumEmbeddeds() {
			emb := ab.convertType(t.EmbeddedType(i))
			fmt.Printf(">> %[1]s (%[1]T)\n", emb) // TODO: Finish
		}
	}
	return ab.registerInterface(it)
}

func (ab *abstractor) convertMap(t *types.Map) typeDesc.TypeDesc {
	key := ab.convertType(t.Key())
	value := ab.convertType(t.Elem())
	return typeDesc.NewSolid(ab.bakeMap(), key, value)
}

func (ab *abstractor) convertNamed(t *types.Named) *typeDesc.Named {
	// TODO: need to handle named better.
	return typeDesc.NewNamed(t.String(), nil)
}

func (ab *abstractor) convertPointer(t *types.Pointer) typeDesc.TypeDesc {
	elem := ab.convertType(t.Elem())
	return typeDesc.NewSolid(ab.bakePointer(), elem)
}

func (ab *abstractor) convertSignature(t *types.Signature) *typeDesc.Signature {
	// Don't output receiver or receiver type here.
	sig := typeDesc.NewSignature()
	sig.Variadic = t.Variadic()

	sig.TypeParams = ab.convertTypeParamList(t.TypeParams())

	sig.Params = ab.convertTuple(t.Params())
	sig.Return = ab.createReturn(ab.convertTuple(t.Results()))

	return ab.registerSignature(sig)
}

func (ab *abstractor) convertSlice(t *types.Slice) typeDesc.TypeDesc {
	elem := ab.convertType(t.Elem())
	return typeDesc.NewSolid(ab.bakeList(), elem)
}

func (ab *abstractor) convertStruct(t *types.Struct) *typeDesc.Struct {
	ts := typeDesc.NewStruct()
	for i := range t.NumFields() {
		f := t.Field(i)
		field := typeDesc.NewNamed(f.Name(), ab.convertType(f.Type()))
		ts.Fields = append(ts.Fields, field)
		if f.Embedded() {
			ts.Embedded = append(ts.Embedded, field)
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
	list := make([]*typeDesc.Named, t.Len())
	for i := range t.Len() {
		list[i] = ab.convertName(t.At(i))
	}
	return list
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
	fmt.Printf("convertTerm: %+v\n", t)

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

func (ab *abstractor) convertTypeParam(t *types.TypeParam) *typeDesc.Named {

	t2 := t.Obj().Type().Underlying()
	fmt.Printf("convertTypeParam: %+v\n", t2)

	// TODO: FIX
	return typeDesc.NewNamed(
		t.Obj().Name(),
		ab.convertType(t2),
	)
	//	Index:      t.Index(),
	//	Constraint: ab.convertType(t.Constraint()),
	//	Type:       ab.convertType(t2),
}

func (ab *abstractor) convertTypeParamList(t *types.TypeParamList) []*typeDesc.Named {
	list := make([]*typeDesc.Named, t.Len())
	for i := range t.Len() {
		list[i] = ab.convertTypeParam(t.At(i))
	}
	return list
}
