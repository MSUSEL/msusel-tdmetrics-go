package constructs

import (
	"go/types"

	"github.com/Snow-Gremlin/goToolbox/collections"

	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/locs"
)

// InterfaceDecl is a named interface typically explicitly defined at the given
// location in the source code. The underlying type description
// can be a class or interface with optional parameter types.
//
// If type parameters are given then the interface is generic.
// Instances with realized versions of the interface,
// are added for each used instance in the source code. If there
// are no instances then the generic interface isn't used.
type InterfaceDecl interface {
	TypeDecl
	IsInterface()

	Interface() InterfaceDesc
	IsNamed() bool
	IsGeneric() bool
	TypeParams() []TypeParam
	AddInstance(inst InterfaceInst) InterfaceInst
	Instances() collections.ReadonlySortedSet[InterfaceInst]
	FindInstance(instanceTypes []TypeDesc) (InterfaceInst, bool)
}

type InterfaceDeclArgs struct {
	RealType types.Type
	Package  Package
	Name     string
	Location locs.Loc

	TypeParams []TypeParam
	Interface  InterfaceDesc
}

type InterfaceDeclFactory interface {
	NewInterfaceDecl(args InterfaceDeclArgs) InterfaceDecl
	InterfaceDecls() collections.ReadonlySortedSet[InterfaceDecl]
}
