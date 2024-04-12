package typeDesc

import "encoding/json"

type Ref struct {
	Ref string
}

func (tr *Ref) MarshalJSON() ([]byte, error) {
	return json.Marshal(tr.Ref)
}

func (tw *Ref) _isTypeDesc() {}
