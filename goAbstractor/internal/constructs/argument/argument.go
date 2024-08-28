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

func (a *argumentImp) IsArgument() {}

func (a *argumentImp) Kind() kind.Kind    { return kind.Argument }
func (a *argumentImp) Index() int         { return a.index }
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

func (a *argumentImp) RemoveTempReferences() {
	a.typ = constructs.ResolvedTempReference(a.typ)
}

func (a *argumentImp) ToJson(ctx *jsonify.Context) jsonify.Datum {
	if ctx.IsOnlyIndex() {
		return jsonify.New(ctx, a.index)
	}
	if ctx.IsShort() {
		return jsonify.NewSprintf(`%s%d`, a.Kind(), a.index)
	}
	return jsonify.NewMap().
		AddIf(ctx, ctx.IsDebugKindIncluded(), `kind`, a.Kind()).
		AddIf(ctx, ctx.IsDebugIndexIncluded(), `index`, a.index).
		AddNonZero(ctx, `name`, a.name).
		Add(ctx.Short(), `type`, a.typ)
}

func (a *argumentImp) String() string {
	if len(a.name) > 0 {
		return a.name + ` ` + a.typ.String()
	}
	return a.typ.String()
}
