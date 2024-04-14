package jsonify

import "encoding/json"

type Jsonable interface {
	ToJson(ctx *Context) Datum
}

func Marshal(ctx *Context, j Jsonable) []byte {
	data := j.ToJson(ctx)
	var b []byte
	var err error
	if ctx.GetBool(`minimize`) {
		b, err = json.MarshalIndent(data, ``, `  `)
	} else {
		b, err = json.Marshal(data)
	}
	if err != nil {
		panic(err)
	}
	return b
}

func ToString(j Jsonable) string {
	ctx := NewContext()
	b, err := json.Marshal(j.ToJson(ctx))
	if err != nil {
		return err.Error()
	}
	return string(b)
}
