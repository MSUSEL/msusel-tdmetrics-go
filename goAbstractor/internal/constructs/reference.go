package constructs

import (
	"errors"
	"fmt"
	"go/types"

	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/jsonify"
	"github.com/Snow-Gremlin/goToolbox/utils"
)

type Reference interface {
	TypeDesc
	_reference()

	PackagePath() string
	Name() string
	SetType(pkg Package, typ TypeDef)
}

func NewReference(reg Types, realType *types.Named, pkgPath, name string) Reference {
	if utils.IsNil(realType) {
		panic(fmt.Errorf(`must provide a real type for %s.%s`, pkgPath, name))
	}

	return reg.RegisterTypeDefRef(&referenceImp{
		realType: realType,
		pkgPath:  pkgPath,
		name:     name,
	})
}

type referenceImp struct {
	realType *types.Named
	pkgPath  string
	name     string
	pkg      Package
	typ      TypeDef
}

func (t *referenceImp) _reference() {}

func (t *referenceImp) Visit(v Visitor) {}

func (t *referenceImp) SetIndex(index int) {
	panic(errors.New(`do not call SetIndex on TypeDefRef`))
}

func (t *referenceImp) GoType() types.Type {
	return t.realType
}

func (t *referenceImp) Equal(other TypeDesc) bool {
	return equalTest(t, other, func(a, b *referenceImp) bool {
		return a.pkgPath == b.pkgPath && a.name == b.name
	})
}

func (t *referenceImp) ToJson(ctx *jsonify.Context) jsonify.Datum {
	ctx2 := ctx.HideKind().Short()
	if ctx.IsReferenceShown() {
		return jsonify.NewMap().
			AddIf(ctx, ctx.IsKindShown(), `kind`, `ref`).
			AddNonZero(ctx2, `packagePath`, t.pkgPath).
			Add(ctx2, `name`, t.name).
			AddNonZero(ctx2, `package`, t.pkg).
			Add(ctx2, `type`, t.typ)
	}

	return jsonify.New(ctx2, t.typ)
}

func (t *referenceImp) String() string {
	return jsonify.ToString(t)
}

func (t *referenceImp) PackagePath() string {
	return t.pkgPath
}

func (t *referenceImp) Name() string {
	return t.name
}

func (t *referenceImp) SetType(pkg Package, typ TypeDef) {
	t.pkg = pkg
	t.typ = typ
}
