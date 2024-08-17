package abstractor

import (
	"go/types"

	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs/kind"
	"github.com/Snow-Gremlin/goToolbox/terrors/terror"
	"github.com/Snow-Gremlin/goToolbox/utils"
)

type instantiator struct {
	proj          constructs.Project
	decl          constructs.Declaration
	instanceTypes []constructs.TypeDesc
	conversion    map[constructs.TypeParam]constructs.TypeDesc
}

func newInstantiator(proj constructs.Project, decl constructs.Declaration, instanceTypes []constructs.TypeDesc) *instantiator {
	typeParams := decl.TypeParams()
	conversion := make(map[constructs.TypeParam]constructs.TypeDesc, len(typeParams))
	for i, tp := range typeParams {
		conversion[tp] = instanceTypes[i]
	}
	return &instantiator{
		decl:          decl,
		instanceTypes: instanceTypes,
		conversion:    conversion,
	}
}

func (i *instantiator) needsInstance() bool {
	// TODO: When creating an instance, the instance they types need
	//       to be checked to be different from the initial generic types.
	//       Example: If `func Foo[T any]() { ... Func[T]() ... }`
	return len(i.instanceTypes) > 0
}

func (ab *abstractor) instantiateTypeDecl(realType types.Type, decl constructs.TypeDecl, instanceTypes ...constructs.TypeDesc) constructs.TypeDesc {
	if instance := ab.instantiateDeclaration(realType, decl, instanceTypes...); !utils.IsNil(instance) {
		return instance
	}
	return decl
}

func (ab *abstractor) instantiateDeclaration(realType types.Type, decl constructs.Declaration, instanceTypes ...constructs.TypeDesc) constructs.TypeDesc {
	i := newInstantiator(ab.proj, decl, instanceTypes)
	if !i.needsInstance() {
		return nil
	}

	var resolved constructs.TypeDesc
	switch decl.Kind() {
	case kind.InterfaceDecl:
		//resolved = i.Interface(decl.(constructs.InterfaceDecl).Interface())
	case kind.Object:
		//resolved = i.Interface(decl.(constructs.InterfaceDecl).Interface())
	case kind.Method:
	default:
		panic(terror.New(`unexpected declaration when creating instances`).
			With(`kind`, decl.Kind()).
			With(`decl`, decl))
	}

	return ab.proj.NewInstance(constructs.InstanceArgs{
		RealType:      realType,
		Generic:       decl,
		Resolved:      resolved,
		InstanceTypes: instanceTypes,
	})
}

func (i *instantiator) TypeDesc(td constructs.TypeDesc) constructs.TypeDesc {
	// TODO: Implement
	return nil
}

/*
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
		Package:   i.curPkg.Source(),
	})
}

func (i *instantiator) Abstract(a constructs.Abstract) constructs.Abstract {
	// TODO: Implement
	return nil
}
*/
