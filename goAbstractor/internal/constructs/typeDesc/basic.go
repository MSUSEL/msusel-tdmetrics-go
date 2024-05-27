package typeDesc

import (
	"go/types"

	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/jsonify"
	"github.com/Snow-Gremlin/goToolbox/utils"
)

type Basic struct {
	typ *types.Basic

	Name string
}

func NewBasic(typ *types.Basic) *Basic {
	return &Basic{
		typ:  typ,
		Name: typ.Name(),
	}
}

func BasicFor[T comparable]() *Basic {
	return &Basic{
		typ:  nil,
		Name: utils.TypeOf[T]().Name(),
	}
}

func (t *Basic) SetIndex(index int) {
	// TODO: add index
}

func (t *Basic) GoType() types.Type {
	return t.typ
}

func (t *Basic) Equal(other TypeDesc) bool {
	return equalTest(t, other, func(a, b *Basic) bool {
		return a.Name == b.Name
	})
}

func (t *Basic) ToJson(ctx *jsonify.Context) jsonify.Datum {
	return jsonify.New(ctx, t.Name)
}

func (t *Basic) String() string {
	return jsonify.ToString(t)
}
