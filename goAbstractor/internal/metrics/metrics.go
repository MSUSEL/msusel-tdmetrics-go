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
		Add(ctx, `complexity`, m.Complexity).
		Add(ctx, `lineCount`, m.LineCount).
		Add(ctx, `codeCount`, m.CodeCount).
		Add(ctx, `indents`, m.Indents)
}

func (m Metrics) String() string {
	return jsonify.ToString(m)
}

func New(fSet *token.FileSet, node ast.Node) Metrics {
	mc := newMetricsCalc(fSet, node)
	mc.calculateMetrics()
	return mc.getMetrics()
}
