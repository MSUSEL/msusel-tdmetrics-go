package typeDesc

import (
	"fmt"
	"go/types"

	"github.com/Snow-Gremlin/goToolbox/collections"

	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/jsonify"
)

// uniqueName returns a unique name that isn't in the set.
// The new unique name will be added to the set.
// This is for naming anonymous fields and unnamed return values.
func uniqueName(names collections.Set[string]) string {
	const (
		attempts = 10_000
		pattern  = `$value%d`
	)
	for offset := 1; offset < attempts; offset++ {
		name := fmt.Sprintf(pattern, offset)
		if !names.Contains(name) {
			names.Add(name)
			return name
		}
	}
	panic(fmt.Errorf(`unable to find unique name in %d attempts`, attempts))
}

type Named interface {
	TypeDesc

	Name() string
	Type() TypeDesc
	EnsureName(names collections.Set[string])
}

type namedImp struct {
	name string
	typ  TypeDesc
}

func NewNamed(name string, typ TypeDesc) Named {
	return &namedImp{
		name: name,
		typ:  typ,
	}
}

func (t *namedImp) Name() string {
	return t.name
}

func (t *namedImp) Type() TypeDesc {
	return t.typ
}

func (t *namedImp) EnsureName(names collections.Set[string]) {
	if len(t.name) <= 0 || t.name == `_` {
		t.name = uniqueName(names)
	}
}

func (t *namedImp) SetIndex(index int) {
	// TODO: add index
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
		return jsonify.New(ctx, t.name)
	}

	return jsonify.NewMap().
		AddIf(ctx, ctx.IsKindShown(), `kind`, `named`).
		Add(ctx, `name`, t.name).
		Add(ctx.ShowKind().Short(), `type`, t.typ)
}

func (t *namedImp) String() string {
	return jsonify.ToString(t)
}
