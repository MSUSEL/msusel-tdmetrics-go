package construct

import "encoding/json"

type TypeDef struct {
	Name string
	Type TypeDesc
}

func (td *TypeDef) MarshalJSON() ([]byte, error) {
	data := map[string]any{
		`name`: td.Name,
		`type`: td.Type,
	}
	return json.Marshal(data)
}
