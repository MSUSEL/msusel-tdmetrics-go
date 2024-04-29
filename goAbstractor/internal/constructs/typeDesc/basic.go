package typeDesc

import "github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/jsonify"

type Basic struct {
	Name string
}

func (t *Basic) _isTypeDesc() {}

func (t *Basic) ToJson(ctx *jsonify.Context) jsonify.Datum {
	return jsonify.New(ctx, t.Name)
}

func (t *Basic) String() string {
	return jsonify.ToString(t)
}
