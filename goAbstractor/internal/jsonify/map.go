package jsonify

import (
	"encoding/json"

	"github.com/Snow-Gremlin/goToolbox/utils"
)

type Map struct {
	data map[string]Datum
}

func NewMap() *Map {
	return &Map{
		data: map[string]Datum{},
	}
}

func (m *Map) _jsonData() {}

func (m *Map) isZero() bool {
	return m == nil || len(m.data) <= 0
}

func (m *Map) Get(key string) Datum {
	return m.data[key]
}

func (m *Map) Seek(path []any) Datum {
	return m.subSeek(newSeeker(path))
}

func (m *Map) subSeek(s *seeker) Datum {
	if s.done() {
		return m
	}

	if s.isCount() {
		return newValue(len(m.data))
	}

	single, selector := s.asSelector()

	sub := NewMap()
	keys := utils.SortedKeys(m.data)
	for _, key := range keys {
		if selector(key) {
			e := m.data[key].subSeek(s.next())
			if single {
				return e
			}
			sub.data[key] = e
		}
	}
	return sub
}

func (m *Map) Add(ctx *Context, key string, value any) *Map {
	if m == nil {
		m = NewMap()
	}
	m.data[key] = New(ctx, value)
	return m
}

func (m *Map) AddIf(ctx *Context, test bool, key string, value any) *Map {
	if test {
		return m.Add(ctx, key, value)
	}
	return m
}

func (m *Map) AddNonZero(ctx *Context, key string, value any) *Map {
	if d := New(ctx, value); !d.isZero() {
		if m == nil {
			m = NewMap()
		}
		m.data[key] = d
	}
	return m
}

func (m *Map) RawValue() any {
	return m.data
}

func (m *Map) MarshalJSON() ([]byte, error) {
	return json.Marshal(m.data)
}
