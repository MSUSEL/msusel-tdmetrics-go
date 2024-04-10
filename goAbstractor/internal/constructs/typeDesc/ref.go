package typeDesc

import (
	"encoding/json"

	"github.com/Snow-Gremlin/goToolbox/utils"
)

type Ref struct {
	Ref string
}

func (tr *Ref) MarshalJSON() ([]byte, error) {
	return json.Marshal(tr.Ref)
}

func (tw *Ref) Equal(other TypeDesc) bool {
	if utils.IsNil(tw) || utils.IsNil(other) {
		return utils.IsNil(tw) && utils.IsNil(other)
	}
	t2, ok := other.(*Ref)
	return ok && tw.Ref == t2.Ref
}
