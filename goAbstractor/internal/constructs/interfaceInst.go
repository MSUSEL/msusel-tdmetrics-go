package constructs

import (
	"go/types"

	"github.com/Snow-Gremlin/goToolbox/collections"
)

// InterfaceInst represents an instantiation of generic interface
// has been resolved to a specific interface with specific type parameters,
// e.g. List[T any] might be resolved to List<int>.
//
// The instance parameter may be referencing a type parameter,
// e.g. List[T any] might be resolved to List<S int|string>, thus the instance
// is also generic. The type may be a non-type parameter on a generic type,
// e.g. List[List[T any]] where List[T any] is the instance type.
type InterfaceInst interface {
	TypeDesc
	TempReferenceContainer
	IsInterfaceInst()

	Generic() InterfaceDecl
	Resolved() InterfaceDesc
	ImplicitTypes() []TypeDesc
	InstanceTypes() []TypeDesc
}

type InterfaceInstArgs struct {
	RealType      types.Type
	Generic       InterfaceDecl
	Resolved      InterfaceDesc
	ImplicitTypes []TypeDesc
	InstanceTypes []TypeDesc
}

type InterfaceInstFactory interface {
	Factory
	NewInterfaceInst(args InterfaceInstArgs) InterfaceInst
	InterfaceInsts() collections.ReadonlySortedSet[InterfaceInst]
}
