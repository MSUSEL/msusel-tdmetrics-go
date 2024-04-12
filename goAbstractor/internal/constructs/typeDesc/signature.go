package typeDesc

import "encoding/json"

type Signature struct {
	Variadic   bool
	Params     []*Param
	Return     TypeDesc
	TypeParams []*TypeParam
}

func (sig *Signature) MarshalJSON() ([]byte, error) {
	data := map[string]any{
		`kind`: `signature`,
	}
	if sig.Variadic {
		data[`variadic`] = true
	}
	if len(sig.Params) > 0 {
		data[`params`] = sig.Params
	}
	if sig.Return != nil {
		data[`return`] = sig.Return
	}
	if len(sig.TypeParams) > 0 {
		data[`typeParams`] = sig.TypeParams
	}
	return json.Marshal(data)
}

func (sig *Signature) _isTypeDesc() {}

type Param struct {
	Name string   `json:"name,omitempty"`
	Type TypeDesc `json:"type"`
}
