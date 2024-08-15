package constructs

import "github.com/Snow-Gremlin/goToolbox/collections"

// Abstract is a named signature in an interface.
//
// The order of the abstract doesn't matter.
type Abstract interface {
	Construct
	IsAbstract()

	Name() string
	Signature() Signature
}

type AbstractArgs struct {
	Name      string
	Signature Signature
}

type AbstractFactory interface {
	NewAbstract(args AbstractArgs) Abstract
	Abstracts() collections.ReadonlySortedSet[Abstract]
}
