package jsonify

import (
	"encoding/json"

	"github.com/Snow-Gremlin/goToolbox/utils"
)

type valueConstraint interface {
	~bool | ~string |
		~int | ~int8 | ~int16 | ~int32 | ~int64 |
		~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 |
		~float32 | ~float64
}

type Value struct {
	data any
}

func NewValue[T valueConstraint](data T) *Value {
	return &Value{data: data}
}

func NewNull() *Value {
	return &Value{data: nil}
}

func (v *Value) _jsonData() {}

func (v *Value) isZero() bool {
	return utils.IsZero(v.data)
}

func (v *Value) Seek(path []any) Datum {
	return v.subSeek(newSeeker(path))
}

func (v *Value) subSeek(s *seeker) Datum {
	if s.done() {
		return v
	}

	if s.isCount() {
		return NewValue(1)
	}

	panic(s.fail(`path continues from a value`).
		With(`value`, v.data))
}

func (v *Value) RawValue() any {
	return v.data
}

func (v *Value) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.data)
}
