package construct

import "encoding/json"

type Interface struct {
	Methods []*TypeFunc
}

func (ti *Interface) MarshalJSON() ([]byte, error) {
	data := map[string]any{
		`kind`:    `interface`,
		`methods`: ti.Methods,
	}
	return json.Marshal(data)
}
