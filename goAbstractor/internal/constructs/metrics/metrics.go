package metrics

import (
	"github.com/Snow-Gremlin/goToolbox/collections"
	"github.com/Snow-Gremlin/goToolbox/comp"

	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/assert"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs/kind"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/jsonify"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/locs"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/stringer"
)

type metricsImp struct {
	loc   locs.Loc
	index int
	alive bool

	complexity int
	lineCount  int
	codeCount  int
	indents    int
	getter     bool
	setter     bool
	sideEffect bool

	reads   collections.SortedSet[constructs.Construct]
	writes  collections.SortedSet[constructs.Construct]
	invokes collections.SortedSet[constructs.Construct]
}

func newMetrics(args constructs.MetricsArgs) constructs.Metrics {
	assert.ArgNotNil(`location`, args.Location)
	assert.ArgHasNoNils(`reads`, args.Reads.ToSlice())
	assert.ArgHasNoNils(`writes`, args.Writes.ToSlice())
	assert.ArgHasNoNils(`invokes`, args.Invokes.ToSlice())

	return &metricsImp{
		loc: args.Location,

		complexity: args.Complexity,
		lineCount:  args.LineCount,
		codeCount:  args.CodeCount,
		indents:    args.Indents,
		getter:     args.Getter,
		setter:     args.Setter,
		sideEffect: args.SideEffect,

		reads:   args.Reads,
		writes:  args.Writes,
		invokes: args.Invokes,
	}
}

func (m *metricsImp) IsMetrics() {}

func (m *metricsImp) Kind() kind.Kind     { return kind.Metrics }
func (m *metricsImp) Index() int          { return m.index }
func (m *metricsImp) SetIndex(index int)  { m.index = index }
func (m *metricsImp) Alive() bool         { return m.alive }
func (m *metricsImp) SetAlive(alive bool) { m.alive = alive }
func (m *metricsImp) Location() locs.Loc  { return m.loc }
func (m *metricsImp) Complexity() int     { return m.complexity }
func (m *metricsImp) LineCount() int      { return m.lineCount }
func (m *metricsImp) CodeCount() int      { return m.codeCount }
func (m *metricsImp) Indents() int        { return m.indents }
func (m *metricsImp) Getter() bool        { return m.getter }
func (m *metricsImp) Setter() bool        { return m.setter }
func (m *metricsImp) SideEffect() bool    { return m.sideEffect }

func (m *metricsImp) Reads() collections.ReadonlySortedSet[constructs.Construct] {
	return m.reads.Readonly()
}

func (m *metricsImp) Writes() collections.ReadonlySortedSet[constructs.Construct] {
	return m.writes.Readonly()
}

func (m *metricsImp) Invokes() collections.ReadonlySortedSet[constructs.Construct] {
	return m.invokes.Readonly()
}

func (m *metricsImp) RemoveTempDeclRefs() {
	constructs.ResolveTempDeclRefSet(m.reads)
	constructs.ResolveTempDeclRefSet(m.writes)
	constructs.ResolveTempDeclRefSet(m.invokes)
}

func (m *metricsImp) CompareTo(other constructs.Construct) int {
	return constructs.CompareTo[constructs.Metrics](m, other, Comparer())
}

func Comparer() comp.Comparer[constructs.Metrics] {
	return func(a, b constructs.Metrics) int {
		aImp, bImp := a.(*metricsImp), b.(*metricsImp)
		aOffset, aFile, aLine := aImp.loc.Info()
		bOffset, bFile, bLine := bImp.loc.Info()
		return comp.Or(
			comp.DefaultPend(aFile, bFile),
			comp.DefaultPend(aLine, bLine),
			comp.DefaultPend(aOffset, bOffset),
		)
	}
}

func (m *metricsImp) ToJson(ctx *jsonify.Context) jsonify.Datum {
	if ctx.IsOnlyIndex() {
		return jsonify.New(ctx, m.index)
	}
	if ctx.IsShort() {
		return jsonify.NewSprintf(`%s%d`, m.Kind(), m.index)
	}
	return jsonify.NewMap().
		AddIf(ctx, ctx.IsDebugKindIncluded(), `kind`, m.Kind()).
		AddIf(ctx, ctx.IsDebugIndexIncluded(), `index`, m.index).
		AddNonZero(ctx, `loc`, m.loc). // Should only be zero for unit-tests.
		AddNonZero(ctx, `complexity`, m.complexity).
		AddNonZero(ctx, `lineCount`, m.lineCount).
		AddNonZero(ctx, `codeCount`, m.codeCount).
		AddNonZero(ctx, `indents`, m.indents).
		AddNonZero(ctx, `getter`, m.getter).
		AddNonZero(ctx, `setter`, m.setter).
		AddNonZero(ctx.Short(), `reads`, m.reads.ToSlice()).
		AddNonZero(ctx.Short(), `writes`, m.writes.ToSlice()).
		AddNonZero(ctx.Short(), `invokes`, m.invokes.ToSlice()).
		AddNonZero(ctx, `sideEffect`, m.sideEffect)
}

func (m *metricsImp) ToStringer(s stringer.Stringer) {
	s.Write(m.String())
}

func (m *metricsImp) String() string {
	return jsonify.ToString(m)
}
