package jsonify

import "encoding/json"

type Jsonable interface {
	ToJson(ctx *Context) Datum
}

func Marshal(ctx *Context, j Jsonable) ([]byte, error) {
	data := j.ToJson(ctx)
	if ctx.IsMinimized() {
		return json.Marshal(data)
	}
	return json.MarshalIndent(data, ``, `  `)
}

func ToString(j Jsonable) string {
	ctx := NewContext()
	b, err := json.Marshal(j.ToJson(ctx))
	if err != nil {
		return err.Error()
	}
	return string(b)
}
