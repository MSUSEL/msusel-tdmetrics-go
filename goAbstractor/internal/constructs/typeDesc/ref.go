package typeDesc

import "github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/jsonify"

type Ref struct {
	Ref string
}

func (tr *Ref) _isTypeDesc() {}

func (tr *Ref) ToJson(ctx *jsonify.Context) jsonify.Datum {
	return jsonify.New(ctx, tr.Ref)
}

func (tr *Ref) String() string {
	return jsonify.ToString(tr)
}
