package constructs

import (
	"github.com/Snow-Gremlin/goToolbox/collections"
)

// TempDeclRef is a temporary reference used while abstracting a project
// to handle when a reference to a declaration is used prior to the
// declaration being defined. All references should be removed during
// resolving the since one thing that is resolved is all references.
type TempDeclRef interface {
	Construct
	IsTypeMethodRef()

	PackagePath() string
	Name() string
	InstanceTypes() []TypeDesc

	ResolvedType() Construct
	Resolved() bool
	SetResolution(con Construct)
}

type TempDeclRefArgs struct {
	PackagePath   string
	Name          string
	InstanceTypes []TypeDesc
}

type TempDeclRefFactory interface {
	NewTempDeclRef(args TempDeclRefArgs) TempDeclRef
	TempDeclRefs() collections.ReadonlySortedSet[TempDeclRef]
	ClearAllTempDeclRefs()
}