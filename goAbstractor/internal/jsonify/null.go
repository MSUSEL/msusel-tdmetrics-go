package jsonify

import (
	"encoding/json"
	"fmt"
)

type null struct{}

func newNull() *null {
	return (*null)(nil)
}

func (n *null) _jsonData() {}

func (n *null) isZero() bool {
	return true
}

func (n *null) Seek(path []any) Datum {
	if len(path) > 0 {
		panic(fmt.Errorf(`path continues from null: %v`, path))
	}
	return n
}

func (n *null) RawValue() any {
	return nil
}

func (n *null) MarshalJSON() ([]byte, error) {
	return json.Marshal(nil)
}
