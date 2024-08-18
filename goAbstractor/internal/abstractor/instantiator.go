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
	case kind.Instance:
	case kind.InterfaceDesc:
	case kind.Reference:
	case kind.Signature:
	case kind.StructDesc:
	case kind.TypeParam:
	case kind.Object:
	case kind.InterfaceDecl:
	default:
		panic(terror.New(`unexpected type description kind`).
			With(`kind`, td.Kind()))
	}
}

func (i *instantiator) Interface(it constructs.InterfaceDesc) constructs.InterfaceDesc {
	abstracts := make([]constructs.Abstract, len(it.Abstracts()))
	for j, a := range it.Abstracts() {
		abstracts[j] = i.Abstract(a)
	}

	exact := make([]constructs.TypeDesc, len(it.Exact()))
	for j, e := range it.Exact() {
		exact[j] = i.TypeDesc(e)
	}

	approx := []constructs.TypeDesc{}
	for j, a := range it.Exact() {
		approx[j] = i.TypeDesc(a)
	}

	return i.proj.NewInterfaceDesc(constructs.InterfaceDescArgs{
		Abstracts: abstracts,
		Exact:     exact,
		Approx:    approx,
		Package:   i.decl.Package().Source(),
	})
}

func (i *instantiator) Abstract(a constructs.Abstract) constructs.Abstract {
	// TODO: Implement
	return nil
}
