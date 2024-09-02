package constructs

import (
	"github.com/Snow-Gremlin/goToolbox/collections"

	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/locs"
)

type Metrics interface {
	Construct
	IsMetrics()

	Location() locs.Loc
	Complexity() int
	LineCount() int
	CodeCount() int
	Indents() int
	Getter() bool
	Setter() bool

	Reads() collections.ReadonlySortedSet[Usage]
	Writes() collections.ReadonlySortedSet[Usage]
	Invokes() collections.ReadonlySortedSet[Usage]
	Defines() collections.ReadonlySortedSet[Usage]
}

// MetricsArgs are measurements taken for a method body or expression.
type MetricsArgs struct {

	// Location is the unique location for the expression.
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

	// Getter indicates the expression only contains no parameters and a
	// single result that is set in a single return of a value.
	//
	// e.g. `func (f Foo) GetX() int { return f.x }`
	//
	// This will not return true for modified result getters,
	// such as offsetting an index,
	// e.g. `func (f Foo) GetX() int { return f.x + 1 }`,
	// or reading a flag,
	// e.g. `func (f Foo) GetX() int { return (f.x & 0xFF) >> 2 }`.
	Getter bool

	// Setter indicates the expression only contains a single parameter
	// and no results that is used in a single assignment of an external value.
	//
	// e.g. `func (f *Foo) SetX(x int) { f.x = x }`
	//
	// This will not return true for setters that modify the value
	// in some way, such as setting an offset index.
	Setter bool

	// Reads are the usages that reads a value or type.
	//
	// e.g. `return point.x + 4` has a read usage of `point.x`.
	Reads collections.SortedSet[Usage]

	// Writes are the usages that writes a value or type.
	// This includes creating a type internal to the function and
	// casting from one type to another.
	//
	// e.g. `point.x = 6` has a write usage of `point.x`.
	//
	// e.g. `return PointType(point)` has a write usage to `PointType`
	//      for the cast as if `point` was first written to a `PointType`
	//      instance prior ot other statements.
	//
	// e.g. `return point.(RealType)` has a write usage to `RealType`
	//      just like the cast in the above example.
	Writes collections.SortedSet[Usage]

	// Invokes are the usages that calls a method, function, or
	// function pointer.
	//
	// e.g. `point.getX()` has an invocation of `point.getX`.
	Invokes collections.SortedSet[Usage]

	// Defines are the usages that define a type internal to the expression.
	// This usage will reference the defined type description.
	//
	// e.g. `var x = struct{x, y int}{x: 42, y: 21}` defines `struct{x, y int}`.
	//
	// e.g. `var doThing = func(x, y int) { ** }` defines `func(x, y int)`.
	//      Note that the internal function will not have its own metrics,
	//      instead it will be counted in this metrics. Meaning a type defined
	//      (or read, write, etc) inside of an internal defined function
	//      will be part of this metrics.
	Defines collections.SortedSet[Usage]
}

type MetricsFactory interface {
	NewMetrics(args MetricsArgs) Metrics
	Metrics() collections.ReadonlySortedSet[Metrics]
}
