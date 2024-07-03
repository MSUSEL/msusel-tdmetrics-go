package metrics

import (
	"go/ast"
	"go/token"
	"math"
)

type metricsCalc struct {
	fSet       *token.FileSet
	node       ast.Node
	complexity int
	minLine    int
	maxLine    int
	indents    int
	minColumn  map[int]int
}

func newMetricsCalc(fSet *token.FileSet, node ast.Node) *metricsCalc {
	return &metricsCalc{
		fSet:       fSet,
		node:       node,
		complexity: 1,
		maxLine:    0,
		minLine:    math.MaxInt,
		indents:    0,
		minColumn:  map[int]int{},
	}
}

func (m *metricsCalc) calculateMetrics() {
	ast.Inspect(m.node, m.addCodePosForNode)
	m.finishIndents()
}

func (m *metricsCalc) finishIndents() {
	leftMostColumn := 10_000 // random large number
	indentSum := 0
	for _, ind := range m.minColumn {
		leftMostColumn = min(ind, leftMostColumn)
		indentSum += ind
	}
	m.indents = indentSum - len(m.minColumn)*leftMostColumn
}

func (m *metricsCalc) getMetrics() Metrics {
	return Metrics{
		Complexity: m.complexity,
		LineCount:  m.maxLine - m.minLine + 1,
		CodeCount:  len(m.minColumn),
		Indents:    m.indents,
	}
}

func (m *metricsCalc) incComplexity(check bool) {
	if check {
		m.complexity++
	}
}

func (m *metricsCalc) addCodePos(pos token.Pos, isEnd bool) {
	p := m.fSet.PositionFor(pos, false)
	lineNo, column := p.Line, p.Column
	m.maxLine = max(m.maxLine, lineNo)
	m.minLine = min(m.minLine, lineNo)
	if isEnd {
		column--
	}
	if otherCol, ok := m.minColumn[lineNo]; ok {
		column = min(column, otherCol)
	}
	m.minColumn[lineNo] = column
}

func (m *metricsCalc) addCodePosForNode(n ast.Node) bool {
	switch t := n.(type) {
	case nil, *ast.Comment, *ast.CommentGroup:
		return true
	case *ast.IfStmt, *ast.ForStmt, *ast.RangeStmt, *ast.GoStmt:
		m.incComplexity(true)
	case *ast.CaseClause:
		m.incComplexity(t.List != nil)
	case *ast.CommClause:
		m.incComplexity(t.Comm != nil)
	case *ast.BinaryExpr:
		m.incComplexity(t.Op == token.LAND || t.Op == token.LOR)
	}

	m.addCodePos(n.Pos(), false)
	if ended, has := n.(interface{ End() token.Pos }); has {
		m.addCodePos(ended.End(), true)
	}
	return true
}
