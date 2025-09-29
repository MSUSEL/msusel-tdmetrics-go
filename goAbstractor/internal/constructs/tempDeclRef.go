package constructs

import (
	"go/types"

	"github.com/Snow-Gremlin/goToolbox/collections"
)

// TempDeclRef is a temporary reference used while abstracting a project
// to handle when a reference to a declaration is used prior to the
// declaration being defined. All references should be removed during
// resolving the since one thing that is resolved is all references.
type TempDeclRef interface {
	Construct
	Nestable
	NestType
	IsTempDeclRef()

	PackagePath() string
	Name() string
	Receiver() string
	ImplicitTypes() []TypeDesc
	InstanceTypes() []TypeDesc

	ResolvedType() Construct
	Resolved() bool
	SetResolution(con Construct)
}

type TempDeclRefArgs struct {
	PackagePath   string
	Name          string
	Receiver      string
	ImplicitTypes []TypeDesc
	InstanceTypes []TypeDesc
	Nest          NestType

	// FuncType is an optional type used when this decl ref is for a nest.
	FuncType *types.Func
}

type TempDeclRefFactory interface {
	Factory
	NewTempDeclRef(args TempDeclRefArgs) TempDeclRef
	TempDeclRefs() collections.ReadonlySortedSet[TempDeclRef]
	ClearAllTempDeclRefs()
}
