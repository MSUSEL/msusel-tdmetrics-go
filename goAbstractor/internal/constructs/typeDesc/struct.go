package typeDesc

import (
	"encoding/json"

	"github.com/Snow-Gremlin/goToolbox/utils"
)

type Struct struct {
	Fields []*Var
}

func (ts *Struct) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]any{
		`kind`:   `struct`,
		`fields`: ts.Fields,
	})
}

func (ts *Struct) Equal(other TypeDesc) bool {
	if utils.IsNil(ts) || utils.IsNil(other) {
		return utils.IsNil(ts) && utils.IsNil(other)
	}
	t2, ok := other.(*Struct)
	return ok && listEqual(ts.Fields, t2.Fields)
}
