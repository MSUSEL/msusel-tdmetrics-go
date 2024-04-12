package typeDesc

import "encoding/json"

type Struct struct {
	Fields []*Field
}

func (ts *Struct) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]any{
		`kind`:   `struct`,
		`fields`: ts.Fields,
	})
}

func (ts *Struct) _isTypeDesc() {}

type Field struct {
	Anonymous bool     `json:"anonymous,omitempty"`
	Name      string   `json:"name,omitempty"`
	Type      TypeDesc `json:"type"`
}
