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
	InstanceTypes() []TypeDesc

	Resolved() bool
	SetType(typ TypeDesc)
}

type ReferenceArgs struct {
	RealType      *types.Named
	PackagePath   string
	Name          string
	InstanceTypes []TypeDesc
}

type ReferenceFactory interface {
	NewReference(args ReferenceArgs) Reference
	References() collections.ReadonlySortedSet[Reference]
}
