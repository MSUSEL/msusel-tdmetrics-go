package reader

import (
	"go/ast"
	"go/token"
	"go/types"
	"sort"
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
	fileSet := map[string]bool{}
	for id := range p.Defs {
		fileSet[p.FileSet.File(id.Pos()).Name()] = true
	}
	files := make([]string, 0, len(fileSet))
	for file := range fileSet {
		files = append(files, file)
	}
	sort.Strings(files)
	return files
}

// Signatures are all the functions and methods.
func (p *Project) Signatures() map[*ast.Ident]*types.Signature {
	signs := make(map[*ast.Ident]*types.Signature)
	for id, obj := range p.Uses {
		switch t := obj.Type().(type) {
		case *types.Signature:
			signs[id] = t
		}
	}
	return signs
}

func (p *Project) TypeDefs() {

}
