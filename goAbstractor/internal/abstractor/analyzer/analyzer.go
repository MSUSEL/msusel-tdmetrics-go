package analyzer

import (
	"go/ast"
	"go/token"
	"math"

	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/locs"
	"github.com/Snow-Gremlin/goToolbox/utils"
)

type Analyzer interface {
	Analyze(node ast.Node) Analyzer
	GetMetrics() constructs.MetricsArgs
}

type analyzerImp struct {
	locs       locs.Set
	loc        locs.Loc
	complexity int
	minLine    int
	maxLine    int
	indents    int
	minColumn  map[int]int
}

func New(locs locs.Set) Analyzer {
	return &analyzerImp{
		locs:       locs,
		loc:        nil,
		complexity: 1,
		maxLine:    0,
		minLine:    math.MaxInt,
		indents:    0,
		minColumn:  map[int]int{},
	}
}

func (a *analyzerImp) Analyze(node ast.Node) Analyzer {
	if utils.IsNil(a.loc) {
		a.loc = a.locs.NewLoc(node.Pos())
	}
	ast.Inspect(node, a.addCodePosForNode)
	return a
}

func (a *analyzerImp) GetMetrics() constructs.MetricsArgs {
	return constructs.MetricsArgs{
		Location:   a.loc,
		Complexity: a.complexity,
		LineCount:  a.maxLine - a.minLine + 1,
		CodeCount:  len(a.minColumn),
		Indents:    a.calcIndents(),
	}
}

func (a *analyzerImp) calcIndents() int {
	leftMostColumn := 10_000 // random large number
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
