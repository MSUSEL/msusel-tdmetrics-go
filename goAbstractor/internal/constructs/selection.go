package constructs

import (
	"github.com/Snow-Gremlin/goToolbox/collections"
)

// Selection is a reference to a field, method, parameter, result, etc
// that is being selected from an origin construct.
type Selection interface {
	Construct
	TempReferenceContainer
	IsSelection()

	Name() string
	Origin() Construct
	Target() Construct
}

type SelectionArgs struct {
	// Name is the name of the field, method, parameter, result, etc
	// that is being selected in the origin.
	Name string

	// Origin is the construct that is being selected from.
	Origin Construct

	// Target is the construct by the given name in the given origin.
	// This may be a field, method, variable, constant, etc.
	// This may be nil if it can't be found.
	Target Construct
}

type SelectionFactory interface {
	Factory
	NewSelection(args SelectionArgs) Selection
	Selections() collections.ReadonlySortedSet[Selection]
}
