package jsonify

import (
	"encoding/json"
	"fmt"
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
	if len(path) <= 0 {
		return l
	}
	index, ok := path[0].(int)
	if !ok || index < 0 || index >= len(l.data) {
		panic(fmt.Errorf(`must have an index in [0 .. %d): %v`, len(l.data), path[0]))
	}
	return l.data[index].Seek(path[1:])
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

func (l *List) MarshalJSON() ([]byte, error) {
	return json.Marshal(l.data)
}
