package value

import (
	"github.com/Snow-Gremlin/goToolbox/comp"
	"github.com/Snow-Gremlin/goToolbox/utils"

	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/assert"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs/kind"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/jsonify"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/locs"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/stringer"
)

type valueImp struct {
	constructs.ConstructCore
	pkg      constructs.Package
	name     string
	exported bool
	loc      locs.Loc
	typ      constructs.TypeDesc
	isConst  bool
	metrics  constructs.Metrics
}

func newValue(args constructs.ValueArgs) constructs.Value {
	if args.Name != `_` {
		// Arguments may be a blank with an underscore.
		assert.ArgValidId(`name`, args.Name)
	}
	assert.ArgNotNil(`package`, args.Package)
	assert.ArgNotNil(`type`, args.Type)
	assert.ArgNotNil(`location`, args.Location)

	return &valueImp{
		pkg:      args.Package,
		name:     args.Name,
		exported: args.Exported,
		loc:      args.Location,
		typ:      args.Type,
		isConst:  args.Const,
		metrics:  args.Metrics,
	}
}

func (v *valueImp) IsDeclaration() {}
func (v *valueImp) IsValue()       {}

func (v *valueImp) Kind() kind.Kind    { return kind.Value }
func (v *valueImp) Const() bool        { return v.isConst }
func (v *valueImp) Name() string       { return v.name }
func (v *valueImp) Exported() bool     { return v.exported }
func (v *valueImp) Location() locs.Loc { return v.loc }

func (v *valueImp) Package() constructs.Package        { return v.pkg }
func (v *valueImp) Type() constructs.TypeDesc          { return v.typ }
func (v *valueImp) Metrics() constructs.Metrics        { return v.metrics }
func (v *valueImp) TypeParams() []constructs.TypeParam { return nil }

func (v *valueImp) HasSideEffect() bool {
	return !utils.IsNil(v.metrics) && v.metrics.SideEffect()
}

func (v *valueImp) CompareTo(other constructs.Construct) int {
	return constructs.CompareTo[constructs.Value](v, other, Comparer())
}

func Comparer() comp.Comparer[constructs.Value] {
	return func(a, b constructs.Value) int {
		aImp, bImp := a.(*valueImp), b.(*valueImp)
		if aImp == bImp {
			return 0
		}
		return comp.Or(
			constructs.ComparerPend(aImp.pkg, bImp.pkg),
			comp.DefaultPend(aImp.name, bImp.name),
			constructs.ComparerPend(aImp.typ, bImp.typ),
			constructs.ComparerPend(aImp.metrics, bImp.metrics),
		)
	}
}

func (v *valueImp) RemoveTempReferences(required bool) (changed bool) {
	v.typ, changed = constructs.ResolvedTempReference(v.typ, required)
	return
}

func (v *valueImp) ToJson(ctx *jsonify.Context) jsonify.Datum {
	if ctx.IsOnlyIndex() {
		return jsonify.New(ctx, v.Index())
	}
	if ctx.IsShort() {
		return jsonify.NewSprintf(`%s%d`, v.Kind(), v.Index())
	}
	if ctx.SkipDead() && !v.Alive() {
		return nil
	}
	if !ctx.KeepDuplicates() && v.Duplicate() {
		return nil
	}
	return jsonify.NewMap().
		AddIf(ctx, ctx.IsDebugKindIncluded(), `kind`, v.Kind()).
		AddIf(ctx, ctx.IsDebugIndexIncluded(), `index`, v.Index()).
		AddIf(ctx, ctx.IsDebugAliveIncluded(), `alive`, v.Alive()).
		Add(ctx.OnlyIndex(), `package`, v.pkg).
		Add(ctx, `name`, v.name).
		Add(ctx.Short(), `type`, v.typ).
		AddNonZero(ctx, `loc`, v.loc).
		AddNonZero(ctx, `const`, v.isConst).
		AddNonZeroIf(ctx, v.exported, `vis`, `exported`).
		AddNonZero(ctx.OnlyIndex(), `metrics`, v.metrics)
}

func (v *valueImp) ToStringer(s stringer.Stringer) {
	s.Write(`var `, v.pkg.Path(), `.`, v.name, ` `, v.typ)
}

func (v *valueImp) String() string {
	return stringer.String(v)
}
