package metrics

import (
	"go/ast"
	"go/token"
)

type Metrics struct {
	Complexity int
	LineCount  int
	CodeCount  int
}

func New(b *ast.BlockStmt) *Metrics {
	m := &Metrics{}
	m.calculateComplexity(b)
	m.calculateLineCounts(b)
	return m
}

func (m *Metrics) incComplexity(check bool) {
	if check {
		m.Complexity++
	}
}

func (m *Metrics) complexityNode(n ast.Node) bool {
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

func (m *Metrics) calculateComplexity(b *ast.BlockStmt) {
	m.Complexity = 1
	ast.Inspect(b, m.complexityNode)
}

func (m *Metrics) calculateLineCounts(b *ast.BlockStmt) {

}
