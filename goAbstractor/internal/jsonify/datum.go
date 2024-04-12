package jsonify

import (
	"fmt"
	"reflect"
)

type Datum interface {
	_jsonData()
	isZero() bool
}

func New(ctx Context, value any) Datum {
	switch v := value.(type) {
	case nil:
		return newNull()
	case bool:
		return newValue(v)
	case string:
		return newValue(v)
	case int:
		return newValue(v)
	case int8:
		return newValue(v)
	case int16:
		return newValue(v)
	case int32:
		return newValue(v)
	case int64:
		return newValue(v)
	case uint:
		return newValue(v)
	case uint8:
		return newValue(v)
	case uint16:
		return newValue(v)
	case uint32:
		return newValue(v)
	case uint64:
		return newValue(v)
	case float32:
		return newValue(v)
	case float64:
		return newValue(v)
	case Datum:
		return v
	case Jsonable:
		return v.ToJson(ctx)
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
				panic(fmt.Errorf(`invalid key type: %T`, k))
			}
			m.Add(ctx, s, it.Value().Interface())
		}
		return m
	}

	panic(fmt.Errorf(`invalid json type: %T`, value))
}
