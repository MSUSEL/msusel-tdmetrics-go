package constructs

import "go/types"

// TypeDesc is a description of a type.
type TypeDesc interface {
	Construct
	IsTypeDesc()

	// GoType gets the Go type associated with this type desc.
	GoType() types.Type
}
