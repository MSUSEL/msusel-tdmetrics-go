package constructs

import (
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs/typeDesc"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/jsonify"
)

type TypeDef struct {
	Name       string
	Type       typeDesc.TypeDesc
	Methods    []*Method
	TypeParams []typeDesc.Named
	Interface  typeDesc.Interface
}

func NewTypeDef(name string, t typeDesc.TypeDesc) *TypeDef {
	return &TypeDef{
		Name: name,
		Type: t,
	}
}

func (td *TypeDef) ToJson(ctx *jsonify.Context) jsonify.Datum {
	return jsonify.NewMap().
		Add(ctx, `name`, td.Name).
		Add(ctx, `type`, td.Type).
		AddNonZero(ctx, `methods`, td.Methods).
		AddNonZero(ctx, `typeParams`, td.TypeParams).
		AddNonZero(ctx, `interface`, td.Interface)
}

func (td *TypeDef) String() string {
	return jsonify.ToString(td)
}
