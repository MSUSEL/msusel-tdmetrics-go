package typeDesc

import (
	"encoding/json"

	"github.com/Snow-Gremlin/goToolbox/utils"
)

type Var struct {
	Name string
	Type TypeDesc
}

func (v *Var) MarshalJSON() ([]byte, error) {
	data := map[string]any{
		`type`: v.Type,
	}
	if len(v.Name) > 0 {
		data[`name`] = v.Name
	}
	return json.Marshal(data)
}

func (ts *Var) Equal(other TypeDesc) bool {
	if utils.IsNil(ts) || utils.IsNil(other) {
		return utils.IsNil(ts) && utils.IsNil(other)
	}
	t2, ok := other.(*Var)
	return ok && ts.Name == t2.Name && ts.Type.Equal(t2.Type)
}
