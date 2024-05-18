package typeDesc

import "go/types"

// TypeDesc is an interface for all type descriptors.
type TypeDesc interface {

	// GoType gets the Go type associated with this type desc.
	GoType() types.Type
}
