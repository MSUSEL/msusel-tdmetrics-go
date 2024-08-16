package abstractor

import (
	"go/types"

	"github.com/Snow-Gremlin/goToolbox/terrors/terror"
	"github.com/Snow-Gremlin/goToolbox/utils"

	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/baker"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs"
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
		panic(terror.New(`unhandled type`).
			WithType(`type`, t).
			With(`value`, t))
	}
}

func (ab *abstractor) convertArray(t *types.Array) constructs.TypeDesc {
	elem := ab.convertType(t.Elem())
	generic := ab.baker.BakeList()
	return ab.proj.NewInstance(constructs.InstanceArgs{
		RealType: t,
		Generic:  generic,
		//Resolved:  // TODO: Fill out
		InstanceTypes: []constructs.TypeDesc{elem},
	})
}

func (ab *abstractor) convertBasic(t *types.Basic) constructs.TypeDesc {
	switch t.Kind() {
	case types.Complex64:
		return ab.baker.BakeComplex64()
	case types.Complex128:
		return ab.baker.BakeComplex128()
	default:
		return ab.proj.NewBasic(constructs.BasicArgs{
			RealType: t,
		})
	}
}

func (ab *abstractor) convertChan(t *types.Chan) constructs.TypeDesc {
	elem := ab.convertType(t.Elem())
	generic := ab.baker.BakeChan()
	return ab.proj.NewInstance(constructs.InstanceArgs{
		RealType: t,
		Generic:  generic,
		//Resolved:  // TODO: Fill out
		InstanceTypes: []constructs.TypeDesc{elem},
	})
}

func (ab *abstractor) convertInterface(t *types.Interface) constructs.InterfaceDesc {
	t = t.Complete()

	abstracts := []constructs.Abstract{}
	for i := range t.NumMethods() {
		f := t.Method(i)
		sig := ab.convertSignature(f.Type().(*types.Signature))
		abstract := ab.proj.NewAbstract(constructs.AbstractArgs{
			Name:      f.Name(),
			Signature: sig,
		})
		abstracts = append(abstracts, abstract)
	}

	var exact, approx []constructs.TypeDesc
	for i := range t.NumEmbeddeds() {
		et := t.EmbeddedType(i)
		if union, ok := et.(*types.Union); ok {
			exact2, approx2 := ab.readUnionTerms(union)
			exact = append(exact, exact2...)
			approx = append(approx, approx2...)
		}
	}

	return ab.proj.NewInterfaceDesc(constructs.InterfaceDescArgs{
		RealType:  t,
		Exact:     exact,
		Approx:    approx,
		Abstracts: abstracts,
	})
}

func (ab *abstractor) convertMap(t *types.Map) constructs.TypeDesc {
	key := ab.convertType(t.Key())
	value := ab.convertType(t.Elem())
	generic := ab.baker.BakeMap()
	return ab.proj.NewInstance(constructs.InstanceArgs{
		RealType: t,
		Generic:  generic,
		//Resolved:  // TODO: Fill out
		InstanceTypes: []constructs.TypeDesc{key, value},
	})
}

func (ab *abstractor) convertNamed(t *types.Named) constructs.TypeDesc {
	pkgPath := ``
	if !utils.IsNil(t.Obj().Pkg()) {
		pkgPath = t.Obj().Pkg().Path()
	}
	name := t.Obj().Name()

	// Check for builtin types that need to be baked.
	if len(pkgPath) <= 0 {
		switch name {
		case `error`:
			return ab.baker.BakeError()
		case `comparable`:
			return ab.baker.BakeComparable()
		}
		pkgPath = baker.BuiltinName
	}

	// Get any type parameters.
	instanceTp := ab.convertInstanceTypes(t.TypeArgs())

	// Check if the reference can already be found.
	_, typ, found := ab.proj.FindType(pkgPath, name, false)
	if !found {
		// Otherwise, create a reference that will be filled later.
		return ab.proj.NewReference(constructs.ReferenceArgs{
			RealType:      t,
			PackagePath:   pkgPath,
			Name:          name,
			InstanceTypes: instanceTp,
		})
	}

	if ab.needsInstance(typ, instanceTp) {
		return ab.proj.NewInstance(constructs.InstanceArgs{
			RealType: t,
			Generic:  typ,
			//Resolved:  // TODO: Fill out?
			InstanceTypes: instanceTp,
		})
	}

	return typ
}

