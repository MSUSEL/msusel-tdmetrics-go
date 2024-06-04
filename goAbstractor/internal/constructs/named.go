package constructs

import (
	"fmt"
	"go/types"

	"github.com/Snow-Gremlin/goToolbox/utils"

	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/jsonify"
)

type Named interface {
	TypeDesc
	_named()

	Name() string
	Type() TypeDesc
}

func NewNamed(reg Register, name string, typ TypeDesc) Named {
	if utils.IsNil(typ) {
		panic(fmt.Errorf(`must have a non-nil type for named, %q`, name))
	}
	return reg.RegisterNamed(&namedImp{
		name: name,
		typ:  typ,
	})
}

type namedImp struct {
	name  string
	typ   TypeDesc
	index int
}

func (t *namedImp) _named() {}

func (t *namedImp) SetIndex(index int) {
	t.index = index
}

func (t *namedImp) GoType() types.Type {
	return t.typ.GoType()
}

func (t *namedImp) Equal(other TypeDesc) bool {
	return equalTest(t, other, func(a, b *namedImp) bool {
		return a.name == b.name &&
			equal(a.typ, b.typ)
	})
}

func (t *namedImp) ToJson(ctx *jsonify.Context) jsonify.Datum {
	if ctx.IsShort() {
		return jsonify.New(ctx, t.index)
	}

	ctx2 := ctx.HideKind().Short()
	return jsonify.NewMap().
		AddIf(ctx, ctx.IsKindShown(), `kind`, `named`).
		Add(ctx2, `name`, t.name).
		Add(ctx2, `type`, t.typ)
}

func (t *namedImp) String() string {
	return jsonify.ToString(t)
}

func (t *namedImp) Name() string {
	return t.name
}

func (t *namedImp) Type() TypeDesc {
	return t.typ
}
