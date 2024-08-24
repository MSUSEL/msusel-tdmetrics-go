package metrics

import (
	"cmp"

	"github.com/Snow-Gremlin/goToolbox/comp"
	"github.com/Snow-Gremlin/goToolbox/utils"

	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs"
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
	loc        locs.Loc
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

func (m *metricsImp) IsMetrics()         {}
func (m *metricsImp) Location() locs.Loc { return m.loc }

func (m *metricsImp) CompareTo(other constructs.Metrics) int {
	if utils.IsNil(m) {
		return utils.Ternary(utils.IsNil(other), 0, -1)
	}
	if utils.IsNil(other) {
		return 1
	}
	return Comparer()(m, other)
}

func Comparer() comp.Comparer[constructs.Metrics] {
	return func(a, b constructs.Metrics) int {
		aImp, bImp := a.(*metricsImp), b.(*metricsImp)
		return cmp.Compare(int(aImp.loc.Pos()), int(bImp.loc.Pos()))
	}
}

func (m *metricsImp) ToJson(ctx *jsonify.Context) jsonify.Datum {
	return jsonify.NewMap().
		AddNonZero(ctx, `complexity`, m.complexity).
		AddNonZero(ctx, `lineCount`, m.lineCount).
		AddNonZero(ctx, `codeCount`, m.codeCount).
		AddNonZero(ctx, `indents`, m.indents)
}

func (m *metricsImp) String() string {
	return jsonify.ToString(m)
}
