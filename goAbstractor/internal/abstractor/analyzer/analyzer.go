package analyzer

import (
	"go/ast"
	"go/types"

	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/abstractor/analyzer/accessor"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/abstractor/analyzer/complexity"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/abstractor/analyzer/usages"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/abstractor/converter"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs"
)

// TODO: Add analytics:
//   - The set of variables with locations that are read from and written
//     to in each method. Used in Tight Class Cohesion (TCC) and
//     Design Recovery (DR).
//   - The set of all methods called in each method. Used for
//     Access to Foreign Data (ATFD) and Design Recovery (DR)
//   - Indicate if a method is an accessor getter or setter (single expression).

func Analyze(info *types.Info, proj constructs.Project, conv converter.Converter, node ast.Node) constructs.Metrics {
	var (
		loc    = proj.Locs().NewLoc(node.Pos())
		cmplx  = complexity.Calculate(node, proj.Locs().FileSet())
		acc    = accessor.Calculate(info, node)
		usages = usages.Calculate(info, proj, conv, node)
	)

	return proj.NewMetrics(constructs.MetricsArgs{
		Location:   loc,
		Complexity: cmplx.Complexity,
		LineCount:  cmplx.LineCount,
		CodeCount:  cmplx.CodeCount,
		Indents:    cmplx.Indents,
		Getter:     acc.Getter,
		Setter:     acc.Setter,
		Reads:      usages.Reads,
		Writes:     usages.Writes,
		Invokes:    usages.Invokes,
	})
}
