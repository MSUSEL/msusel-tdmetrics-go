package analyzer

import (
	"go/ast"
	"go/token"
	"math"

	"github.com/Snow-Gremlin/goToolbox/collections"
	"github.com/Snow-Gremlin/goToolbox/collections/sortedSet"
	"github.com/Snow-Gremlin/goToolbox/utils"

	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs/usage"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/locs"
)

// TODO: Add analytics:
//   - The set of variables with locations that are read from and written
//     to in each method. Used in Tight Class Cohesion (TCC) and
//     Design Recovery (DR).
//   - The set of all methods called in each method. Used for
//     Access to Foreign Data (ATFD) and Design Recovery (DR)
//   - Indicate if a method is an accessor getter or setter (single expression).

func Analyze(locs locs.Set, factory constructs.MetricsFactory, node ast.Node) constructs.Metrics {
	return factory.NewMetrics(newAnalyzer(locs).Analyze(node).GetMetricsArgs())
}

type analyzerImp struct {
	locs locs.Set
	loc  locs.Loc

	complexity int
	minLine    int
	maxLine    int
	indents    int
	minColumn  map[int]int
	getter     bool
	setter     bool

	reads   collections.SortedSet[constructs.Usage]
	writes  collections.SortedSet[constructs.Usage]
	invokes collections.SortedSet[constructs.Usage]
	defines collections.SortedSet[constructs.Usage]
}

func newAnalyzer(locs locs.Set) *analyzerImp {
	return &analyzerImp{
		locs: locs,
		loc:  nil,

		complexity: 1,
		maxLine:    0,
		minLine:    math.MaxInt,
		indents:    0,
		minColumn:  map[int]int{},

		reads:   sortedSet.New(usage.Comparer()),
		writes:  sortedSet.New(usage.Comparer()),
		invokes: sortedSet.New(usage.Comparer()),
		defines: sortedSet.New(usage.Comparer()),
	}
}

func (a *analyzerImp) Analyze(node ast.Node) *analyzerImp {
	if utils.IsNil(a.loc) {
		a.loc = a.locs.NewLoc(node.Pos())
	}
	// gather positional information for indents and cyclomatic complexity.
	ast.Inspect(node, a.addCodePosForNode)
	a.checkForGetter(node)
	a.checkForSetter(node)
	return a
}

func (a *analyzerImp) GetMetricsArgs() constructs.MetricsArgs {
	return constructs.MetricsArgs{
		Location:   a.loc,
		Complexity: a.complexity,
		LineCount:  a.maxLine - a.minLine + 1,
		CodeCount:  len(a.minColumn),
		Indents:    a.calcIndents(),
		Reads:      a.reads,
		Writes:     a.writes,
		Invokes:    a.invokes,
		Defines:    a.defines,
	}
}

func (a *analyzerImp) calcIndents() int {
	leftMostColumn := math.MaxInt
	indentSum := 0
	for _, ind := range a.minColumn {
		leftMostColumn = min(ind, leftMostColumn)
		indentSum += ind
	}
	return indentSum - len(a.minColumn)*leftMostColumn
}

func (a *analyzerImp) incComplexity(check bool) {
	if check {
		a.complexity++
	}
}

func (a *analyzerImp) addCodePos(pos token.Pos, isEnd bool) {
	p := a.locs.FileSet().PositionFor(pos, false)
	lineNo, column := p.Line, p.Column
	a.maxLine = max(a.maxLine, lineNo)
	a.minLine = min(a.minLine, lineNo)
	if isEnd {
		column--
	}
	if otherCol, ok := a.minColumn[lineNo]; ok {
		column = min(column, otherCol)
	}
	a.minColumn[lineNo] = column
}

func (a *analyzerImp) addCodePosForNode(n ast.Node) bool {
	switch t := n.(type) {
	case nil, *ast.Comment, *ast.CommentGroup:
		return true
	case *ast.IfStmt, *ast.ForStmt, *ast.RangeStmt, *ast.GoStmt:
		a.incComplexity(true)
	case *ast.CaseClause:
		a.incComplexity(t.List != nil)
	case *ast.CommClause:
		a.incComplexity(t.Comm != nil)
	case *ast.BinaryExpr:
		a.incComplexity(t.Op == token.LAND || t.Op == token.LOR)
	}

	a.addCodePos(n.Pos(), false)
	if ended, has := n.(interface{ End() token.Pos }); has {
		a.addCodePos(ended.End(), true)
	}
	return true
}

func (a *analyzerImp) checkForGetter(n ast.Node) {
	ast.Print(a.locs.FileSet(), n)
}

func (a *analyzerImp) checkForSetter(n ast.Node) {

}
