package constructs

import (
	"errors"
	"go/types"

	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/jsonify"
)

type TypeDefRef interface {
	TypeDesc
	_typeDefRef()

	PackagePath() string
	Name() string
	SetType(typ TypeDef)
}

func NewTypeDefRef(reg Register, pkgPath, name string) TypeDefRef {
	return reg.RegisterTypeDefRef(&typeDefRefImp{
		pkgPath: pkgPath,
		name:    name,
	})
}

type typeDefRefImp struct {
	pkgPath string
	name    string
	typ     TypeDef
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
		return a.pkgPath == b.pkgPath && a.name == b.name
	})
}

func (t *typeDefRefImp) ToJson(ctx *jsonify.Context) jsonify.Datum {
	if ctx.IsReferenceShown() {
		return jsonify.NewMap().
			AddIf(ctx, ctx.IsKindShown(), `kind`, `ref`).
			Add(ctx, `packagePath`, t.pkgPath).
			Add(ctx, `name`, t.name).
			Add(ctx, `type`, t.typ)
	}

	return jsonify.New(ctx, t.typ)
}

func (t *typeDefRefImp) String() string {
	return jsonify.ToString(t)
}

func (t *typeDefRefImp) PackagePath() string {
	return t.pkgPath
}

func (t *typeDefRefImp) Name() string {
	return t.name
}

func (t *typeDefRefImp) SetType(typ TypeDef) {
	t.typ = typ
}
