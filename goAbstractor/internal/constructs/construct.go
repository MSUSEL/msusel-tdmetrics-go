package constructs

import (
	"github.com/Snow-Gremlin/goToolbox/utils"

	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs/kind"
)

// Construct is a single unit of a software design or type.
type Construct interface {
	utils.Comparable[Construct]

	// Kind gets the kind of the construct.
	Kind() kind.Kind

	// SetIndex sets the unique index of construct.
	// Indices will be 1 based so that 0 is unset.
	SetIndex(index int)
}

// Compare two constructs together.
func Compare[T Construct](a, b T) int {
	m, n := utils.IsNil(a), utils.IsNil(b)
	if m && n {
		return 0
	}
	if m {
		return -1
	}
	if n {
		return 1
	}

	if cmp := a.Kind().CompareTo(b.Kind()); cmp != 0 {
		return cmp
	}

	return a.CompareTo(b)
}
