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
	Nestable
	IsTypeReference()

	PackagePath() string
	Name() string
	ImplicitTypes() []TypeDesc
	InstanceTypes() []TypeDesc

	ResolvedType() TypeDesc
	Resolved() bool
	SetResolution(typ TypeDesc)
}

type TempReferenceArgs struct {
	RealType      types.Type
	PackagePath   string
	Name          string
	ImplicitTypes []TypeDesc
	InstanceTypes []TypeDesc
	Nest          NestType

	// Package is needed when the real type isn't given.
	// The package is used to help create the real type.
	Package *packages.Package
}

type TempReferenceFactory interface {
	NewTempReference(args TempReferenceArgs) TempReference
	TempReferences() collections.ReadonlySortedSet[TempReference]
	ClearAllTempReferences()
}
