package value

import (
	"github.com/Snow-Gremlin/goToolbox/comp"

	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/assert"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs/kind"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/jsonify"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/locs"
)

type valueImp struct {
	pkg     constructs.Package
	name    string
	loc     locs.Loc
	typ     constructs.TypeDesc
	isConst bool
	metrics constructs.Metrics
	index   int
}

func newValue(args constructs.ValueArgs) constructs.Value {
	assert.ArgNotNil(`package`, args.Package)
	assert.ArgValidId(`name`, args.Name)
	assert.ArgNotNil(`type`, args.Type)
	assert.ArgNotNil(`location`, args.Location)

	return &valueImp{
		pkg:     args.Package,
		name:    args.Name,
		loc:     args.Location,
		typ:     args.Type,
		isConst: args.Const,
		metrics: args.Metrics,
	}
}

func (v *valueImp) IsDeclaration() {}
func (v *valueImp) IsValue()       {}

func (v *valueImp) Kind() kind.Kind    { return kind.Value }
func (v *valueImp) Index() int         { return v.index }
func (v *valueImp) SetIndex(index int) { v.index = index }

func (v *valueImp) Name() string                { return v.name }
func (v *valueImp) Location() locs.Loc          { return v.loc }
func (v *valueImp) Package() constructs.Package { return v.pkg }

func (v *valueImp) Type() constructs.TypeDesc          { return v.typ }
func (v *valueImp) Const() bool                        { return v.isConst }
func (v *valueImp) Metrics() constructs.Metrics        { return v.metrics }
func (v *valueImp) TypeParams() []constructs.TypeParam { return nil }

func (v *valueImp) CompareTo(other constructs.Construct) int {
	return constructs.CompareTo[constructs.Value](v, other, Comparer())
}

func Comparer() comp.Comparer[constructs.Value] {
	return func(a, b constructs.Value) int {
		aImp, bImp := a.(*valueImp), b.(*valueImp)
		return comp.Or(
			constructs.ComparerPend(aImp.pkg, bImp.pkg),
			comp.DefaultPend(aImp.name, bImp.name),
			constructs.ComparerPend(aImp.typ, bImp.typ),
		)
	}
}

func (v *valueImp) RemoveTempReferences() {
	v.typ = constructs.ResolvedTempReference(v.typ)
}

func (v *valueImp) ToJson(ctx *jsonify.Context) jsonify.Datum {
	if ctx.IsOnlyIndex() {
		return jsonify.New(ctx, v.index)
	}
	if ctx.IsShort() {
		return jsonify.NewSprintf(`%s%d`, v.Kind(), v.index)
	}
	return jsonify.NewMap().
		AddIf(ctx, ctx.IsDebugKindIncluded(), `kind`, v.Kind()).
		AddIf(ctx, ctx.IsDebugIndexIncluded(), `index`, v.index).
		Add(ctx.OnlyIndex(), `package`, v.pkg).
		Add(ctx, `name`, v.name).
		Add(ctx.Short(), `type`, v.typ).
		AddNonZero(ctx, `loc`, v.loc).
		AddNonZero(ctx, `const`, v.isConst).
		AddNonZero(ctx.OnlyIndex(), `metrics`, v.metrics)
}

func (v *valueImp) String() string {
	return `var ` + v.name + ` ` + v.typ.String()
}
