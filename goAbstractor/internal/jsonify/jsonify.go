package jsonify

import "encoding/json"

type Jsonable interface {
	ToJson(ctx *Context) Datum
}

func Marshal(ctx *Context, data any) ([]byte, error) {
	datum := New(ctx, data)
	if ctx.IsMinimized() {
		return json.Marshal(datum)
	}
	return json.MarshalIndent(datum, ``, `  `)
}

func ToString(data any) string {
	ctx := NewContext().SetMinimize(true)
	b, err := Marshal(ctx, data)
	if err != nil {
		return err.Error()
	}
	return string(b)
}
