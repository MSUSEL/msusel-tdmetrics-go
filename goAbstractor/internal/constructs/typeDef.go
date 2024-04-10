package constructs

import "github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs/typeDesc"

type TypeDef struct {
	Name string            `json:"name,omitempty"`
	Type typeDesc.TypeDesc `json:"type,omitempty"`
}
