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
	s := scope.New().SetRange(c.enter, c.exit)
	addScope(s, b)
	return c
}

func (c *Cyclomatic) String() string {
	return node.Graph(c.enter)
}

func (c *Cyclomatic) Mermaid() string {
	return node.Mermaid(c.enter)
}

func addScope(s scope.Scope, stmt ast.Stmt) {
	if !addStatements(s, stmt) {
		s.Get(scope.Begin).AddNext(s.Get(scope.End))
	}
}

func setScopeLabels(s scope.Scope, statements []ast.Stmt) {
	for _, statement := range statements {
		if stmt, ok := statement.(*ast.LabeledStmt); ok {
			name := stmt.Label.Name
			labelNode := node.New(`label_`+name, stmt.Pos())
			s.Set(name, labelNode)
		}
	}
}

func addStatements(s scope.Scope, statements ...ast.Stmt) bool {
	setScopeLabels(s, statements)

	for _, statement := range statements {
		switch stmt := statement.(type) {
		case *ast.DeclStmt, *ast.EmptyStmt, *ast.ExprStmt,
			*ast.SendStmt, *ast.IncDecStmt, *ast.AssignStmt,
			*ast.GoStmt:
			// Non-branching
			break
		case *ast.BlockStmt:
			addStatements(s, stmt.List...)
		case *ast.LabeledStmt:
			addLabeledStmt(s, stmt)
		case *ast.DeferStmt:
			addDeferStmt(s, stmt)
		case *ast.ReturnStmt:
			addReturnStmt(s, stmt)
			return true
		case *ast.BranchStmt:
			addBranchStmt(s, stmt)
			return true
		case *ast.IfStmt:
			addIfStmt(s, stmt)
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
	return false
}

func addLabeledStmt(s scope.Scope, stmt *ast.LabeledStmt) {
	labelNode := s.Get(stmt.Label.Name)
	s.Get(scope.Begin).AddNext(labelNode)
	s.Set(scope.Begin, labelNode)
	if stmt.Stmt != nil {
		addStatements(s, stmt.Stmt)
	}
}

func addDeferStmt(s scope.Scope, stmt *ast.DeferStmt) {
	deferNode := node.New(`defer`, stmt.Pos())
	deferNode.AddNext(s.Get(scope.End))
	s.Set(scope.End, deferNode)
	// TODO: need to handle a function expression with a body.
}

func addReturnStmt(s scope.Scope, _ *ast.ReturnStmt) {
	s.Get(scope.Begin).AddNext(s.Get(scope.End))
}

func addBranchStmt(s scope.Scope, stmt *ast.BranchStmt) {
	switch stmt.Tok {
	case token.BREAK:
		panic(errors.New(`TODO: Implement Break branch`))
	case token.CONTINUE:
		panic(errors.New(`TODO: Implement Continue branch`))
	case token.GOTO:
		addGotoBranch(s, stmt)
	case token.FALLTHROUGH:
		panic(errors.New(`TODO: Implement Fallthrough branch`))
	default:
		panic(fmt.Errorf(`unexpected branch type: %s`, stmt.Tok.String()))
	}
}

func addGotoBranch(s scope.Scope, stmt *ast.BranchStmt) {
	name := stmt.Label.Name
	gtNode := node.New(`goto_`+name, stmt.Pos())
	s.Get(scope.Begin).AddNext(gtNode)
	label := s.Get(name)
	gtNode.AddNext(label)
	s.Set(scope.Begin, label)
}

func addIfStmt(s scope.Scope, stmt *ast.IfStmt) {
	// TODO: handle stmt.Cond
	// TODO: handle stmt.Init

	ifNode := node.New(`if`, stmt.Pos())
	s.Get(scope.Begin).AddNext(ifNode)

	ifBodyNode := node.New(`ifBody`, stmt.Body.Pos())
	endIfNode := node.New(`endIf`, stmt.End())
	ifNode.AddNext(ifBodyNode)

	s2 := s.Push().SetRange(ifBodyNode, endIfNode)
	addScope(s2, stmt.Body)

	if stmt.Else == nil {
		ifNode.AddNext(endIfNode)
	} else {
		fmt.Printf("stmt.Else %T: %v\n", stmt.Else, stmt.Else)

		elseBodyNode := node.New(`elseBody`, stmt.Else.Pos())
		ifNode.AddNext(elseBodyNode)

		s3 := s.Push().SetRange(elseBodyNode, endIfNode)
		addScope(s3, stmt.Else)
	}

	s.Set(scope.Begin, endIfNode)
}
