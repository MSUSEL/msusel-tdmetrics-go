package metrics

import (
	"cmp"

	"github.com/Snow-Gremlin/goToolbox/collections"
	"github.com/Snow-Gremlin/goToolbox/comp"

	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/assert"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs/kind"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/jsonify"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/locs"
)

// TODO: Add analytics:
//   - The set of variables with locations that are read from and written
//     to in each method. Used in Tight Class Cohesion (TCC) and
//     Design Recovery (DR).
//   - The set of all methods called in each method. Used for
//     Access to Foreign Data (ATFD) and Design Recovery (DR)
//   - Indicate if a method is an accessor getter or setter (single expression).

type metricsImp struct {
	loc   locs.Loc
	index int
	alive bool

	complexity int
	lineCount  int
	codeCount  int
	indents    int

	reads   collections.SortedSet[constructs.Usage]
	writes  collections.SortedSet[constructs.Usage]
	invokes collections.SortedSet[constructs.Usage]
	defines collections.SortedSet[constructs.Usage]
}

func newMetrics(args constructs.MetricsArgs) constructs.Metrics {
	assert.ArgNotNil(`location`, args.Location)

	return &metricsImp{
		loc: args.Location,

		complexity: args.Complexity,
		lineCount:  args.LineCount,
		codeCount:  args.CodeCount,
		indents:    args.Indents,

		reads:   args.Reads.Clone(),
		writes:  args.Writes.Clone(),
		invokes: args.Invokes.Clone(),
		defines: args.Defines.Clone(),
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

func (m *metricsImp) Reads() collections.ReadonlySortedSet[constructs.Usage] {
	return m.reads.Readonly()
}

func (m *metricsImp) Writes() collections.ReadonlySortedSet[constructs.Usage] {
	return m.writes.Readonly()
}

func (m *metricsImp) Invokes() collections.ReadonlySortedSet[constructs.Usage] {
	return m.invokes.Readonly()
}

func (m *metricsImp) Defines() collections.ReadonlySortedSet[constructs.Usage] {
	return m.defines.Readonly()
}

func (m *metricsImp) CompareTo(other constructs.Construct) int {
	return constructs.CompareTo[constructs.Metrics](m, other, Comparer())
}

func Comparer() comp.Comparer[constructs.Metrics] {
	return func(a, b constructs.Metrics) int {
		aImp, bImp := a.(*metricsImp), b.(*metricsImp)
		return cmp.Compare(int(aImp.loc.Pos()), int(bImp.loc.Pos()))
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
		AddNonZero(ctx, `loc`, m.loc).
		AddNonZero(ctx, `complexity`, m.complexity).
		AddNonZero(ctx, `lineCount`, m.lineCount).
		AddNonZero(ctx, `codeCount`, m.codeCount).
		AddNonZero(ctx, `indents`, m.indents).
		AddNonZero(ctx, `reads`, m.reads.ToSlice()).
		AddNonZero(ctx, `writes`, m.writes.ToSlice()).
		AddNonZero(ctx, `invokes`, m.invokes.ToSlice()).
		AddNonZero(ctx, `defines`, m.defines.ToSlice())
}

func (m *metricsImp) String() string {
	return jsonify.ToString(m)
}
