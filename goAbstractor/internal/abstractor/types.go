package abstractor

import (
	"fmt"
	"go/types"

	"github.com/Snow-Gremlin/goToolbox/collections"
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
	return typeDesc.NewSolid(ab.proj, t, ab.bakeList(), elem)
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
	return typeDesc.NewSolid(ab.proj, t, ab.bakeChan(), elem)
}

func (ab *abstractor) convertInterface(t *types.Interface) typeDesc.Interface {
	t = t.Complete()

	methods := map[string]typeDesc.TypeDesc{}
	for i := range t.NumMethods() {
		f := t.Method(i)
		sig := ab.convertSignature(f.Type().(*types.Signature))
		methods[f.Name()] = sig
	}

	var union typeDesc.Union
	if t.IsImplicit() {
		for i := range t.NumEmbeddeds() {
			et := t.EmbeddedType(i)
			switch et.(type) {
			case *types.Union:
				union = ab.convertType(et).(typeDesc.Union)
			}
		}
	}

	return typeDesc.NewInterface(ab.proj, typeDesc.InterfaceArgs{
		RealType: t,
		Union:    union,
		Methods:  methods,
	})
}

func (ab *abstractor) convertMap(t *types.Map) typeDesc.TypeDesc {
	key := ab.convertType(t.Key())
	value := ab.convertType(t.Elem())
	return typeDesc.NewSolid(ab.proj, t, ab.bakeMap(), key, value)
}

func (ab *abstractor) convertNamed(t *types.Named) typeDesc.TypeDefRef {
	// TODO: Update
	return typeDesc.NewTypeDefRef(ab.proj, t.String())
}

func (ab *abstractor) convertPointer(t *types.Pointer) typeDesc.TypeDesc {
	elem := ab.convertType(t.Elem())
	return typeDesc.NewSolid(ab.proj, t, ab.bakePointer(), elem)
}

func (ab *abstractor) convertSignature(t *types.Signature) typeDesc.Signature {
	// Don't output receiver or receiver type here.
	return typeDesc.NewSignature(ab.proj, typeDesc.SignatureArgs{
		RealType:   t,
		Variadic:   t.Variadic(),
		TypeParams: ab.convertTypeParamList(t.TypeParams()),
		Params:     ab.convertTuple(t.Params()),
		Return:     ab.createReturn(ab.convertTuple(t.Results())),
	})
}

func (ab *abstractor) convertSlice(t *types.Slice) typeDesc.TypeDesc {
	elem := ab.convertType(t.Elem())
	return typeDesc.NewSolid(ab.proj, t, ab.bakeList(), elem)
}

func (ab *abstractor) convertStruct(t *types.Struct) typeDesc.Struct {
	fields := []typeDesc.Named{}
	for i := range t.NumFields() {
		f := t.Field(i)
		field := typeDesc.NewNamed(ab.proj, f.Name(), ab.convertType(f.Type()))
		fields = append(fields, field)
		// Nothing needs to be done with f.Embedded() here.
	}
	return typeDesc.NewStruct(ab.proj, typeDesc.StructArgs{
		RealType: t,
		//TypeParams: , // TODO: Handle type params
		Fields: fields,
	})
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
		return typeDesc.NewStruct(ab.proj, typeDesc.StructArgs{
			Fields: returns,
		})
	}
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

func blankName(name string) bool {
	return len(name) <= 0 || name == `_` || name == `.`
}

func (ab *abstractor) convertTuple(t *types.Tuple) []typeDesc.Named {
	count := t.Len()
	names := make([]string, count)
	types := make([]typeDesc.TypeDesc, count)
	filledNames := set.New[string](count)
	for i := range count {
		t2 := t.At(i)
		name := t2.Name()
		names[i] = name
		types[i] = ab.convertType(t2.Type())
		if !blankName(name) {
			filledNames.Add(name)
		}
	}

	list := make([]typeDesc.Named, count)
	for i := range count {
		name := names[i]
		if blankName(name) {
			name = uniqueName(filledNames)
		}
		list[i] = typeDesc.NewNamed(ab.proj, name, types[i])
	}
	return list
}

func (ab *abstractor) convertUnion(t *types.Union) typeDesc.Union {
	exact := []typeDesc.TypeDesc{}
	approx := []typeDesc.TypeDesc{}
	for i := range t.Len() {
		term := t.Term(i)
		it := ab.convertType(term.Type())
		if term.Tilde() {
			approx = append(approx, it)
		} else {
			exact = append(exact, it)
		}
	}
	return typeDesc.NewUnion(ab.proj, typeDesc.UnionArgs{
		RealType: t,
		Exact:    exact,
		Approx:   approx,
	})
}

func (ab *abstractor) convertTypeParam(t *types.TypeParam) typeDesc.Named {
	t2 := t.Obj().Type().Underlying()
	return typeDesc.NewNamed(ab.proj, t.Obj().Name(), ab.convertType(t2))
}

func (ab *abstractor) convertTypeParamList(t *types.TypeParamList) []typeDesc.Named {
	list := make([]typeDesc.Named, t.Len())
	for i := range t.Len() {
		list[i] = ab.convertTypeParam(t.At(i))
	}
	return list
}
