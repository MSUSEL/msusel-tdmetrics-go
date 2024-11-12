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
}

type SelectionArgs struct {
	// Name is the name of the field, method, parameter, result, etc
	// that is being selected in the origin.
	Name string

	// Origin is the construct that is being selected from.
	Origin Construct
}

type SelectionFactory interface {
	NewSelection(args SelectionArgs) Selection
	Selections() collections.ReadonlySortedSet[Selection]
}
