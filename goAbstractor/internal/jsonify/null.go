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

func (n *null) Seek(path []any) Datum {
	return newSeeker(path).StepNull(n)
}

func (n *null) RawValue() any {
	return nil
}

func (n *null) MarshalJSON() ([]byte, error) {
	return json.Marshal(nil)
}
