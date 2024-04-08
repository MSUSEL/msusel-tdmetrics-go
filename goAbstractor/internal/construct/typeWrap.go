package construct

import "encoding/json"

type TypeWrap struct {
	// TODO: Kind should probably be a set of constants
	Kind string
	Elem TypeDesc
}

func (tw *TypeWrap) MarshalJSON() ([]byte, error) {
	data := map[string]any{
		`kind`: tw.Kind,
		`elem`: tw.Elem,
	}
	return json.Marshal(data)
}