func (ab *abstractor) needsInstance(_ constructs.TypeDecl, tp []constructs.TypeDesc) bool {
	// TODO: When creating an instance, the instance they types need
	//       to be checked to be different from the initial generic types.
	//       Example: If `func Foo[T any]() { ... Func[T]() ... }`
	return len(tp) > 0
}

func (ab *abstractor) convertPointer(t *types.Pointer) constructs.TypeDesc {
	elem := ab.convertType(t.Elem())
	generic := ab.baker.BakePointer()
	return ab.proj.NewInstance(constructs.InstanceArgs{
		RealType: t,
		Generic:  generic,
		//Resolved:  // TODO: Fill out
		InstanceTypes: []constructs.TypeDesc{elem},
	})
}

func (ab *abstractor) convertSignature(t *types.Signature) constructs.Signature {
	// Don't output receiver or receiver type here.
	// Don't convert type parameters here.
	return ab.proj.NewSignature(constructs.SignatureArgs{
		RealType: t,
		Variadic: t.Variadic(),
		Params:   ab.convertArguments(t.Params()),
		Results:  ab.convertArguments(t.Results()),
	})
}

func (ab *abstractor) convertSlice(t *types.Slice) constructs.TypeDesc {
	elem := ab.convertType(t.Elem())
	generic := ab.baker.BakeList()
	return ab.proj.NewInstance(constructs.InstanceArgs{
		RealType: t,
		Generic:  generic,
		//Resolved:  // TODO: Fill out
		InstanceTypes: []constructs.TypeDesc{elem},
	})
}

func (ab *abstractor) convertStruct(t *types.Struct) constructs.StructDesc {
	fields := make([]constructs.Field, 0, t.NumFields())
	for i := range t.NumFields() {
		f := t.Field(i)
		if !blankName(f.Name()) {
			field := ab.proj.NewField(constructs.FieldArgs{
				Name:     f.Name(),
				Type:     ab.convertType(f.Type()),
				Embedded: f.Embedded(),
			})
			fields = append(fields, field)
		}
	}

	return ab.proj.NewStructDesc(constructs.StructDescArgs{
		RealType: t,
		Fields:   fields,
	})
}

func (ab *abstractor) convertArguments(t *types.Tuple) []constructs.Argument {
	count := t.Len()
	list := make([]constructs.Argument, count)
	for i := range count {
		t2 := t.At(i)
		list[i] = ab.proj.NewArgument(constructs.ArgumentArgs{
			Name: t2.Name(),
			Type: ab.convertType(t2.Type()),
		})
	}
	return list
}

func (ab *abstractor) convertUnion(t *types.Union) constructs.InterfaceDesc {
	exact, approx := ab.readUnionTerms(t)
	return ab.proj.NewInterfaceDesc(constructs.InterfaceDescArgs{
		Exact:  exact,
		Approx: approx,
	})
}

func (ab *abstractor) readUnionTerms(t *types.Union) (exact, approx []constructs.TypeDesc) {
	for i := range t.Len() {
		term := t.Term(i)
		it := ab.convertType(term.Type())
		if term.Tilde() {
			approx = append(approx, it)
		} else {
			exact = append(exact, it)
		}
	}
	return exact, approx
}

func (ab *abstractor) convertTypeParam(t *types.TypeParam) constructs.TypeParam {
	if tr, ok := ab.typeParamReplacer[t]; ok {
		t = tr
	}

	t2 := t.Obj().Type().Underlying()
	return ab.proj.NewTypeParam(constructs.TypeParamArgs{
		Name: t.Obj().Name(),
		Type: ab.convertType(t2),
	})
}

func (ab *abstractor) convertInstanceTypes(t *types.TypeList) []constructs.TypeDesc {
	list := make([]constructs.TypeDesc, t.Len())
	for i := range t.Len() {
		list[i] = ab.convertType(t.At(i))
	}
	return list
}

func blankName(name string) bool {
	return len(name) <= 0 || name == `_` || name == `.`
}
