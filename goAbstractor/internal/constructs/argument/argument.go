package argument

import (
	"github.com/Snow-Gremlin/goToolbox/comp"

	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/assert"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs/kind"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/jsonify"
)

type argumentImp struct {
	name  string
	typ   constructs.TypeDesc
	index int
}

func newArgument(args constructs.ArgumentArgs) constructs.Argument {
	if len(args.Name) > 0 && args.Name != `_` {
		// Arguments may be blank (unnamed or named with underscore).
		assert.ArgValidId(`name`, args.Name)
	}
	assert.ArgNotNil(`type`, args.Type)

	return &argumentImp{
		name: args.Name,
		typ:  args.Type,
	}
}

func (a *argumentImp) IsArgument()        {}
func (a *argumentImp) Kind() kind.Kind    { return kind.Argument }
func (a *argumentImp) SetIndex(index int) { a.index = index }

func (a *argumentImp) Name() string              { return a.name }
func (a *argumentImp) Type() constructs.TypeDesc { return a.typ }

func (a *argumentImp) CompareTo(other constructs.Construct) int {
	return constructs.CompareTo[constructs.Argument](a, other, Comparer())
}

func Comparer() comp.Comparer[constructs.Argument] {
	return func(a, b constructs.Argument) int {
		aImp, bImp := a.(*argumentImp), b.(*argumentImp)
		return comp.Or(
			comp.DefaultPend(aImp.name, bImp.name),
			constructs.ComparerPend(aImp.typ, bImp.typ),
		)
	}
}

func (a *argumentImp) ToJson(ctx *jsonify.Context) jsonify.Datum {
	if ctx.IsShort() {
		return jsonify.New(ctx, a.index)
	}

	ctx2 := ctx.HideKind().Short()
	return jsonify.NewMap().
		AddIf(ctx, ctx.IsKindShown(), `kind`, a.Kind()).
		AddIf(ctx, ctx.IsIndexShown(), `index`, a.index).
		Add(ctx2, `name`, a.name).
		Add(ctx2, `type`, a.typ)
}
