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
	// This is used as a key to determine different metrics.
	//
	// Metrics for values can be attached to zero or more values,
	// (`var _ = func() int { ** }`) or (`var x, y = func()(int, int) { ** }`).
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
	// The expression must be a single function with or without a receiver,
	// have no parameter and a single result, and only a return with the
	// right hand side only identifiers (`f`), selectors (`f.x.a`), literals,
	// reference/dereference, or explicit/implicit casts.
	//
	// e.g. `func (f Foo) GetX() int { return f.x }`
	//
	// e.g. `func (f Foo) Kind() string { return "literal" }`
	//
	// This will not return true for modified result getters, such as
	// offsetting an index (`func (f Foo) GetX() int { return f.x + 1 }`),
	// reading a flag (`func (f Foo) GetX() int { return (f.x & 0xFF) >> 2 }`),
	// indexing (`func (f Foo) GetFirst() int { return f.list[0] }`),
	// or creating an instance of a type.
	Getter bool

	// Setter indicates the expression only contains an optional single parameter
	// and no results that is used in a single assignment of an external value.
	//
	// The expression must be a single function with or without a receiver,
	// with zero or one parameters, if the parameter is given it may only be
	// used on the right hand side of the assignment, a single assignment,
	// with only identifiers (`f`), selectors (`f.x.a`), literals,
	// reference/dereference, or explicit/implicit casts.
	//
	// e.g. `func (f *Foo) SetX(x int) { f.x = x }`
	//
	// e.g. `func (f *Foo) SetAsHidden() { f.state = "hidden" }`
	//
	// This will not return true for setters that modify the value in some way,
	// such as setting an offset index, setting an indexed value, etc.
	// This will not return true for reverse setters
	// (`func (f Foo) SetBar(b *Bar) { *b = f.x }`).
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
