package typeDesc

import (
	"go/types"

	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/jsonify"
	"github.com/Snow-Gremlin/goToolbox/utils"
)

type Basic interface {
	TypeDesc
}

type basicImp struct {
	typ  *types.Basic
	name string
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

func (t *basicImp) SetIndex(index int) {
	// TODO: add index
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
	return jsonify.New(ctx, t.name)
}

func (t *basicImp) String() string {
	return jsonify.ToString(t)
}
