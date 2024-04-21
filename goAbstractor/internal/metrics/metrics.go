package metrics

import (
	"fmt"
	"go/ast"
	"go/token"
)

type Metrics struct {
	Complexity int
	LineCount  int
	CodeCount  int
	Indents    int
}

func (m Metrics) String() string {
	return fmt.Sprintf(
		"Complexity: %d\nLineCount:  %d\nCodeCount:  %d\nIndents:    %d",
		m.Complexity, m.LineCount, m.CodeCount, m.Indents)
}

type metricsCalc struct {
	*Metrics
	fSet  *token.FileSet
	block *ast.BlockStmt
	data  map[int]int
}

func New(fSet *token.FileSet, block *ast.BlockStmt) Metrics {
	met := Metrics{}
	m := &metricsCalc{
		Metrics: &met,
		fSet:    fSet,
		block:   block,
	}
	m.calculateComplexity()
	m.calculateLineCount()
	m.calculateCodeCounts()
	return met
}

func (m *metricsCalc) incComplexity(check bool) {
	if check {
		m.Complexity++
	}
}

func (m *metricsCalc) complexityNode(n ast.Node) bool {
	switch n := n.(type) {
	case *ast.IfStmt, *ast.ForStmt, *ast.RangeStmt:
		m.incComplexity(true)
	case *ast.CaseClause:
		m.incComplexity(n.List != nil)
	case *ast.CommClause:
		m.incComplexity(n.Comm != nil)
	case *ast.BinaryExpr:
		m.incComplexity(n.Op == token.LAND || n.Op == token.LOR)
	}
	return true
}

func (m *metricsCalc) calculateComplexity() {
	m.Complexity = 1
	ast.Inspect(m.block, m.complexityNode)
}

func (m *metricsCalc) linePos(pos token.Pos) (int, int) {
	p := m.fSet.PositionFor(pos, false)
	return p.Line, p.Column
}

func (m *metricsCalc) calculateLineCount() {
	first, _ := m.linePos(m.block.Lbrace)
	last, _ := m.linePos(m.block.Rbrace)
	m.LineCount = last - first + 1
}

func (m *metricsCalc) addCodePos(p token.Pos) {
	lineNo, column := m.linePos(p)
	if otherCol, ok := m.data[lineNo]; ok {
		m.data[lineNo] = min(column, otherCol)
	} else {
		m.data[lineNo] = column
	}
}

func (m *metricsCalc) addCodePosForNode(n ast.Node) bool {
	switch n.(type) {
	case nil, *ast.Comment, *ast.CommentGroup:
		return true
	}
	m.addCodePos(n.Pos())
	if ended, has := n.(interface{ End() token.Pos }); has {
		m.addCodePos(ended.End())
	}
	return true
}

func (m *metricsCalc) calculateCodeCounts() {
	// See https://codescene.com/engineering-blog/bumpy-road-code-complexity-in-context/
	// See https://codescene.io/docs/guides/technical/complexity-trends.html
	m.data = map[int]int{}
	for _, item := range m.block.List {
		ast.Inspect(item, m.addCodePosForNode)
	}
	fmt.Println(m.data)
	m.CodeCount = len(m.data)
	m.Indents = 0
	for _, indent := range m.data {
		m.Indents += indent - 1
	}
}
