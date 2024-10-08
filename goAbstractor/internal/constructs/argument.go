package constructs

import "github.com/Snow-Gremlin/goToolbox/collections"

// Argument is a parameter or result in a method signature.
//
// The order of the arguments matters.
type Argument interface {
	Construct
	TempReferenceContainer
	IsArgument()

	Name() string
	Type() TypeDesc
}

type ArgumentArgs struct {
	Name string
	Type TypeDesc
}

type ArgumentFactory interface {
	NewArgument(args ArgumentArgs) Argument
	Arguments() collections.ReadonlySortedSet[Argument]
}
