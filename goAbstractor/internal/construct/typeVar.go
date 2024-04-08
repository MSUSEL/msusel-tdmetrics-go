package construct

import "encoding/json"

type TypeVar struct {
	Name string
	Type TypeDesc
}

func (v *TypeVar) MarshalJSON() ([]byte, error) {
	data := map[string]any{
		`type`: v.Type,
	}
	if len(v.Name) > 0 {
		data[`name`] = v.Name
	}
	return json.Marshal(data)
}
