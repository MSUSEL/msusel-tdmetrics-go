package typeDesc

import "github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/jsonify"

type Chan struct {
	Elem TypeDesc
}

func (tc *Chan) _isTypeDesc() {}

func (tc *Chan) ToJson(ctx *jsonify.Context) jsonify.Datum {
	return jsonify.NewMap().
		Add(ctx, `kind`, `chan`).
		Add(ctx, `elem`, tc.Elem)
}

func (tc *Chan) String() string {
	return jsonify.ToString(tc)
}
