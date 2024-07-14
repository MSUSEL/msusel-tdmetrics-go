package jsonify

import "encoding/json"

type Jsonable interface {
	ToJson(ctx *Context) Datum
}

func Marshal(ctx *Context, data any) ([]byte, error) {
	if j, ok := data.(Jsonable); ok {
		data = j.ToJson(ctx)
	}
	if ctx.IsMinimized() {
		return json.Marshal(data)
	}
	return json.MarshalIndent(data, ``, `  `)
}

func ToString(data any) string {
	ctx := NewContext().SetMinimize(true)
	b, err := Marshal(ctx, data)
	if err != nil {
		return err.Error()
	}
	return string(b)
}
