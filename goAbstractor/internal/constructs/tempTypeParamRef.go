package constructs

import (
	"go/types"

	"github.com/Snow-Gremlin/goToolbox/collections"
)

// TempTypeParamRef is a temporary reference used while abstracting a type
// to handle when a type parameter references itself.
type TempTypeParamRef interface {
	TypeDesc
	IsTypeTypeParamRef()

	ResolvedType() TypeDesc
	Context() string
	Name() string
	Resolved() bool
	SetResolution(typ TypeDesc)
}

type TempTypeParamRefArgs struct {
	RealType types.Type
	Context  string
	Name     string
}

type TempTypeParamRefFactory interface {
	Factory
	NewTempTypeParamRef(args TempTypeParamRefArgs) TempTypeParamRef
	TempTypeParamRefs() collections.ReadonlySortedSet[TempTypeParamRef]
	ClearAllTempTypeParamRefs()
}
