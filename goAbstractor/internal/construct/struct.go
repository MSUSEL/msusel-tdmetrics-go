package construct

import "encoding/json"

type Struct struct {
	Fields []*TypeVar
}

func (ts *Struct) MarshalJSON() ([]byte, error) {
	data := map[string]any{
		`kind`:   `struct`,
		`fields`: ts.Fields,
	}
	return json.Marshal(data)
}
