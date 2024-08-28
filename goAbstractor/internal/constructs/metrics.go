package constructs

import (
	"github.com/Snow-Gremlin/goToolbox/collections"

	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/locs"
)

type Metrics interface {
	Construct
	IsMetrics()

	Location() locs.Loc
}

// MetricsArgs are measurements taken for a method.
type MetricsArgs struct {
	// Location is the unique location for the top level function or method.
	Location locs.Loc

	// Complexity is the McCabe's Cyclomatic Complexity value for the method.
	Complexity int

	// LineCount is the number of lines in a method
	// including function definition line and close bracket line.
	LineCount int

	// CodeCount is the number of lines that have code on it
	// including function definition line and close bracket line
	// but excluding blank lines and comment lines.
	CodeCount int

	// Indent is the indent complexity count using the amount of whitespace
	// to the left of the any line with code on it.
	Indents int
}

type MetricsFactory interface {
	NewMetrics(args MetricsArgs) Metrics
	Metrics() collections.ReadonlySortedSet[Metrics]
}
