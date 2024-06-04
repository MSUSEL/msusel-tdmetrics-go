package constructs

import (
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/jsonify"
)

type ValueDef interface{}

type valueDefImp struct {
	name    string
	isConst bool
	typ     TypeDesc
}

func NewValueDef(name string, isConst bool, typ TypeDesc) ValueDef {
	return &valueDefImp{
		name:    name,
		isConst: isConst,
		typ:     typ,
	}
}

func (vd *valueDefImp) ToJson(ctx *jsonify.Context) jsonify.Datum {
	ctx2 := ctx.HideKind().Short()
	return jsonify.NewMap().
		Add(ctx2, `name`, vd.name).
		AddNonZero(ctx2, `const`, vd.isConst).
		Add(ctx2, `type`, vd.typ)
}

func (vd *valueDefImp) String() string {
	return jsonify.ToString(vd)
}
