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

// Metrics are measurements taken for a method.
type Metrics struct {

	// Complexity is the McCabe's Cyclomatic Complexity value for the method.
	Complexity int

	// LineCount is the number of lines in a method
	// including function definition line and close bracket line.
	LineCount int

	// CodeCount is the number of lines that have code on it
	// including function definition line and close bracket line
	// but excluding blank lines and comment lines.
	CodeCount int

	// Indent is the indent complexity count using the amount of whitespace
	// to the left of the any line with code on it.
	Indents int
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
