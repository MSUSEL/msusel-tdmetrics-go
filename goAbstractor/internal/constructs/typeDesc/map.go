package typeDesc

import "github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/jsonify"

type Map struct {
	Key   TypeDesc
	Value TypeDesc
}

func (tm *Map) _isTypeDesc() {}

func (tm *Map) ToJson(ctx *jsonify.Context) jsonify.Datum {
	return jsonify.NewMap().
		Add(ctx, `kind`, `map`).
		Add(ctx, `key`, tm.Key).
		Add(ctx, `value`, tm.Value)
}

func (tm *Map) String() string {
	return jsonify.ToString(tm)
}
