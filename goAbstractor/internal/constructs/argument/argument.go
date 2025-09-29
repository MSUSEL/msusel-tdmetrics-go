package argument

import (
	"github.com/Snow-Gremlin/goToolbox/comp"

	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/assert"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs/kind"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/jsonify"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/stringer"
)

type argumentImp struct {
	constructs.ConstructCore
	name string
	typ  constructs.TypeDesc
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

func (a *argumentImp) Kind() kind.Kind           { return kind.Argument }
func (a *argumentImp) Name() string              { return a.name }
func (a *argumentImp) Type() constructs.TypeDesc { return a.typ }

func (a *argumentImp) ReplaceDuplicate(m map[constructs.Construct]constructs.Construct) {
	constructs.FindReplacement(m, &a.typ)
}

func (a *argumentImp) CompareTo(other constructs.Construct) int {
	return constructs.CompareTo[constructs.Argument](a, other, Comparer())
}

func Comparer() comp.Comparer[constructs.Argument] {
	return func(a, b constructs.Argument) int {
		aImp, bImp := a.(*argumentImp), b.(*argumentImp)
		if aImp == bImp {
			return 0
		}
		return comp.Or(
			comp.DefaultPend(aImp.name, bImp.name),
			constructs.ComparerPend(aImp.typ, bImp.typ),
		)
	}
}

func (a *argumentImp) RemoveTempReferences(required bool) (changed bool) {
	a.typ, changed = constructs.ResolvedTempReference(a.typ, required)
	return
}

func (a *argumentImp) ToJson(ctx *jsonify.Context) jsonify.Datum {
	if ctx.IsOnlyIndex() {
		return jsonify.New(ctx, a.Index())
	}
	if ctx.IsShort() {
		return jsonify.NewSprintf(`%s%d`, a.Kind(), a.Index())
	}
	if ctx.SkipDead() && !a.Alive() {
		return nil
	}
	if !ctx.KeepDuplicates() && a.Duplicate() {
		return nil
	}
	return jsonify.NewMap().
		AddIf(ctx, ctx.IsDebugKindIncluded(), `kind`, a.Kind()).
		AddIf(ctx, ctx.IsDebugIndexIncluded(), `index`, a.Index()).
		AddIf(ctx, ctx.IsDebugAliveIncluded(), `alive`, a.Alive()).
		AddNonZero(ctx, `name`, a.name).
		Add(ctx.Short(), `type`, a.typ)
}

func (a *argumentImp) ToStringer(s stringer.Stringer) {
	if len(a.name) > 0 {
		s.Write(a.name, ` `)
	}
	s.Write(a.typ)
}

func (a *argumentImp) String() string {
	return stringer.String(a)
}
