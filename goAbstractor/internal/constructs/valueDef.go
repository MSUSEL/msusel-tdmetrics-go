package constructs

import (
	"encoding/json"

	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs/typeDesc"
)

type ValueDef struct {
	Name string
	Type typeDesc.TypeDesc
}

func (vd *ValueDef) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]any{
		`name`: vd.Name,
		`type`: vd.Type,
	})
}
