package construct

import (
	"encoding/json"

	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/construct/typeKind"
)

type TypeWrap struct {
	Kind typeKind.TypeKind
	Elem TypeDesc
}

func (tw *TypeWrap) MarshalJSON() ([]byte, error) {
	data := map[string]any{
		`kind`: string(tw.Kind),
		`elem`: tw.Elem,
	}
	return json.Marshal(data)
}
