package jsonify

type Jsonable interface {
	ToJson(ctx Context) Datum
}
