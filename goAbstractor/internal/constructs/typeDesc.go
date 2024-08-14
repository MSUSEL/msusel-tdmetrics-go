package constructs

import "go/types"

// TypeDesc is a description of a type.
type TypeDesc interface {
	Construct

	// IsTypeDesc indicates that the type is a TypeDesc at compile time.
	// This prevents anything else from duck-typing into a TypeDecs.
	IsTypeDesc()

	// GoType gets the Go type associated with this type desc.
	GoType() types.Type
}
