package constructs

import "github.com/Snow-Gremlin/goToolbox/collections"

// Field is a variable inside of a struct or class.
//
// For the abstraction, the order of the fields and
// the tags of the fields don't matter.
type Field interface {
	Construct
	TempReferenceContainer
	IsField()

	Name() string
	Exported() bool
	Type() TypeDesc
	Embedded() bool
}

type FieldArgs struct {
	Name     string
	Exported bool
	Type     TypeDesc
	Embedded bool
}

type FieldFactory interface {
	Factory
	NewField(args FieldArgs) Field
	Fields() collections.ReadonlySortedSet[Field]
}
