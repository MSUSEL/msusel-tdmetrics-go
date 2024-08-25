package constructs

import (
	"go/types"

	"github.com/Snow-Gremlin/goToolbox/collections"
)

// Instance represents a generic type that has been resolved to a specific type
// with specific type parameters, e.g. List<T> might be resolved to List<int>.
// The type parameter resolution may be referencing another type parameter,
// e.g. a method signature inside a generic interface.
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
