package abstractor

import (
	"go/types"

	"github.com/Snow-Gremlin/goToolbox/terrors/terror"
	"github.com/Snow-Gremlin/goToolbox/utils"

	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs/kind"
)

type instantiator struct {
	proj          constructs.Project
	decl          constructs.Declaration
	instanceTypes []constructs.TypeDesc
	conversion    map[constructs.TypeParam]constructs.TypeDesc
}

func newInstantiator(proj constructs.Project, decl constructs.Declaration, instanceTypes []constructs.TypeDesc) (*instantiator, constructs.Instance, bool) {
	typeParams := decl.TypeParams()
	count := len(typeParams)
	if count != len(instanceTypes) {
		panic(terror.New(`the amount of type params must match the instance types`).
			With(`type params`, count).
			With(`instance types`, len(instanceTypes)))
	}

	// Check declaration is a generic type.
	if count <= 0 {
		return nil, nil, false
	}

	// Check if instance types match the declaration types.
	// For example if `func Foo[T any]() { ... Func[T]() ... }` is given such that
	// the call to `Foo` doesn't need an instance since it matches the generic.
	same := true
	for i, tp := range typeParams {
		if tp != instanceTypes[i] {
			same = false
			break
		}
	}
	if same {
		return nil, nil, false
	}

	// Check that a matching instance doesn't already exist.
	instance, found := decl.Instances().Enumerate().
		Where(func(in constructs.Instance) bool {
			for i, it := range in.InstanceTypes() {
				if it != instanceTypes[i] {
					return false
				}
			}
			return true
		}).First()
	if found {
		return nil, instance, true
	}

	// Create a new instantiator for the new instance.
	conversion := make(map[constructs.TypeParam]constructs.TypeDesc, count)
	for i, tp := range typeParams {
		conversion[tp] = instanceTypes[i]
	}
	return &instantiator{
		proj:          proj,
		decl:          decl,
		instanceTypes: instanceTypes,
		conversion:    conversion,
	}, nil, true
}

func (ab *abstractor) instantiateTypeDecl(realType types.Type, decl constructs.TypeDecl, instanceTypes ...constructs.TypeDesc) constructs.TypeDesc {
	if instance := ab.instantiateDeclaration(realType, decl, instanceTypes...); !utils.IsNil(instance) {
		return instance
	}
	return decl
}

func (ab *abstractor) instantiateDeclaration(realType types.Type, decl constructs.Declaration, instanceTypes ...constructs.TypeDesc) constructs.TypeDesc {
	i, existing, needsInstance := newInstantiator(ab.proj, decl, instanceTypes)
	if !needsInstance {
		return nil
	}
	if !utils.IsNil(existing) {
		return existing
	}
	return ab.proj.NewInstance(constructs.InstanceArgs{
		RealType:      realType,
		Generic:       decl,
		Resolved:      i.TypeDesc(decl.Type()),
		InstanceTypes: instanceTypes,
	})
}

func (i *instantiator) TypeDesc(td constructs.TypeDesc) constructs.TypeDesc {
	switch td.Kind() {
	case kind.Basic:
		return td
	case kind.Instance:
		return i.Instance(td.(constructs.Instance))
	case kind.InterfaceDesc:
		return i.InterfaceDesc(td.(constructs.InterfaceDesc))
	case kind.Reference:
		return i.Reference(td.(constructs.Reference))
	case kind.Signature:
		return i.Signature(td.(constructs.Signature))
	case kind.StructDesc:
		return i.StructDesc(td.(constructs.StructDesc))
	case kind.TypeParam:
		return i.TypeParam(td.(constructs.TypeParam))
	case kind.Object:
		return i.Object(td.(constructs.Object))
	case kind.InterfaceDecl:
		return i.InterfaceDecl(td.(constructs.InterfaceDecl))
	default:
		panic(terror.New(`unexpected type description kind`).
			With(`kind`, td.Kind()).
			With(`type desc`, td))
	}
}

func mapSlice[T any, S ~[]T](s S, handle func(T) T) S {
	result := make(S, len(s))
	for i, e := range s {
		result[i] = handle(e)
	}
	return result
}

func (i *instantiator) Abstract(a constructs.Abstract) constructs.Abstract {
	return i.proj.NewAbstract(constructs.AbstractArgs{
		Name:      a.Name(),
		Signature: i.Signature(a.Signature()),
	})
}

func (i *instantiator) Argument(a constructs.Argument) constructs.Argument {
	return i.proj.NewArgument(constructs.ArgumentArgs{
		Name: a.Name(),
		Type: i.TypeDesc(a.Type()),
	})
}

func (i *instantiator) Field(f constructs.Field) constructs.Field {
	return i.proj.NewField(constructs.FieldArgs{
		Name:     f.Name(),
		Type:     i.TypeDesc(f.Type()),
		Embedded: f.Embedded(),
	})
}

func (i *instantiator) Instance(in constructs.Instance) constructs.Instance {
	// TODO: Check if instance is a match to the current instance types.
	//       Otherwise, look up the declaration and create a new instantiator.
	return nil
}

func (i *instantiator) InterfaceDecl(s constructs.InterfaceDecl) constructs.InterfaceDecl {
	// TODO: Use the declaration to create a new instantiator.
	return nil
}

func (i *instantiator) InterfaceDesc(it constructs.InterfaceDesc) constructs.InterfaceDesc {
	return i.proj.NewInterfaceDesc(constructs.InterfaceDescArgs{
		Abstracts: mapSlice(it.Abstracts(), i.Abstract),
		Exact:     mapSlice(it.Exact(), i.TypeDesc),
		Approx:    mapSlice(it.Approx(), i.TypeDesc),
		Package:   i.decl.Package().Source(),
	})
}

func (i *instantiator) Object(obj constructs.Object) constructs.Object {
	// TODO: Use the declaration to create a new instantiator.
	return nil
}

func (i *instantiator) Reference(r constructs.Reference) constructs.TypeDesc {
	if r.Resolved() {
		// The reference will probably not be resolved, but just in case
		// it has been, just return the instantiated resolved type and
		// skip the instantiated reference.
		return i.TypeDesc(r.ResolvedType())
	}

	return i.proj.NewReference(constructs.ReferenceArgs{
		PackagePath:   r.PackagePath(),
		Name:          r.Name(),
		InstanceTypes: mapSlice(r.InstanceTypes(), i.TypeDesc),
		Package:       i.decl.Package().Source(),
	})
}

func (i *instantiator) Signature(s constructs.Signature) constructs.Signature {
	return i.proj.NewSignature(constructs.SignatureArgs{
		Variadic: s.Variadic(),
		Params:   mapSlice(s.Params(), i.Argument),
		Results:  mapSlice(s.Results(), i.Argument),
		Package:  i.decl.Package().Source(),
	})
}

func (i *instantiator) StructDesc(s constructs.StructDesc) constructs.StructDesc {
	return i.proj.NewStructDesc(constructs.StructDescArgs{
		Fields:  mapSlice(s.Fields(), i.Field),
		Package: i.decl.Package().Source(),
	})
}

func (i *instantiator) TypeParam(tp constructs.TypeParam) constructs.TypeDesc {
	if t, has := i.conversion[tp]; has {
		return t
	}
	return tp
}
