package jsonify

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/Snow-Gremlin/goToolbox/terrors/terror"
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

var seekKVPattern = utils.LazyRegex(`^\s*(\w+)\s*=(.+)$`)

func (l *List) Seek(path []any) Datum {
	if len(path) <= 0 {
		return l
	}

	if index, ok := path[0].(int); ok {
		if index < 0 || index >= len(l.data) {
			panic(terror.New(`index is out of bounds`).
				With(`count`, len(l.data)).
				With(`path`, path))
		}
		return l.data[index].Seek(path[1:])
	}

	if kv, ok := path[0].(string); ok {
		parts := seekKVPattern().FindAllStringSubmatch(kv, -1)
		if len(parts) != 1 || len(parts[0]) != 3 {
			panic(terror.New(`invalid key/value in path`).
				With(`pattern`, `^(\w+)=(.+)$`).
				With(`gotten`, parts).
				With(`path`, path))
		}

		key := parts[0][1]
		value := strings.TrimSpace(parts[0][2])
		if strings.HasPrefix(value, `"`) {
			var err error
			value, err = strconv.Unquote(value)
			if err != nil {
				panic(terror.New(`must have a value in key/value that is unquotable`, err).
					With(`value`, value))
			}
		}

		foundKey := false
		for _, item := range l.data {
			if m, ok := item.(*Map); ok {
				if v := m.Get(key); !utils.IsNil(v) {
					foundKey = true
					if fmt.Sprint(v.RawValue()) == value {
						return item
					}
				}
			}
		}

		if foundKey {
			panic(terror.New(`no value found with the given key`).
				With(`value`, value).
				With(`key`, key).
				With(`path`, path))
		}
		panic(terror.New(`no key found`).
			With(`key`, key).
			With(`path`, path))
	}

	panic(terror.New(`must have an index (int) or key/value (string)`).
		With(`path`, path[0]))
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
