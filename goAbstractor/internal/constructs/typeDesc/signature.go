package typeDesc

import (
	"encoding/json"

	"github.com/Snow-Gremlin/goToolbox/utils"
)

type Signature struct {
	Variadic   bool
	Params     []*Var
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

func (sig *Signature) Equal(other TypeDesc) bool {
	if utils.IsNil(sig) || utils.IsNil(other) {
		return utils.IsNil(sig) && utils.IsNil(other)
	}
	t2, ok := other.(*Signature)
	return ok &&
		sig.Variadic == t2.Variadic &&
		listEqual(sig.Params, t2.Params) &&
		sig.Return.Equal(t2.Return) &&
		listEqual(sig.TypeParams, t2.TypeParams)
}
