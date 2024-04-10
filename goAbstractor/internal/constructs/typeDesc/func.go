package typeDesc

import (
	"encoding/json"

	"github.com/Snow-Gremlin/goToolbox/utils"
)

type Func struct {
	Name      string
	Signature *Signature
}

func (tf *Func) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]any{
		`name`:      tf.Name,
		`signature`: tf.Signature,
	})
}

func (tw *Func) Equal(other TypeDesc) bool {
	if utils.IsNil(tw) || utils.IsNil(other) {
		return utils.IsNil(tw) && utils.IsNil(other)
	}
	t2, ok := other.(*Func)
	return ok && tw.Name == t2.Name && tw.Signature.Equal(t2.Signature)
}
