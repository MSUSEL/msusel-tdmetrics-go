package constructs

import (
	"go/types"

	"golang.org/x/tools/go/packages"

	"github.com/Snow-Gremlin/goToolbox/collections"
)

// Instance represents a generic type that has been resolved to a specific type
// with specific type parameters, e.g. List<T> might be resolved to List<int>.
// The type parameter resolution may be referencing another type parameter,
// e.g. a method signature inside a generic interface.
type Instance interface {
	TypeDesc
	IsInstance()

	InstanceTypes() []TypeDesc
}

type InstanceArgs struct {
	RealType      types.Type
	Generic       TypeDecl
	Resolved      TypeDesc
	InstanceTypes []TypeDesc

	// Package is needed when the real type isn't given.
	// The package is used to help create the real type.
	Package *packages.Package
}

type InstanceFactory interface {
	NewInstance(args InstanceArgs) Instance
	Instances() collections.ReadonlySortedSet[Instance]
}
