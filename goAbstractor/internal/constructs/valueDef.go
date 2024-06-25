package constructs

import (
	"fmt"

	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/jsonify"
	"github.com/Snow-Gremlin/goToolbox/utils"
)

type ValueDef interface {
	Visitable
}

type valueDefImp struct {
	name    string
	isConst bool
	typ     TypeDesc
}

func NewValueDef(name string, isConst bool, typ TypeDesc) ValueDef {
	if len(name) <= 0 || utils.IsNil(typ) {
		constStr := ``
		if isConst {
			constStr = `const `
		}
		panic(fmt.Errorf(`must have a name and type for a value definition: %s%q %v`,
			constStr, name, typ))
	}

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

func (vd *valueDefImp) Visit(v Visitor) {
	visitTest(v, vd.typ)
}

func (vd *valueDefImp) String() string {
	return jsonify.ToString(vd)
}
