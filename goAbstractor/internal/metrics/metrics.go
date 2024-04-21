package metrics

import (
	"fmt"
	"go/ast"
	"go/token"
	"math"
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
	fSet    *token.FileSet
	node    ast.Node
	minLine int
	maxLine int
	data    map[int]int
}

func New(fSet *token.FileSet, node ast.Node) Metrics {
	met := Metrics{}
	m := &metricsCalc{
		Metrics: &met,
		fSet:    fSet,
		node:    node,
		maxLine: 0,
		minLine: math.MaxInt,
		data:    map[int]int{},
	}

	m.Complexity = 1
	ast.Inspect(m.node, m.addCodePosForNode)

	fmt.Println(m.data) // TODO: REMOVE

	m.LineCount = m.maxLine - m.minLine + 1
	m.CodeCount = len(m.data)
	m.Indents = 0
	for _, indent := range m.data {
		m.Indents += indent - 1
	}

	return met
}

func (m *metricsCalc) incComplexity(check bool) {
	if check {
		m.Complexity++
	}
}

func (m *metricsCalc) linePos(pos token.Pos) (int, int) {
	p := m.fSet.PositionFor(pos, false)
	return p.Line, p.Column
}

func (m *metricsCalc) addCodePos(p token.Pos) {
	lineNo, column := m.linePos(p)
	m.maxLine = max(m.maxLine, lineNo)
	m.minLine = min(m.minLine, lineNo)
	if otherCol, ok := m.data[lineNo]; ok {
		m.data[lineNo] = min(column, otherCol)
	} else {
		m.data[lineNo] = column
	}
}

func (m *metricsCalc) addCodePosForNode(n ast.Node) bool {
	switch t := n.(type) {
	case nil, *ast.Comment, *ast.CommentGroup:
		return true
	case *ast.IfStmt, *ast.ForStmt, *ast.RangeStmt:
		m.incComplexity(true)
	case *ast.CaseClause:
		m.incComplexity(t.List != nil)
	case *ast.CommClause:
		m.incComplexity(t.Comm != nil)
	case *ast.BinaryExpr:
		m.incComplexity(t.Op == token.LAND || t.Op == token.LOR)
	}

	m.addCodePos(n.Pos())
	if ended, has := n.(interface{ End() token.Pos }); has {
		m.addCodePos(ended.End())
	}
	return true
}
