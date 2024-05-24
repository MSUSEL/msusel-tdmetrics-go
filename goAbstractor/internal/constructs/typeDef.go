package constructs

import (
	"go/types"

	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs/typeDesc"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/jsonify"
)

type TypeDef struct {
	typ *types.Interface

	Name       string
	Type       typeDesc.TypeDesc
	Inherits   []*typeDesc.Interface
	Methods    []*Method
	TypeParams []*typeDesc.Named

	Interface *typeDesc.Interface
}

func NewTypeDef(name string, t typeDesc.TypeDesc) *TypeDef {
	return &TypeDef{
		Name: name,
		Type: t,
	}
}

func (td *TypeDef) SetInheritance(typ *types.Interface, in *typeDesc.Interface) {
	td.typ = typ
	td.Interface = in
}

func (td *TypeDef) ToJson(ctx *jsonify.Context) jsonify.Datum {
	return jsonify.NewMap().
		Add(ctx, `name`, td.Name).
		Add(ctx, `type`, td.Type).
		AddNonZero(ctx, `inherits`, td.Inherits).
		AddNonZero(ctx, `methods`, td.Methods).
		AddNonZero(ctx, `typeParams`, td.TypeParams)
}

func (td *TypeDef) String() string {
	return jsonify.ToString(td)
}
