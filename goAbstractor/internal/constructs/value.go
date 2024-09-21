package constructs

import (
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/locs"
	"github.com/Snow-Gremlin/goToolbox/collections"
)

type Value interface {
	Declaration
	TempReferenceContainer
	IsValue()

	Const() bool
	Metrics() Metrics

	// HasSideEffect indicates that the value initialization expression
	// contains a usage that affects a variable outside of itself, arguments,
	// or receiver. Or that a function/method invoked in the expression
	// that these metrics are for, has a side effect.
	// If a variable is only read from it is not considered a side effect.
	//
	// For example:
	//    var callCount = 0
	//	  func foo() int {
	//       callCount++
	//       return 42
	//    }
	//
	//    // The below metrics for `bar` would show the invocation of `foo`
	//    // that has a side effect of changing `callCount`.
	//    // Therefore, even if `bar` isn't used, the initialization
	//    // of be is alive because it has side effects.
	//    var bar = foo()
	HasSideEffect() bool
}

type ValueArgs struct {
	Package  Package
	Name     string
	Exported bool
	Location locs.Loc
	Type     TypeDesc
	Const    bool

	// Metrics are optional and may be nil. These metrics are for
	// a variable initialized with an anonymous function.
	// (e.g. `var x = func() int { ⋯ }()`)
	Metrics Metrics
}

type ValueFactory interface {
	NewValue(args ValueArgs) Value
	Values() collections.ReadonlySortedSet[Value]
}
