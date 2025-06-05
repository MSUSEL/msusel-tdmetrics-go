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
	IsTempDeclRef()

	PackagePath() string
	Name() string
	ImplicitTypes() []TypeDesc
	InstanceTypes() []TypeDesc
	Nest() NestType

	ResolvedType() Construct
	Resolved() bool
	SetResolution(con Construct)
}

type TempDeclRefArgs struct {
	PackagePath   string
	Name          string
	ImplicitTypes []TypeDesc
	InstanceTypes []TypeDesc
	Nest          NestType
}

type TempDeclRefFactory interface {
	NewTempDeclRef(args TempDeclRefArgs) TempDeclRef
	TempDeclRefs() collections.ReadonlySortedSet[TempDeclRef]
	ClearAllTempDeclRefs()
}
