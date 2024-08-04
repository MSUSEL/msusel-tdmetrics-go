package constructs

import (
	"cmp"
	"strings"

	"github.com/Snow-Gremlin/goToolbox/utils"

	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs/kind"
)

// Construct is a single unit of a software design or type.
type Construct interface {
	// compareTo compares this construct to another.
	// The other should be of the same kind when this compare is called.
	compareTo(other Construct) int

	// Kind gets the kind of the construct.
	Kind() kind.Kind

	// SetIndex sets the unique index of construct.
	// Indices will be 1 based so that 0 is unset.
	setIndex(index int)
}

// or will return the first non-zero value returned
// by a comparison or it will return zero.
func or(comparisons ...func() int) int {
	for _, cmp := range comparisons {
		if c := cmp(); c != 0 {
			return c
		}
	}
	return 0
}

func ternary[T any](test bool, first, second T) T {
	if test {
		return first
	}
	return second
}

func boolCompare(a, b bool) int {
	return ternary(b, 1, 0) - ternary(a, 1, 0)
}

// Compare two constructs together.
func Compare[T Construct](a, b T) int {
	m, n := utils.IsNil(a), utils.IsNil(b)
	if m && n {
		return 0
	}
	return or(
		func() int { return boolCompare(m, n) },
		func() int { return strings.Compare(string(a.Kind()), string(b.Kind())) },
		func() int { return a.compareTo(b) },
	)
}

// compareSlice two slices of constructs together.
func compareSlice[T Construct, S ~[]T](a, b S) int {
	ca, cb := len(a), len(b)
	cMin := min(ca, cb)
	for i := 0; i < cMin; i++ {
		if cmp := Compare(a[i], b[i]); cmp != 0 {
			return cmp
		}
	}
	return cmp.Compare[int](ca, cb)
}
