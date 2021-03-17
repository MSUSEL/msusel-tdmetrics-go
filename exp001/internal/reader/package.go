package reader

import (
	"go/ast"
	"go/types"

	"github.com/MSUSEL/msusel-tdmetrics-go/exp001/internal/utils"
)

// Package is the collection of compiled data for the package.
type Package struct {
	Name       string
	Project    *Project
	Package    *types.Package
	Types      map[ast.Expr]types.TypeAndValue
	Defs       map[*ast.Ident]types.Object
	Uses       map[*ast.Ident]types.Object
	Implicits  map[ast.Node]types.Object
	Selections map[*ast.SelectorExpr]*types.Selection
}

// SourceFilePaths gets the list source file paths for this package.
func (p *Package) SourceFilePaths() []string {
	files := utils.NewStringSet()
	for id := range p.Defs {
		files.Add(p.Project.FileSet.File(id.Pos()).Name())
	}
	return files.Values()
}

// UsedSignatures are all the functions and methods in all packages referenced by this package.
func (p *Package) UsedSignatures() map[*ast.Ident]*types.Signature {
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

// DefinedFuncs are all the functions and methods in this package.
func (p *Package) DefinedFuncs() map[*ast.Ident]*types.Func {
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
// and aliases from this package.
func (p *Package) TypeDefs() map[*ast.Ident]*types.TypeName {
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

// Participation gets all the functions which each type definition
// from the main package has participated in.
func (p *Package) Participation() map[*ast.Ident][]*ast.Ident {
	result := map[*ast.Ident][]*ast.Ident{}
	typeDefs := p.TypeDefs()
	bt := NewBaseTypes()
	defFuncs := p.DefinedFuncs()
	for defID, def := range typeDefs {
		funcIDs := []*ast.Ident{}
		for funcID, funcObj := range defFuncs {
			parts := bt.getBaseTypes(funcObj.Type())
			if parts[def.Type()] {
				funcIDs = append(funcIDs, funcID)
			}
			result[defID] = funcIDs
		}
	}

	// for defID, def := range typeDefs {
	// 	def.Type()
	// }

	// TODO: Propagate to sub-types

	//fmt.Println(typeDefs)

	return result
}
