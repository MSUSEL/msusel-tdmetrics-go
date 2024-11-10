package constructs

import (
	"go/types"

	"golang.org/x/tools/go/packages"

	"github.com/Snow-Gremlin/goToolbox/collections"
)

type StructDesc interface {
	TypeDesc
	IsStructDesc()

	// Synthetic indicates this is a structure created around
	// a non-struct data for an object, e.g. `type Cat int`
	// has a synthetic `struct { $data int }`.
	Synthetic() bool

	Fields() []Field
}

type StructDescArgs struct {
	RealType types.Type

	Fields []Field

	// Package is needed when the real type isn't given.
	// The package is used to help create the real type.
	Package *packages.Package
}

type StructDescFactory interface {
	NewStructDesc(args StructDescArgs) StructDesc
	StructDescs() collections.ReadonlySortedSet[StructDesc]
}
