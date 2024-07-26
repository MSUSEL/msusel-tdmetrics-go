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
	return n.subSeek(newSeeker(path))
}

func (n *null) subSeek(s *seeker) Datum {
	if !s.done() {
		return n
	}

	if s.asString() == `#` {
		return newValue(1)
	}

	panic(s.fail(`path continues from a null`))
}

func (n *null) RawValue() any {
	return nil
}

func (n *null) MarshalJSON() ([]byte, error) {
	return json.Marshal(nil)
}
