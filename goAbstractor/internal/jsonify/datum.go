package jsonify

import (
	"reflect"

	"github.com/Snow-Gremlin/goToolbox/terrors/terror"
	"github.com/Snow-Gremlin/goToolbox/utils"
)

type Datum interface {
	_jsonData()
	isZero() bool
	Seek(path []any) Datum
	subSeek(s *seeker) Datum
	RawValue() any
}

func New(ctx *Context, value any) Datum {
	if utils.IsNil(value) {
		return NewNull()
	}
	switch v := value.(type) {
	case nil:
		return NewNull()
	case bool:
		return NewValue(v)
	case string:
		return NewValue(v)
	case int:
		return NewValue(v)
	case int8:
		return NewValue(v)
	case int16:
		return NewValue(v)
	case int32:
		return NewValue(v)
	case int64:
		return NewValue(v)
	case uint:
		return NewValue(v)
	case uint8:
		return NewValue(v)
	case uint16:
		return NewValue(v)
	case uint32:
		return NewValue(v)
	case uint64:
		return NewValue(v)
	case float32:
		return NewValue(v)
	case float64:
		return NewValue(v)
	case Datum:
		return v
	case Jsonable:
		d := v.ToJson(ctx)
		if utils.IsNil(d) {
			return NewNull()
		}
		return d
	}

	f := reflect.ValueOf(value)
	switch f.Kind() {
	case reflect.Array, reflect.Slice:
		l := NewList()
		for i := range f.Len() {
			l.Append(ctx, f.Index(i).Interface())
		}
		return l
	case reflect.Map:
		m := NewMap()
		it := f.MapRange()
		for it.Next() {
			k := it.Key().Interface()
			s, ok := k.(string)
			if !ok {
				panic(terror.New(`invalid key type`).
					WithType(`type`, k))
			}
			m.Add(ctx, s, it.Value().Interface())
		}
		return m
	}

	panic(terror.New(`invalid json type`).
		WithType(`type`, value))
}
