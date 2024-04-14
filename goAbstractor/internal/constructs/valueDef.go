package constructs

import (
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs/typeDesc"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/jsonify"
)

type ValueDef struct {
	Name  string
	Const bool
	Type  typeDesc.TypeDesc
}

func (vd *ValueDef) ToJson(ctx *jsonify.Context) jsonify.Datum {
	return jsonify.NewMap().
		Add(ctx, `name`, vd.Name).
		AddNonZero(ctx, `const`, vd.Const).
		Add(ctx, `type`, vd.Type)
}

func (vd *ValueDef) String() string {
	return jsonify.ToString(vd)
}
