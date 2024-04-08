package construct

import "encoding/json"

type ValueDef struct {
	Name string
	Type TypeDesc
}

func (vd *ValueDef) MarshalJSON() ([]byte, error) {
	data := map[string]any{
		`name`: vd.Name,
		`type`: vd.Type,
	}
	return json.Marshal(data)
}
