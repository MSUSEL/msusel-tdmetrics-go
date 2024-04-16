package cyclomatic

import "go/ast"

type Cyclomatic struct {
	start *node
	end   *node
}

type node struct {
	kind       string
	statements []ast.Stmt
	next       []*node
}

func New(b *ast.BlockStmt) *Cyclomatic {
	c := &Cyclomatic{
		start: &node{kind: `start`},
		end:   &node{kind: `end`},
	}
	addStatements(c.start, c.end, b)
	return c
}

func addStatements(start, end *node, b *ast.BlockStmt) {
	for _, statement := range b.List {
		switch statement.(type) {
		case *ast.DeclStmt, *ast.EmptyStmt, *ast.ExprStmt,
			*ast.SendStmt, *ast.IncDecStmt, *ast.AssignStmt:
			break
		case *ast.LabeledStmt:
		case *ast.GoStmt:
		case *ast.DeferStmt:
		case *ast.ReturnStmt:
		case *ast.BranchStmt:
		case *ast.BlockStmt:
		case *ast.IfStmt:
		case *ast.SwitchStmt:
		case *ast.TypeSwitchStmt:
		case *ast.SelectStmt:
		case *ast.ForStmt:
		case *ast.RangeStmt:

		default: // *ast.BadStmt, *ast.CommClause, *ast.CaseClause:

		}
	}
}
