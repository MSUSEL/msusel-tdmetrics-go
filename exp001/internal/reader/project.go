package reader

import (
	"go/ast"
	"go/token"
	"go/types"

	"github.com/MSUSEL/msusel-tdmetrics-go/exp001/internal/utils"
)

// Project is the collection of compiled data for the project.
type Project struct {
	BasePath   string
	FileSet    *token.FileSet
	Package    *types.Package
	Types      map[ast.Expr]types.TypeAndValue
	Defs       map[*ast.Ident]types.Object
	Uses       map[*ast.Ident]types.Object
	Implicits  map[ast.Node]types.Object
	Selections map[*ast.SelectorExpr]*types.Selection
}

// SourceFilePaths gets the list source file paths for this project.
func (p *Project) SourceFilePaths() []string {
	files := utils.NewStringSet()
	for id := range p.Defs {
		files.Add(p.FileSet.File(id.Pos()).Name())
	}
	return files.Values()
}

// UsedSignatures are all the functions and methods in all packages
func (p *Project) UsedSignatures() map[*ast.Ident]*types.Signature {
	signs := make(map[*ast.Ident]*types.Signature)
	for id, obj := range p.Uses {
		if obj != nil {
			switch t := obj.Type().(type) {
			case *types.Signature:
				signs[id] = t
			}
		}
	}
	return signs
}

// DefinedFuncs are all the functions and methods in the main package.
func (p *Project) DefinedFuncs() map[*ast.Ident]*types.Func {
	funcs := make(map[*ast.Ident]*types.Func)
	for id, obj := range p.Defs {
		if obj != nil {
			switch t := obj.(type) {
			case *types.Func:
				funcs[id] = t
			}
		}
	}
	return funcs
}

// TypeDefs are all the type definitions for structures, interfaces, functions signatures, and aliases.
func (p *Project) TypeDefs() map[*ast.Ident]*types.TypeName {
	typeDefs := make(map[*ast.Ident]*types.TypeName)
	for id, obj := range p.Defs {
		if obj != nil {
			switch t := obj.(type) {
			case *types.TypeName:
				typeDefs[id] = t
			}
		}
	}
	return typeDefs
}
