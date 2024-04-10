package typeDesc

import (
	"encoding/json"

	"github.com/Snow-Gremlin/goToolbox/utils"
)

type Map struct {
	Key   TypeDesc
	Value TypeDesc
}

func (tm *Map) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]any{
		`kind`:  `map`,
		`key`:   tm.Key,
		`value`: tm.Value,
	})
}

func (tp *Map) Equal(other TypeDesc) bool {
	if utils.IsNil(tp) || utils.IsNil(other) {
		return utils.IsNil(tp) && utils.IsNil(other)
	}
	t2, ok := other.(*Map)
	return ok && tp.Key.Equal(t2.Key) && tp.Value.Equal(t2.Value)
}
