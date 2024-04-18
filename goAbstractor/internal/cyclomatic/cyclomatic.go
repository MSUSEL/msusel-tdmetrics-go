package cyclomatic

import (
	"errors"
	"fmt"
	"go/ast"
	"go/token"

	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/cyclomatic/node"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/cyclomatic/scope"
)

type Cyclomatic struct {
	enter node.Node
	exit  node.Node
}

func New(b *ast.BlockStmt) *Cyclomatic {
	c := &Cyclomatic{
		enter: node.New(`enter`, b.Pos()),
		exit:  node.New(`exit`, b.End()),
	}
	s := scope.New()
	s.Set(scope.Enter, c.enter)
	s.Set(scope.Exit, c.exit)
	addStatements(s, b.List)
	return c
}

func (c *Cyclomatic) String() string {
	return c.enter.FullString()
}

func addStatements(s scope.Scope, statements []ast.Stmt) {
	for _, statement := range statements {
		if stmt, ok := statement.(*ast.LabeledStmt); ok {
			labelNode := node.New(`label`, stmt.Pos())
			s.Set(stmt.Label.Name, labelNode)
		}
	}

	for i, statement := range statements {
		switch stmt := statement.(type) {
		case *ast.DeclStmt, *ast.EmptyStmt, *ast.ExprStmt,
			*ast.SendStmt, *ast.IncDecStmt, *ast.AssignStmt,
			*ast.GoStmt:
			// Non-branching
			break
		case *ast.LabeledStmt:
			addLabeledStmt(s, stmt)
		case *ast.DeferStmt:
			addDeferStmt(s, stmt)
		case *ast.ReturnStmt:
			addReturnStmt(s, stmt, i, statements)
			return
		case *ast.BranchStmt:
			addBranchStmt(s, stmt)
		case *ast.BlockStmt:
			panic(fmt.Errorf(`TODO: Implement %T`, stmt))
		case *ast.IfStmt:
			panic(fmt.Errorf(`TODO: Implement %T`, stmt))
		case *ast.SwitchStmt:
			panic(fmt.Errorf(`TODO: Implement %T`, stmt))
		case *ast.TypeSwitchStmt:
			panic(fmt.Errorf(`TODO: Implement %T`, stmt))
		case *ast.SelectStmt:
			panic(fmt.Errorf(`TODO: Implement %T`, stmt))
		case *ast.ForStmt:
			panic(fmt.Errorf(`TODO: Implement %T`, stmt))
		case *ast.RangeStmt:
			panic(fmt.Errorf(`TODO: Implement %T`, stmt))
		default: // *ast.BadStmt, *ast.CommClause, *ast.CaseClause:
			panic(fmt.Errorf(`unexpected statement in block %T: %s`, statement, statement))
		}
	}
	s.Get(scope.Enter).AddNext(s.Get(scope.Exit))
}

func addLabeledStmt(s scope.Scope, stmt *ast.LabeledStmt) {
	labelNode := s.Get(stmt.Label.Name)
	s.Get(scope.Enter).AddNext(labelNode)
	s.Set(scope.Enter, labelNode)
	if stmt.Stmt != nil {
		// TODO: determine what the stmt.Stmt in the label is.
		//       I assume this is the switch, select, for being labelled.
		panic(fmt.Errorf(`TODO: Implement label statement: %T`, stmt.Stmt))
	}
}

func addDeferStmt(s scope.Scope, stmt *ast.DeferStmt) {
	deferNode := node.New(`defer`, stmt.Pos())
	deferNode.AddNext(s.Get(scope.Exit))
	s.Set(scope.Exit, deferNode)
	// TODO: need to handle a function expression with a body.
}

func addReturnStmt(s scope.Scope, _ *ast.ReturnStmt, i int, statements []ast.Stmt) {
	if remainder := len(statements) - 1 - i; remainder > 0 {
		panic(fmt.Errorf(`unexpected %d statements after return`, remainder))
	}
	s.Get(scope.Enter).AddNext(s.Get(scope.Exit))
}

func addBranchStmt(_ scope.Scope, stmt *ast.BranchStmt) {
	switch stmt.Tok {
	case token.BREAK:
		panic(errors.New(`TODO: Implement Break branch`))
	case token.CONTINUE:
		panic(errors.New(`TODO: Implement Continue branch`))
	case token.GOTO:
		panic(errors.New(`TODO: Implement Goto branch`))
	case token.FALLTHROUGH:
		panic(errors.New(`TODO: Implement Fallthrough branch`))
	default:
		panic(fmt.Errorf(`unexpected branch type: %s`, stmt.Tok.String()))
	}
}
