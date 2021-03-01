package reader

import (
	"go/ast"
	"go/token"
	"go/types"
)

// Project is the collection of compiled data for the project.
type Project struct {
	Fileset    *token.FileSet
	Package    *types.Package
	Types      map[ast.Expr]types.TypeAndValue
	Defs       map[*ast.Ident]types.Object
	Uses       map[*ast.Ident]types.Object
	Implicits  map[ast.Node]types.Object
	Selections map[*ast.SelectorExpr]*types.Selection
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

func addPrticipationType(parts map[types.Type]map[*ast.Ident]bool, t types.Type, id *ast.Ident) {
	// if types.Id() types.IsUntyped(t) {
	// 	return
	// }
	if idMap := parts[t]; idMap != nil {
		idMap[id] = true
		return
	}
	idMap := map[*ast.Ident]bool{id: true}
	parts[t] = idMap
}

func (p *Project) Participation() map[types.Type][]*ast.Ident {
	parts := make(map[types.Type]map[*ast.Ident]bool)
	for id, sign := range p.Signatures() {
		if sign.Recv() != nil {
			addPrticipationType(parts, sign.Recv().Type(), id)
		}
		for i := sign.Params().Len() - 1; i >= 0; i-- {
			addPrticipationType(parts, sign.Params().At(i).Type(), id)
		}
	}

	result := make(map[types.Type][]*ast.Ident)
	for t, idMap := range parts {
		ids := make([]*ast.Ident, 0, len(idMap))
		for id := range idMap {
			ids = append(ids, id)
		}
		result[t] = ids
	}
	return result
}
