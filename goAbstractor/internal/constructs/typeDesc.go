package constructs

import "go/types"

// TypeDesc is an interface for all type descriptors.
type TypeDesc interface {
	Construct

	// GoType gets the Go type associated with this type desc.
	GoType() types.Type
}
