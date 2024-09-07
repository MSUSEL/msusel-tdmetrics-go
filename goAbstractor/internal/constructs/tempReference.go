package constructs

import (
	"go/types"

	"golang.org/x/tools/go/packages"

	"github.com/Snow-Gremlin/goToolbox/collections"
)

// TempReference is a temporary reference used while abstracting a project
// to handle when a reference to a type description is used prior to the
// type description being defined. All references should be removed during
// resolving the since one thing that is resolved is all references.
type TempReference interface {
	TypeDesc
	IsTypeReference()

	PackagePath() string
	Name() string
	InstanceTypes() []TypeDesc

	ResolvedType() TypeDesc
	Resolved() bool
	SetResolution(typ TypeDesc)
}

type TempReferenceArgs struct {
	RealType      *types.Named
	PackagePath   string
	Name          string
	InstanceTypes []TypeDesc

	// Package is needed when the real type isn't given.
	// The package is used to help create the real type.
	Package *packages.Package
}

type TempReferenceFactory interface {
	NewTempReference(args TempReferenceArgs) TempReference
	TempReferences() collections.ReadonlySortedSet[TempReference]
	ClearAllTempReferences()
}
