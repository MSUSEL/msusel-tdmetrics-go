package construct

import "encoding/json"

type TypeRef struct {
	Ref string
}

func (tr *TypeRef) MarshalJSON() ([]byte, error) {
	return json.Marshal(tr.Ref)
}
