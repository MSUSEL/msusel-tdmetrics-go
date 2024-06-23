package jsonify

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

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
			panic(fmt.Errorf(`must have an index in [0 .. %d): %v`, len(l.data), path[0]))
		}
		return l.data[index].Seek(path[1:])
	}

	if kv, ok := path[0].(string); ok {
		parts := seekKVPattern().FindAllStringSubmatch(kv, -1)
		if len(parts) != 1 || len(parts[0]) != 3 {
			panic(fmt.Errorf(`must have key/value string match '^(\w+)=(.+)$' got %v: %v`, parts, path[0]))
		}

		key := parts[0][1]
		value := strings.TrimSpace(parts[0][2])
		if strings.HasPrefix(value, `"`) {
			var err error
			value, err = strconv.Unquote(value)
			if err != nil {
				panic(fmt.Errorf(`must have a value in key/value that is unquotable, %q: %w`, value, err))
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
			panic(fmt.Errorf(`no value found for %q with the given key %q: %v`, value, key, path[0]))
		}
		panic(fmt.Errorf(`no key found for %q: %v`, key, path[0]))
	}

	panic(fmt.Errorf(`must have an index (int) or key/value (string): %v`, path[0]))
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
