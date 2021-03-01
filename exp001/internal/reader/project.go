package reader

import (
	"go/ast"
	"go/token"
	"go/types"
)

type Project struct {
	Fileset    *token.FileSet
	Package    *types.Package
	Types      map[ast.Expr]types.TypeAndValue
	Defs       map[*ast.Ident]types.Object
	Uses       map[*ast.Ident]types.Object
	Implicits  map[ast.Node]types.Object
	Selections map[*ast.SelectorExpr]*types.Selection
}
