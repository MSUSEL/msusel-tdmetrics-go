package constructs

import (
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/locs"
	"github.com/Snow-Gremlin/goToolbox/collections"
)

type Value interface {
	Declaration
	IsValue()

	Type() TypeDesc
	Const() bool
}

type ValueArgs struct {
	Package  Package
	Name     string
	Location locs.Loc
	Type     TypeDesc
	Const    bool
}

type ValueFactory interface {
	NewValue(args ValueArgs) Value
	Values() collections.ReadonlySortedSet[Value]
}
