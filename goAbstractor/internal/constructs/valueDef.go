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
	return jsonify.NewMap().
		Add(ctx, `name`, vd.name).
		AddNonZero(ctx, `const`, vd.isConst).
		Add(ctx, `type`, vd.typ)
}

func (vd *valueDefImp) String() string {
	return jsonify.ToString(vd)
}
