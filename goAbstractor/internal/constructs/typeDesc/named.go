package typeDesc

import (
	"go/types"

	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/jsonify"
)

type Named struct {
	Name string
	Type TypeDesc
}

func NewNamed(name string, t TypeDesc) *Named {
	return &Named{
		Name: name,
		Type: t,
	}
}

func (t *Named) GoType() types.Type {
	return t.Type.GoType()
}

func (t *Named) Equal(other TypeDesc) bool {
	return equalTest(t, other, func(a, b *Named) bool {
		return a.Name == b.Name &&
			equal(a.Type, b.Type)
	})
}

func (t *Named) ToJson(ctx *jsonify.Context) jsonify.Datum {
	if ctx.IsShort() {
		return jsonify.New(ctx, t.Name)
	}

	return jsonify.NewMap().
		AddIf(ctx, ctx.IsKindShown(), `kind`, `named`).
		Add(ctx, `name`, t.Name).
		Add(ctx.ShowKind().Short(), `type`, t.Type)
}

func (t *Named) String() string {
	return jsonify.ToString(t)
}
