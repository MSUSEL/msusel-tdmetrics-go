package constructs

import (
	"go/types"

	"github.com/Snow-Gremlin/goToolbox/collections"
)

// Instance represents an instantiation of generic type
// has been resolved to a specific type with specific type parameters,
// e.g. List[T any] might be resolved to List<int>.
//
// The instance parameter may be referencing a type parameter,
// e.g. List[T any] might be resolved to List<S int|string>, thus the instance
// is also generic. The type may be a non-type parameter on a generic type,
// e.g. List[List[T any]] where List[T any] is the instance type.
type Instance interface {
	TypeDesc
	Identifiable
	TempReferenceContainer
	IsInstance()

	Generic() Declaration
	Resolved() TypeDesc
	InstanceTypes() []TypeDesc
}

type InstanceArgs struct {
	RealType      types.Type
	Generic       Declaration
	Resolved      TypeDesc
	InstanceTypes []TypeDesc
}

type InstanceFactory interface {
	NewInstance(args InstanceArgs) Instance
	Instances() collections.ReadonlySortedSet[Instance]
}
