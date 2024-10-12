package usages

import (
	"go/ast"
	"go/token"
	"go/types"

	"github.com/Snow-Gremlin/goToolbox/terrors/terror"
)

type posReader interface {
	Pos() token.Pos
}

func isLocal(root ast.Node, query posReader) bool {
	pos := query.Pos()
	return root.Pos() <= pos && pos <= root.End()
}

func isLocalType(root ast.Node, t types.Type) bool {
	named, ok := stripNamed(t)
	return ok && isLocal(root, named.Obj())
}

func unparen(node ast.Node) ast.Node {
	if p, ok := node.(*ast.ParenExpr); ok {
		return p.X
	}
	return node
}

func stripNamed(typ types.Type) (*types.Named, bool) {
	if pointer, ok := typ.(*types.Pointer); ok {
		typ = pointer.Elem()
	}
	named, ok := typ.(*types.Named)
	return named, ok
}

func getName(fSet *token.FileSet, expr ast.Expr) string {
	exp := unparen(expr)
	if id, ok := exp.(*ast.Ident); ok {
		return id.Name
	}
	if sel, ok := exp.(*ast.SelectorExpr); ok {
		src := unparen(sel.X)
		if id, ok := src.(*ast.Ident); ok {
			return id.Name + `.` + sel.Sel.Name
		}
		panic(terror.New(`unexpected expression in selection for name`).
			WithType(`type`, src).
			With(`expression`, src).
			With(`selection`, sel).
			With(`position`, fSet.Position(expr.Pos())))
	}
	panic(terror.New(`unexpected expression for name`).
		WithType(`type`, exp).
		With(`expression`, exp).
		With(`position`, fSet.Position(expr.Pos())))
}
