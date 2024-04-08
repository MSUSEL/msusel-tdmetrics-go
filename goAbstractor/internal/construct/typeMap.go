package construct

import "encoding/json"

type TypeMap struct {
	Key   TypeDesc
	Value TypeDesc
}

func (tm *TypeMap) MarshalJSON() ([]byte, error) {
	data := map[string]any{
		`kind`:  `map`,
		`key`:   tm.Key,
		`value`: tm.Value,
	}
	return json.Marshal(data)
}
