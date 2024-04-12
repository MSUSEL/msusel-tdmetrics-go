package typeDesc

import "encoding/json"

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

func (tp *Map) _isTypeDesc() {}
