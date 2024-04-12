package jsonify

import "encoding/json"

type null struct{}

func newNull() *null {
	return (*null)(nil)
}

func (n *null) _jsonData() {}

func (n *null) isZero() bool {
	return true
}

func (n *null) MarshalJSON() ([]byte, error) {
	return json.Marshal(nil)
}
