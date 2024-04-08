package construct

import "encoding/json"

type TypeFunc struct {
	Name      string
	Signature *Signature
}

func (tf *TypeFunc) MarshalJSON() ([]byte, error) {
	data := map[string]any{
		`name`:      tf.Name,
		`signature`: tf.Signature,
	}
	return json.Marshal(data)
}
