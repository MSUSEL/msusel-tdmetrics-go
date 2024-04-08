package construct

import "encoding/json"

type Method struct {
	Name      string
	Signature *Signature
	Receiver  TypeDesc
}

func (m *Method) MarshalJSON() ([]byte, error) {
	data := map[string]any{
		`name`:      m.Name,
		`signature`: m.Signature,
	}
	if m.Receiver != nil {
		data[`receiver`] = m.Receiver
	}
	return json.Marshal(data)
}
