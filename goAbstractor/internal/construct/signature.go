package construct

import "encoding/json"

type Signature struct {
	ShowKind   bool
	Variadic   bool
	Params     []*TypeVar
	Returns    []*TypeVar
	TypeParams []*TypeParam
}

func (sig *Signature) MarshalJSON() ([]byte, error) {
	data := map[string]any{}
	if sig.ShowKind {
		data[`kind`] = `signature`
	}
	if sig.Variadic {
		data[`variadic`] = true
	}
	if len(sig.Params) > 0 {
		data[`params`] = sig.Params
	}
	if len(sig.Returns) > 0 {
		data[`returns`] = sig.Returns
	}
	if len(sig.TypeParams) > 0 {
		data[`typeParams`] = sig.TypeParams
	}
	return json.Marshal(data)
}
