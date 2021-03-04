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
	result := make(map[*ast.Ident]*types.Signature)
	for id, obj := range p.Uses {
		if obj != nil {
			switch t := obj.Type().(type) {
			case *types.Signature:
				result[id] = t
			}
		}
	}
	return result
}

// DefinedFuncs are all the functions and methods in the main package.
func (p *Project) DefinedFuncs() map[*ast.Ident]*types.Func {
	result := make(map[*ast.Ident]*types.Func)
	for id, obj := range p.Defs {
		if obj != nil {
			switch t := obj.(type) {
			case *types.Func:
				result[id] = t
			}
		}
	}
	return result
}

// TypeDefs are all the type definitions for structures, interfaces, functions signatures,
// and aliases from the main package.
func (p *Project) TypeDefs() map[*ast.Ident]*types.TypeName {
	result := make(map[*ast.Ident]*types.TypeName)
	for id, obj := range p.Defs {
		if obj != nil {
			switch t := obj.(type) {
			case *types.TypeName:
				result[id] = t
			}
		}
	}
	return result
}

// getParticipants gets all the participants (receiver types and parameters types).
func getParticipants(f *types.Func) map[types.Type]bool {
	result := map[types.Type]bool{}
	if sig, ok := f.Type().(*types.Signature); ok {
		if recv := sig.Recv(); recv != nil {
			result[recv.Type()] = true
		}
		params := sig.Params()
		for i := 0; i < params.Len(); i++ {
			param := params.At(i)
			result[param.Type()] = true
		}
	}
	return result
}

// Participation gets all the functions which each type definition
// from the main package has participated in.
func (p *Project) Participation() map[*ast.Ident][]*ast.Ident {
	result := map[*ast.Ident][]*ast.Ident{}
	typeDefs := p.TypeDefs()
	defFuncs := p.DefinedFuncs()
	for defID, def := range typeDefs {
		funcIDs := []*ast.Ident{}
		for funcID, funcObj := range defFuncs {
			parts := getParticipants(funcObj)
			if parts[def.Type()] {
				funcIDs = append(funcIDs, funcID)
			}
			result[defID] = funcIDs
		}
	}
	return result
}
