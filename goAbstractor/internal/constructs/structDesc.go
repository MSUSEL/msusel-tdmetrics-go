package constructs

import (
	"go/types"

	"golang.org/x/tools/go/packages"

	"github.com/Snow-Gremlin/goToolbox/collections"
)

type StructDesc interface {
	TypeDesc
	IsStructDesc()

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
