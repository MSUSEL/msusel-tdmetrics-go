package metrics

import (
	"go/ast"
	"go/token"

	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/jsonify"
)

// TODO: Add analytics:
//   - The set of variables with locations that are read from and written
//     to in each method. Used in Tight Class Cohesion (TCC) and
//     Design Recovery (DR).
//   - The set of all methods called in each method. Used for
//     Access to Foreign Data (ATFD) and Design Recovery (DR)
//   - Indicate if a method is an accessor getter or setter (single expression).

type Metrics struct {
	Complexity int
	LineCount  int
	CodeCount  int
	Indents    int
}

func (m Metrics) ToJson(ctx *jsonify.Context) jsonify.Datum {
	return jsonify.NewMap().
		AddNonZero(ctx, `complexity`, m.Complexity).
		AddNonZero(ctx, `lineCount`, m.LineCount).
		AddNonZero(ctx, `codeCount`, m.CodeCount).
		AddNonZero(ctx, `indents`, m.Indents)
}

func (m Metrics) String() string {
	return jsonify.ToString(m)
}

func New(fSet *token.FileSet, node ast.Node) Metrics {
	mc := newMetricsCalc(fSet, node)
	mc.calculateMetrics()
	return mc.getMetrics()
}
