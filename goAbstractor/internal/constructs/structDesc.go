package constructs

import (
	"go/types"

	"github.com/Snow-Gremlin/goToolbox/collections"
)

type StructDesc interface {
	TypeDesc
	IsStructDesc()
}

type StructDescArgs struct {
	RealType types.Type

	Fields []Field
}

type StructDescFactory interface {
	NewStructDesc(args StructDescArgs) StructDesc
	StructDescs() collections.ReadonlySortedSet[StructDesc]
}
