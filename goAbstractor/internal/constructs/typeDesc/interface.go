package typeDesc

import "encoding/json"

type Interface struct {
	Methods []*Func
}

func (ti *Interface) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]any{
		`kind`:    `interface`,
		`methods`: ti.Methods,
	})
}

func (ti *Interface) _isTypeDesc() {}

type Func struct {
	Name      string     `json:"name"`
	Signature *Signature `json:"signature"`
}
