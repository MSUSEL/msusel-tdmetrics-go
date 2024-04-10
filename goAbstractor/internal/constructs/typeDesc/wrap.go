package typeDesc

import (
	"encoding/json"

	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs/typeDesc/wrapKind"
	"github.com/Snow-Gremlin/goToolbox/utils"
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

func (tw *Wrap) Equal(other TypeDesc) bool {
	if utils.IsNil(tw) || utils.IsNil(other) {
		return utils.IsNil(tw) && utils.IsNil(other)
	}
	t2, ok := other.(*Wrap)
	return ok && tw.Kind == t2.Kind && tw.Elem.Equal(t2.Elem)
}
