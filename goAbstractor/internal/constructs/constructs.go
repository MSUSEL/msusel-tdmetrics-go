package constructs

import (
	"github.com/Snow-Gremlin/goToolbox/comp"
)

// Construct is a description of a type.
type Construct interface {
	comp.Comparable[Construct]

	// Kind gets a string unique to each construct type.
	Kind() string

	// setIndex sets the unique index of construct.
	// Indices will be 1 based so that 0 is unset.
	SetIndex(index int)
}
