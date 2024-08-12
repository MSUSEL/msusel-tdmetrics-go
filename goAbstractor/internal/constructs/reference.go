package constructs

import (
	"go/types"

	"github.com/Snow-Gremlin/goToolbox/collections"
)

type Reference interface {
	TypeDesc
	IsReference()

	PackagePath() string
	Name() string
	Resolved() bool
	SetType(typ TypeDecl)
}

type ReferenceArgs struct {
	RealType    *types.Named
	PackagePath string
	Name        string
}

type ReferenceFactory interface {
	NewReference(args ReferenceArgs) Reference
	References() collections.ReadonlySortedSet[Reference]
}
