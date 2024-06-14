package abstractor

import (
	"fmt"
	"go/types"

	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs"
	"github.com/Snow-Gremlin/goToolbox/collections"
	"github.com/Snow-Gremlin/goToolbox/collections/set"
	"github.com/Snow-Gremlin/goToolbox/utils"
)

func (ab *abstractor) convertType(t types.Type) constructs.TypeDesc {
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

func (ab *abstractor) convertArray(t *types.Array) constructs.TypeDesc {
	elem := ab.convertType(t.Elem())
	return constructs.NewSolid(ab.proj.Types(), t, ab.bakeList(), elem)
}

func (ab *abstractor) convertBasic(t *types.Basic) constructs.TypeDesc {
	switch t.Kind() {
	case types.Complex64:
		return ab.bakeComplex64()
	case types.Complex128:
		return ab.bakeComplex128()
	default:
		return constructs.NewBasic(ab.proj.Types(), t)
	}
}

func (ab *abstractor) convertChan(t *types.Chan) constructs.TypeDesc {
	elem := ab.convertType(t.Elem())
	return constructs.NewSolid(ab.proj.Types(), t, ab.bakeChan(), elem)
}

func (ab *abstractor) convertInterface(t *types.Interface) constructs.Interface {
	t = t.Complete()

	methods := map[string]constructs.TypeDesc{}
	for i := range t.NumMethods() {
		f := t.Method(i)
		sig := ab.convertSignature(f.Type().(*types.Signature))
		methods[f.Name()] = sig
	}

	var union constructs.Union
	if t.IsImplicit() {
		for i := range t.NumEmbeddeds() {
			et := t.EmbeddedType(i)
			switch et.(type) {
			case *types.Union:
				union = ab.convertType(et).(constructs.Union)
			}
		}
	}

	return constructs.NewInterface(ab.proj.Types(), constructs.InterfaceArgs{
		RealType: t,
		Union:    union,
		Methods:  methods,
	})
}

func (ab *abstractor) convertMap(t *types.Map) constructs.TypeDesc {
	key := ab.convertType(t.Key())
	value := ab.convertType(t.Elem())
	return constructs.NewSolid(ab.proj.Types(), t, ab.bakeMap(), key, value)
}

func (ab *abstractor) convertNamed(t *types.Named) constructs.TypeDefRef {
	pkgPath := ``
	if !utils.IsNil(t.Obj().Pkg()) {
		pkgPath = t.Obj().Pkg().Path()
	}
	name := t.Obj().Name()
	return constructs.NewTypeDefRef(ab.proj.Types(), t, pkgPath, name)
}

func (ab *abstractor) convertPointer(t *types.Pointer) constructs.TypeDesc {
	elem := ab.convertType(t.Elem())
	return constructs.NewSolid(ab.proj.Types(), t, ab.bakePointer(), elem)
}

func (ab *abstractor) convertSignature(t *types.Signature) constructs.Signature {
	// Don't output receiver or receiver type here.
	tp := ab.convertTypeParamList(t.TypeParams())
	return constructs.NewSignature(ab.proj.Types(), constructs.SignatureArgs{
		RealType:   t,
		Variadic:   t.Variadic(),
		TypeParams: tp,
		Params:     ab.convertTuple(t.Params()),
		Return:     ab.createReturn(ab.convertTuple(t.Results())),
	})
}

func (ab *abstractor) convertSlice(t *types.Slice) constructs.TypeDesc {
	elem := ab.convertType(t.Elem())
	return constructs.NewSolid(ab.proj.Types(), t, ab.bakeList(), elem)
}

func (ab *abstractor) convertStruct(t *types.Struct) constructs.Struct {
	fields := []constructs.Named{}
	for i := range t.NumFields() {
		f := t.Field(i)
		field := constructs.NewNamed(ab.proj.Types(), f.Name(), ab.convertType(f.Type()))
		fields = append(fields, field)
		// Nothing needs to be done with f.Embedded() here.
	}

	return constructs.NewStruct(ab.proj.Types(), constructs.StructArgs{
		RealType: t,
		Fields:   fields,
	})
}

func (ab *abstractor) createReturn(returns []constructs.Named) constructs.TypeDesc {
	switch len(returns) {
	case 0:
		return nil
	case 1:
		return returns[0].Type()
	default:
		return constructs.NewStruct(ab.proj.Types(), constructs.StructArgs{
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

func (ab *abstractor) convertTuple(t *types.Tuple) []constructs.Named {
	count := t.Len()
	names := make([]string, count)
	types := make([]constructs.TypeDesc, count)
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

	list := make([]constructs.Named, count)
	for i := range count {
		name := names[i]
		if blankName(name) {
			name = uniqueName(filledNames)
		}
		list[i] = constructs.NewNamed(ab.proj.Types(), name, types[i])
	}
	return list
}

func (ab *abstractor) convertUnion(t *types.Union) constructs.Union {
	exact := []constructs.TypeDesc{}
	approx := []constructs.TypeDesc{}
	for i := range t.Len() {
		term := t.Term(i)
		it := ab.convertType(term.Type())
		if term.Tilde() {
			approx = append(approx, it)
		} else {
			exact = append(exact, it)
		}
	}
	return constructs.NewUnion(ab.proj.Types(), constructs.UnionArgs{
		RealType: t,
		Exact:    exact,
		Approx:   approx,
	})
}

func (ab *abstractor) convertTypeParam(t *types.TypeParam) constructs.Named {
	if tr, ok := ab.typeParamReplacer[t]; ok {
		t = tr
	}

	t2 := t.Obj().Type().Underlying()
	return constructs.NewNamed(ab.proj.Types(), t.Obj().Name(), ab.convertType(t2))
}

func (ab *abstractor) convertTypeParamList(t *types.TypeParamList) []constructs.Named {
	list := make([]constructs.Named, t.Len())
	for i := range t.Len() {
		list[i] = ab.convertTypeParam(t.At(i))
	}
	return list
}
