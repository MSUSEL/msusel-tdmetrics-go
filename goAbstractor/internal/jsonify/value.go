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

type value[T valueConstraint] struct {
	data T
}

func newValue[T valueConstraint](data T) *value[T] {
	return &value[T]{data: data}
}

func (v *value[T]) _jsonData() {}

func (v *value[T]) isZero() bool {
	return utils.IsZero(v.data)
}

func (v *value[T]) Seek(path []any) Datum {
	return v.subSeek(newSeeker(path))
}

func (v *value[T]) subSeek(s *seeker) Datum {
	if s.done() {
		return v
	}

	if s.isCount() {
		return newValue(1)
	}

	panic(s.fail(`path continues from a value`).
		With(`value`, v.data))
}

func (v *value[T]) RawValue() any {
	return v.data
}

func (v *value[T]) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.data)
}
