package constructs

import (
	"go/types"

	"github.com/Snow-Gremlin/goToolbox/collections"
)

// Basic is a base type (e.g. bool, int, string, float64).
//
// This does not contain complex types, those are treated as an interface.
type Basic interface {
	TypeDesc
	IsBasic()
}

type BasicArgs struct {
	// RealType is the basic type underlying this type.
	RealType *types.Basic
}

type BasicFactory interface {
	Factory
	NewBasic(args BasicArgs) Basic
	Basics() collections.ReadonlySortedSet[Basic]
}
