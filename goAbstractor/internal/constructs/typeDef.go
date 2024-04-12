package constructs

import (
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs/typeDesc"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/jsonify"
)

type TypeDef struct {
	Name string
	Type typeDesc.TypeDesc
}

func (td *TypeDef) ToJson(ctx jsonify.Context) jsonify.Datum {
	return jsonify.NewMap().
		AddNonZero(ctx, `name`, td.Name).
		AddNonZero(ctx, `type`, td.Type)
}
