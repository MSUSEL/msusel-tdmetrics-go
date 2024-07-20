package jsonify

import (
	"encoding/json"

	"github.com/Snow-Gremlin/goToolbox/terrors/terror"
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
		panic(terror.New(`path continues from null`).
			With(`path`, path))
	}
	return n
}

func (n *null) RawValue() any {
	return nil
}

func (n *null) MarshalJSON() ([]byte, error) {
	return json.Marshal(nil)
}
