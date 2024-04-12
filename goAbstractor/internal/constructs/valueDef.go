package constructs

import (
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs/typeDesc"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/jsonify"
)

type ValueDef struct {
	Name string
	Type typeDesc.TypeDesc
}

func (vd *ValueDef) ToJson(ctx *jsonify.Context) jsonify.Datum {
	return jsonify.NewMap().
		Add(ctx, `name`, vd.Name).
		Add(ctx, `type`, vd.Type)
}
