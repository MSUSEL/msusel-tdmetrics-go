package construct

import "encoding/json"

// TODO: Need to rework
// Type param is defined on each param/return and signature right now.
// At minimum the params/returns could be just index references.
// Need to rework to use minimum common interfaces to be like Java.
// This means things like `int` need to have a pseudo interface.

type TypeParam struct {
	Index      int
	Constraint TypeDesc
}

func (tp *TypeParam) MarshalJSON() ([]byte, error) {
	data := map[string]any{
		`kind`:       `typeParam`,
		`fields`:     tp.Index,
		`constraint`: tp.Constraint,
	}
	return json.Marshal(data)
}
