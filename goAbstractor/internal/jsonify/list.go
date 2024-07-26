package jsonify

import (
	"encoding/json"
	"fmt"

	"github.com/Snow-Gremlin/goToolbox/utils"
)

type List struct {
	data []Datum
}

func NewList() *List {
	return &List{}
}

func (l *List) _jsonData() {}

func (l *List) isZero() bool {
	return l == nil || len(l.data) <= 0
}

func (l *List) Seek(path []any) Datum {
	return l.subSeek(newSeeker(path))
}

func (l *List) subSeek(s *seeker) Datum {
	if s.done() {
		return l
	}

	length := len(l.data)
	if s.asString() == `#` {
		return newValue(length)
	}

	if index, ok := s.asIndex(length); ok {
		return l.data[index].subSeek(s.next())
	}

	if start, end, ok := s.asRange(length); ok {
		sub := NewList()
		for i := start; i <= end; i++ {
			d := l.data[i].subSeek(s.next())
			sub.data = append(sub.data, d)
		}
		return sub
	}

	if key, single, selector, ok := s.asKeyValue(); ok {
		foundKey := false
		sub := NewList()
		for _, item := range l.data {
			if m, ok := item.(*Map); ok {
				if v := m.Get(key); !utils.IsNil(v) {
					foundKey = true
					if selector(fmt.Sprint(v.RawValue())) {
						e := item.subSeek(s.next())
						if single {
							return e
						}
						sub.data = append(sub.data, e)
					}
				}
			}
		}
		if !foundKey {
			panic(s.fail(`no key found`).
				With(`key`, key))
		}
		return sub
	}

	panic(s.fail(`must have an index, range, or key/value`))
}

func (l *List) Append(ctx *Context, values ...any) *List {
	if l == nil {
		l = NewList()
	}
	if count := len(values); count > 0 {
		data := make([]Datum, count)
		for i, value := range values {
			data[i] = New(ctx, value)
		}
		l.data = append(l.data, data...)
	}
	return l
}

func (l *List) AppendNonZero(ctx *Context, values ...any) *List {
	if l == nil {
		l = NewList()
	}
	if count := len(values); count > 0 {
		data := make([]Datum, 0, count)
		for _, value := range values {
			if d := New(ctx, value); !d.isZero() {
				data = append(data, d)
			}
		}
		l.data = append(l.data, data...)
	}
	return l
}

func (l *List) RawValue() any {
	return l.data
}

func (l *List) MarshalJSON() ([]byte, error) {
	return json.Marshal(l.data)
}
