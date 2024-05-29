package typeDesc

import (
	"go/types"

	"github.com/Snow-Gremlin/goToolbox/utils"

	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/jsonify"
)

type Basic interface {
	TypeDesc
}

func NewBasic(typ *types.Basic) Basic {
	return &basicImp{
		typ:  typ,
		name: typ.Name(),
	}
}

func BasicFor[T comparable]() Basic {
	return &basicImp{
		typ:  nil,
		name: utils.TypeOf[T]().Name(),
	}
}

type basicImp struct {
	typ   *types.Basic
	name  string
	index int
}

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

	return jsonify.New(ctx, t.name)
}

func (t *basicImp) String() string {
	return jsonify.ToString(t)
}
