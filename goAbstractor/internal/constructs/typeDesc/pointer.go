package typeDesc

import "github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/jsonify"

type Pointer struct {
	Elem TypeDesc
}

func (tp *Pointer) _isTypeDesc() {}

func (tp *Pointer) ToJson(ctx *jsonify.Context) jsonify.Datum {
	return jsonify.NewMap().
		Add(ctx, `kind`, `pointer`).
		Add(ctx, `elem`, tp.Elem)
}

func (tp *Pointer) String() string {
	return jsonify.ToString(tp)
}
