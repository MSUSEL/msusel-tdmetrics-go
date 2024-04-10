package typeDesc

import (
	"encoding/json"

	"github.com/Snow-Gremlin/goToolbox/utils"
)

type Interface struct {
	Methods []*Func
}

func (ti *Interface) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]any{
		`kind`:    `interface`,
		`methods`: ti.Methods,
	})
}

func (ti *Interface) Equal(other TypeDesc) bool {
	if utils.IsNil(ti) || utils.IsNil(other) {
		return utils.IsNil(ti) && utils.IsNil(other)
	}
	t2, ok := other.(*Interface)
	return ok && listEqual(ti.Methods, t2.Methods)
}
