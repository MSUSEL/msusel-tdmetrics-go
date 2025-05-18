package constructs

import (
	"go/types"

	"github.com/Snow-Gremlin/goToolbox/collections"
)

// ObjectInst represents an instantiation of generic object
// has been resolved to a specific object with specific type parameters,
// e.g. List[T any] might be resolved to List<int>.
//
// The instance parameter may be referencing a type parameter,
// e.g. List[T any] might be resolved to List<S int|string>, thus the instance
// is also generic. The type may be a non-type parameter on a generic type,
// e.g. List[List[T any]] where List[T any] is the instance type.
type ObjectInst interface {
	TypeDesc
	TempReferenceContainer
	IsObjectInst()

	Generic() Object
	ResolvedData() StructDesc
	ResolvedInterface() InterfaceDesc
	ImplicitTypes() []TypeDesc
	InstanceTypes() []TypeDesc

	Methods() collections.ReadonlySortedSet[MethodInst]
	AddMethod(method MethodInst) MethodInst
	SetResolvedInterface(it InterfaceDesc)
}

type ObjectInstArgs struct {
	RealType      types.Type
	Generic       Object
	ResolvedData  StructDesc
	ImplicitTypes []TypeDesc
	InstanceTypes []TypeDesc
}

type ObjectInstFactory interface {
	NewObjectInst(args ObjectInstArgs) ObjectInst
	ObjectInsts() collections.ReadonlySortedSet[ObjectInst]
}
