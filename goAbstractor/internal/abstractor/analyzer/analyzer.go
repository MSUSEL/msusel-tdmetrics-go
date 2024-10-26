package analyzer

import (
	"go/ast"
	"go/types"

	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/abstractor/analyzer/accessor"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/abstractor/analyzer/complexity"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/abstractor/analyzer/usages"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/abstractor/baker"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/abstractor/converter"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/assert"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs"
)

func Analyze(info *types.Info, proj constructs.Project, curPkg constructs.Package, baker baker.Baker, conv converter.Converter, node ast.Node) constructs.Metrics {
	assert.ArgNotNil(`curPkg`, curPkg)

	switch curPkg.Path() {
	case `runtime`, `unsafe`, `reflect`:
		// Skip analyzing things in these packages.
		return nil
	}

	var (
		loc    = proj.Locs().NewLoc(node.Pos())
		cmplx  = complexity.Calculate(node, proj.Locs().FileSet())
		acc    = accessor.Calculate(info, node)
		usages = usages.Calculate(info, proj, curPkg, baker, conv, node)
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
		SideEffect: usages.SideEffect,
	})
}
