package typeDesc

import (
	"errors"
	"go/types"

	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/jsonify"
)

type TypeDefRef interface {
	TypeDesc
	_typeDefRef()

	Name() string
	SetType(typ TypeDesc)
}

func NewTypeDefRef(reg Register, name string) TypeDefRef {
	return reg.RegisterTypeDefRef(&typeDefRefImp{
		name: name,
	})
}

type typeDefRefImp struct {
	name string
	typ  TypeDesc
}

func (t *typeDefRefImp) _typeDefRef() {}

func (t *typeDefRefImp) SetIndex(index int) {
	panic(errors.New(`do not call SetIndex on TypeDefRef`))
}

func (t *typeDefRefImp) GoType() types.Type {
	return t.typ.GoType()
}

func (t *typeDefRefImp) Equal(other TypeDesc) bool {
	return equalTest(t, other, func(a, b *typeDefRefImp) bool {
		return a.name == b.name
	})
}

func (t *typeDefRefImp) ToJson(ctx *jsonify.Context) jsonify.Datum {
	if ctx.IsReferenceShown() {
		return jsonify.NewMap().
			AddIf(ctx, ctx.IsKindShown(), `kind`, `ref`).
			Add(ctx, `name`, t.name).
			Add(ctx, `type`, t.typ)
	}

	return jsonify.New(ctx, t.typ)
}

func (t *typeDefRefImp) String() string {
	return jsonify.ToString(t)
}

func (t *typeDefRefImp) Name() string {
	return t.name
}

func (t *typeDefRefImp) SetType(typ TypeDesc) {
	t.typ = typ
}
