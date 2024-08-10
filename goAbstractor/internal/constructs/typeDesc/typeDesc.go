package typeDesc

import (
	"go/types"

	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs"
)

// TypeDesc is a description of a type.
type TypeDesc interface {
	constructs.Construct

	// GoType gets the Go type associated with this type desc.
	GoType() types.Type
}
