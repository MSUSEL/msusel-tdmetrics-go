package constructs

import (
	"github.com/Snow-Gremlin/goToolbox/collections"
)

// Usage is a reference of type usage inside of expressions.
//
// The usage describes:
//   - the read or write of a field, variable, or parameter
//   - the read or write of a function pointer
//   - the invocation/call of a method, function, function pointer, or closure
//   - the definition of an expression local type
//   - the casting or type checking of one type to another
//   - the creation of a type
type Usage interface {
	Construct
	IsUsage()

	Target() Construct
	Origin() Usage
	HasOrigin() bool
}

type UsageArgs struct {

	// Target is the construct that is being used.
	Target Construct

	// Origin is the optional construct that this usage is part of.
	// Origin may be a receiver, parent, or root type.
	Origin Usage
}

type UsageFactory interface {
	NewUsage(args UsageArgs) Usage
	Usages() collections.ReadonlySortedSet[Usage]
}
