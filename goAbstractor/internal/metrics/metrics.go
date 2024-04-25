package metrics

import (
	"go/ast"
	"go/token"

	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/jsonify"
)

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
