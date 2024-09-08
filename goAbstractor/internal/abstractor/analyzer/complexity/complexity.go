package complexity

import (
	"go/ast"
	"go/token"
	"math"
)

type Complexity struct {
	Complexity int
	LineCount  int
	CodeCount  int
	Indents    int
}

type complexityImp struct {
	fSet       *token.FileSet
	complexity int
	minLine    int
	maxLine    int
	indents    int
	minColumn  map[int]int
}

// Calculate gathers positional information for indents and cyclomatic complexity.
func Calculate(node ast.Node, fSet *token.FileSet) Complexity {
	c := &complexityImp{
		fSet:       fSet,
		complexity: 1,
		maxLine:    0,
		minLine:    math.MaxInt,
		indents:    0,
		minColumn:  map[int]int{},
	}

	ast.Inspect(node, c.addCodePosForNode)

	return Complexity{
		Complexity: c.complexity,
		LineCount:  c.maxLine - c.minLine + 1,
		CodeCount:  len(c.minColumn),
		Indents:    c.calcIndents(),
	}
}

func (c *complexityImp) calcIndents() int {
	leftMostColumn := math.MaxInt
	indentSum := 0
	for _, ind := range c.minColumn {
		leftMostColumn = min(ind, leftMostColumn)
		indentSum += ind
	}
	return indentSum - len(c.minColumn)*leftMostColumn
}

func (c *complexityImp) incComplexity(check bool) {
	if check {
		c.complexity++
	}
}

func (c *complexityImp) addCodePos(pos token.Pos, isEnd bool) {
	p := c.fSet.PositionFor(pos, false)
	lineNo, column := p.Line, p.Column
	c.maxLine = max(c.maxLine, lineNo)
	c.minLine = min(c.minLine, lineNo)
	if isEnd {
		column--
	}
	if otherCol, ok := c.minColumn[lineNo]; ok {
		column = min(column, otherCol)
	}
	c.minColumn[lineNo] = column
}

func (c *complexityImp) addCodePosForNode(n ast.Node) bool {
	switch t := n.(type) {
	case nil, *ast.Comment, *ast.CommentGroup:
		return true
	case *ast.IfStmt, *ast.ForStmt, *ast.RangeStmt, *ast.GoStmt:
		c.incComplexity(true)
	case *ast.CaseClause:
		c.incComplexity(t.List != nil)
	case *ast.CommClause:
		c.incComplexity(t.Comm != nil)
	case *ast.BinaryExpr:
		c.incComplexity(t.Op == token.LAND || t.Op == token.LOR)
	}

	c.addCodePos(n.Pos(), false)
	if ended, has := n.(interface{ End() token.Pos }); has {
		c.addCodePos(ended.End(), true)
	}
	return true
}
