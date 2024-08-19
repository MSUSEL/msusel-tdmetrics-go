package constructs

import (
	"go/types"

	"golang.org/x/tools/go/packages"

	"github.com/Snow-Gremlin/goToolbox/collections"
)

type Reference interface {
	TypeDesc
	IsReference()

	PackagePath() string
	Name() string
	InstanceTypes() []TypeDesc

	ResolvedType() TypeDesc
	Resolved() bool
	SetType(typ TypeDesc)
}

type ReferenceArgs struct {
	RealType      *types.Named
	PackagePath   string
	Name          string
	InstanceTypes []TypeDesc

	// Package is needed when the real type isn't given.
	// The package is used to help create the real type.
	Package *packages.Package
}

type ReferenceFactory interface {
	NewReference(args ReferenceArgs) Reference
	References() collections.ReadonlySortedSet[Reference]
}
