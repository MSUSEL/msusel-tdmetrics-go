package typeDesc

import (
	"encoding/json"

	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs/typeDesc/wrapKind"
)

type Wrap struct {
	Kind wrapKind.WrapKind
	Elem TypeDesc
}

func (tw *Wrap) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]any{
		`kind`: string(tw.Kind),
		`elem`: tw.Elem,
	})
}

func (tw *Wrap) _isTypeDesc() {}
