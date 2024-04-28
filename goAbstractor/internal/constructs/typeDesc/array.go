package typeDesc

import "github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/jsonify"

type Array struct {
	Length int
	Elem   TypeDesc
}

func (ta *Array) _isTypeDesc() {}

func (ta *Array) ToJson(ctx *jsonify.Context) jsonify.Datum {
	return jsonify.NewMap().
		Add(ctx, `kind`, `array`).
		Add(ctx, `length`, ta.Length).
		Add(ctx, `elem`, ta.Elem)
}

func (ta *Array) String() string {
	return jsonify.ToString(ta)
}
