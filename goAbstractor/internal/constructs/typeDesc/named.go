package typeDesc

import "github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/jsonify"

type Named struct {
	Name string
}

func NewNamed(name string) *Named {
	return &Named{Name: name}
}

func (t *Named) _isTypeDesc() {}

func (t *Named) ToJson(ctx *jsonify.Context) jsonify.Datum {
	return jsonify.New(ctx, t.Name)
}

func (t *Named) String() string {
	return jsonify.ToString(t)
}
