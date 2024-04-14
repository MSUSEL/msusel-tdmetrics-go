package constructs

import (
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs/typeDesc"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/jsonify"
)

type TypeDef struct {
	Name     string
	Type     typeDesc.TypeDesc
	Inherits []*typeDesc.Interface
	Methods  []*Method
}

func (td *TypeDef) ToJson(ctx *jsonify.Context) jsonify.Datum {
	return jsonify.NewMap().
		Add(ctx, `name`, td.Name).
		Add(ctx, `type`, td.Type).
		AddNonZero(ctx, `inherits`, td.Inherits).
		AddNonZero(ctx, `methods`, td.Methods)
}
