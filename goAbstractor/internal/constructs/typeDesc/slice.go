package typeDesc

import "github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/jsonify"

type Slice struct {
	Elem TypeDesc
}

func (ts *Slice) _isTypeDesc() {}

func (ts *Slice) ToJson(ctx *jsonify.Context) jsonify.Datum {
	return jsonify.NewMap().
		Add(ctx, `kind`, `slice`).
		Add(ctx, `elem`, ts.Elem)
}

func (ts *Slice) String() string {
	return jsonify.ToString(ts)
}
