package cyclomatic

import (
	"errors"
	"fmt"
	"go/ast"
	"go/token"
	"strings"
)

type Cyclomatic struct {
	enter *node
	exit  *node
}

func New(b *ast.BlockStmt) *Cyclomatic {
	c := &Cyclomatic{
		enter: newNode(`enter`, b.Pos()),
		exit:  newNode(`exit`, b.End()),
	}
	s := newScope()
	s.setTag(enterTag, c.enter)
	s.setTag(exitTag, c.exit)
	addStatements(s, b.List)
	return c
}

func (c *Cyclomatic) String() string {
	buf := &strings.Builder{}
	touched := map[string]bool{}
	c.enter.format(buf, touched, `─`, ` `)
	return buf.String()
}

func addStatements(s *scope, statements []ast.Stmt) {
	for _, statement := range statements {
		if stmt, ok := statement.(*ast.LabeledStmt); ok {
			labelNode := newNode(`label`, stmt.Pos())
			s.setTag(stmt.Label.Name, labelNode)
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
	s.getTag(enterTag).addNext(s.getTag(exitTag))
}

func addLabeledStmt(s *scope, stmt *ast.LabeledStmt) {
	labelNode := s.getTag(stmt.Label.Name)
	s.getTag(enterTag).addNext(labelNode)
	s.setTag(enterTag, labelNode)
	if stmt.Stmt != nil {
		// TODO: determine what the stmt.Stmt in the label is.
		//       I assume this is the switch, select, for being labelled.
		panic(fmt.Errorf(`TODO: Implement label statement: %T`, stmt.Stmt))
	}
}

func addDeferStmt(s *scope, stmt *ast.DeferStmt) {
	deferNode := newNode(`defer`, stmt.Pos())
	deferNode.addNext(s.getTag(exitTag))
	s.setTag(exitTag, deferNode)
	// TODO: need to handle a function expression with a body.
}

func addReturnStmt(s *scope, _ *ast.ReturnStmt, i int, statements []ast.Stmt) {
	if remainder := len(statements) - 1 - i; remainder > 0 {
		panic(fmt.Errorf(`unexpected %d statements after return`, remainder))
	}
	s.getTag(enterTag).addNext(s.getTag(exitTag))
}

func addBranchStmt(_ *scope, stmt *ast.BranchStmt) {
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
