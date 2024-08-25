package metrics

import (
	"cmp"

	"github.com/Snow-Gremlin/goToolbox/comp"

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
	loc locs.Loc
	id  any

	complexity int
	lineCount  int
	codeCount  int
	indents    int
}

func newMetrics(args constructs.MetricsArgs) constructs.Metrics {
	return &metricsImp{
		loc:        args.Location,
		complexity: args.Complexity,
		lineCount:  args.LineCount,
		codeCount:  args.CodeCount,
		indents:    args.Indents,
	}
}

func (m *metricsImp) IsMetrics() {}

func (m *metricsImp) Kind() kind.Kind    { return kind.Metrics }
func (m *metricsImp) Id() any            { return m.id }
func (m *metricsImp) SetId(id any)       { m.id = id }
func (m *metricsImp) Location() locs.Loc { return m.loc }

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
	if ctx.IsShort() {
		return jsonify.New(ctx, m.id)
	}

	ctx2 := ctx.HideKind().Short()
	return jsonify.NewMap().
		AddIf(ctx2, ctx.IsKindShown(), `kind`, m.Kind()).
		AddIf(ctx2, ctx.IsIdShown(), `id`, m.id).
		AddNonZero(ctx2, `complexity`, m.complexity).
		AddNonZero(ctx2, `lineCount`, m.lineCount).
		AddNonZero(ctx2, `codeCount`, m.codeCount).
		AddNonZero(ctx2, `indents`, m.indents)
}

func (m *metricsImp) String() string {
	return jsonify.ToString(m)
}
