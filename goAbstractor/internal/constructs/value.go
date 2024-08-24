package constructs

import (
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/locs"
	"github.com/Snow-Gremlin/goToolbox/collections"
)

type Value interface {
	Declaration
	IsValue()

	Const() bool
	Metrics() Metrics
}

type ValueArgs struct {
	Package  Package
	Name     string
	Location locs.Loc
	Type     TypeDesc
	Const    bool

	// Metrics are optional and may be nil. These metrics are for
	// a variable initialized with an anonymous function.
	// (e.g. `var x = func() int { ** }()`)
	Metrics Metrics
}

type ValueFactory interface {
	NewValue(args ValueArgs) Value
	Values() collections.ReadonlySortedSet[Value]
}
