package constructs

import (
	"errors"
	"fmt"
	"go/token"
	"go/types"

	"github.com/Snow-Gremlin/goToolbox/utils"
	"golang.org/x/tools/go/packages"

	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/jsonify"
)

type Basic interface {
	TypeDesc
	_basic()
}

func NewBasic(reg Register, typ *types.Basic) Basic {
	if utils.IsNil(typ) {
		panic(errors.New(`may not create a new basic with a nil type`))
	}
	return reg.RegisterBasic(&basicImp{
		typ:  typ,
		name: typ.Name(),
	})
}

func BasicFromName(reg Register, pkg *packages.Package, typeName string) Basic {
	tv, err := types.Eval(pkg.Fset, pkg.Types, token.NoPos, `(*`+typeName+`)(nil)`)
	if err != nil {
		panic(fmt.Errorf(`unable to create basic type of %s: %w`, typeName, err))
	}
	typ := tv.Type.(*types.Pointer).Elem().(*types.Basic)
	return NewBasic(reg, typ)
}

type basicImp struct {
	typ   *types.Basic
	name  string
	index int
}

func (t *basicImp) _basic() {}

func (t *basicImp) Visit(v Visitor) {}

func (t *basicImp) SetIndex(index int) {
	t.index = index
}

func (t *basicImp) GoType() types.Type {
	return t.typ
}

func (t *basicImp) Equal(other TypeDesc) bool {
	return equalTest(t, other, func(a, b *basicImp) bool {
		return a.name == b.name
	})
}

func (t *basicImp) ToJson(ctx *jsonify.Context) jsonify.Datum {
	if ctx.IsShort() {
		return jsonify.New(ctx, t.index)
	}

	if ctx.IsKindShown() {
		return jsonify.NewMap().
			AddIf(ctx, ctx.IsKindShown(), `kind`, `basic`).
			Add(ctx, `name`, t.name)
	}

	return jsonify.New(ctx, t.name)
}

func (t *basicImp) String() string {
	return jsonify.ToString(t)
}
