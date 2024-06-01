package abstractor

import (
	"fmt"
	"go/types"

	"github.com/Snow-Gremlin/goToolbox/collections/enumerator"
	"github.com/Snow-Gremlin/goToolbox/collections/set"

	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs/typeDesc"
)

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
	return typeDesc.NewSolid(t, ab.bakeList(), elem)
}

func (ab *abstractor) convertBasic(t *types.Basic) typeDesc.TypeDesc {
	switch t.Kind() {
	case types.Complex64:
		return ab.bakeComplex64()
	case types.Complex128:
		return ab.bakeComplex128()
	default:
		return typeDesc.NewBasic(ab.proj, t)
	}
}

func (ab *abstractor) convertChan(t *types.Chan) typeDesc.TypeDesc {
	elem := ab.convertType(t.Elem())
	return typeDesc.NewSolid(t, ab.bakeChan(), elem)
}

func (ab *abstractor) convertInterface(t *types.Interface) typeDesc.Interface {
	t = t.Complete()

	it := typeDesc.NewInterface(t)

	for i := range t.NumMethods() {
		f := t.Method(i)
		sig := ab.convertSignature(f.Type().(*types.Signature))
		it.AddFunc(f.Name(), sig)
	}

	if t.IsImplicit() {
		for i := range t.NumEmbeddeds() {
			et := t.EmbeddedType(i)

			switch et.(type) {
			case *types.Union:
				it.SetUnion(ab.convertType(et).(typeDesc.Union))
			}
		}
	}
	return ab.proj.RegisterInterface(it)
}

func (ab *abstractor) convertMap(t *types.Map) typeDesc.TypeDesc {
	key := ab.convertType(t.Key())
	value := ab.convertType(t.Elem())
	return typeDesc.NewSolid(t, ab.bakeMap(), key, value)
}

func (ab *abstractor) convertNamed(t *types.Named) typeDesc.Named {
	// TODO: need to handle named better.
	return typeDesc.NewNamed(t.String(), nil)
}

func (ab *abstractor) convertPointer(t *types.Pointer) typeDesc.TypeDesc {
	elem := ab.convertType(t.Elem())
	return typeDesc.NewSolid(t, ab.bakePointer(), elem)
}

func (ab *abstractor) convertSignature(t *types.Signature) typeDesc.Signature {
	// Don't output receiver or receiver type here.
	sig := typeDesc.NewSignature(t)
	sig.SetVariadic(t.Variadic())
	sig.AppendTypeParam(ab.convertTypeParamList(t.TypeParams())...)
	sig.AppendParam(ab.convertTuple(t.Params())...)
	sig.SetReturn(ab.createReturn(ab.convertTuple(t.Results())))

	return ab.proj.RegisterSignature(sig)
}

func (ab *abstractor) convertSlice(t *types.Slice) typeDesc.TypeDesc {
	elem := ab.convertType(t.Elem())
	return typeDesc.NewSolid(t, ab.bakeList(), elem)
}

func (ab *abstractor) convertStruct(t *types.Struct) typeDesc.Struct {
	ts := typeDesc.NewStruct(t)
	for i := range t.NumFields() {
		f := t.Field(i)
		field := typeDesc.NewNamed(f.Name(), ab.convertType(f.Type()))
		ts.AppendField(f.Embedded(), field)
	}
	return ab.proj.RegisterStruct(ts)
}

func (ab *abstractor) createReturn(returns []typeDesc.Named) typeDesc.TypeDesc {
	// TODO: Need to handle adding type parameters in struct
	//       or returning a solid type if single return has type parameters.
	switch len(returns) {
	case 0:
		return nil
	case 1:
		return returns[0].Type()
	default:
		names := set.From(enumerator.Select(enumerator.Enumerate(returns...),
			func(f typeDesc.Named) string { return f.Name() }).NotZero())
		for _, f := range returns {
			f.EnsureName(names)
		}
		st := typeDesc.NewStruct(nil)
		st.AppendField(false, returns...)
		return ab.proj.RegisterStruct(st)
	}
}

func (ab *abstractor) convertTuple(t *types.Tuple) []typeDesc.Named {
	list := make([]typeDesc.Named, t.Len())
	for i := range t.Len() {
		list[i] = ab.convertName(t.At(i))
	}
	return list
}

func (ab *abstractor) convertName(t *types.Var) typeDesc.Named {
	return typeDesc.NewNamed(t.Name(), ab.convertType(t.Type()))
}

func (ab *abstractor) convertUnion(t *types.Union) typeDesc.Union {
	union := typeDesc.NewUnion(t)
	for i := range t.Len() {
		term := t.Term(i)
		it := ab.convertType(term.Type())
		union.AddType(term.Tilde(), it)
	}
	return union
}

func (ab *abstractor) convertTypeParam(t *types.TypeParam) typeDesc.Named {
	t2 := t.Obj().Type().Underlying()
	return typeDesc.NewNamed(
		t.Obj().Name(),
		ab.convertType(t2),
	)
}

func (ab *abstractor) convertTypeParamList(t *types.TypeParamList) []typeDesc.Named {
	list := make([]typeDesc.Named, t.Len())
	for i := range t.Len() {
		list[i] = ab.convertTypeParam(t.At(i))
	}
	return list
}
