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
	IsInstance()
}

type InstanceArgs struct {
	RealType   *types.Signature
	Resolved   TypeDesc
	TypeParams []TypeDesc
}

type InstanceFactory interface {
	NewInstance(args InstanceArgs) Instance
	Instances() collections.ReadonlySortedSet[Instance]
}
