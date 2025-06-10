package querier

import (
	"go/ast"
	"go/token"
	"go/types"

	"github.com/Snow-Gremlin/goToolbox/terrors/terror"
	"golang.org/x/tools/go/packages"
)

type Querier struct {
	packages []*packages.Package
	info     *types.Info
	fSet     *token.FileSet
	fnScopes map[*types.Scope]*types.Func
}

func New(packages []*packages.Package) *Querier {
	if len(packages) == 0 {
		panic(terror.New(`no packages provided`))
	}

	q := &Querier{
		packages: packages,
		info:     packages[0].TypesInfo,
		fSet:     packages[0].Fset,
		fnScopes: map[*types.Scope]*types.Func{},
	}
	for _, obj := range q.info.Defs {
		if fn, ok := obj.(*types.Func); ok {
			q.fnScopes[fn.Scope()] = fn
		}
	}
	return q
}

func NewSimple(info *types.Info, fSet *token.FileSet) *Querier {
	return &Querier{
		info: info,
		fSet: fSet,
	}
}

func (q *Querier) Packages() []*packages.Package    { return q.packages }
func (q *Querier) Info() *types.Info                { return q.info }
func (q *Querier) FileSet() *token.FileSet          { return q.fSet }
func (q *Querier) Pos(pos token.Pos) token.Position { return q.fSet.Position(pos) }

func (q *Querier) ForeachPackage(handle func(*packages.Package)) {
	packages.Visit(q.packages, func(pkg *packages.Package) bool {
		handle(pkg)
		return true
	}, nil)
}

func (q *Querier) GetType(e ast.Expr) types.Type {
	if tv, has := q.info.Types[e]; has {
		return tv.Type
	}
	panic(terror.New(`type expression not found in types info`).
		WithType(`expression`, e).
		With(`pos`, q.Pos(e.Pos())))
}

func (q *Querier) GetDef(id *ast.Ident) types.Object {
	if obj, has := q.info.Defs[id]; has {
		return obj
	}
	panic(terror.New(`type declaration not found in types info`).
		WithType(`identifier`, id.Name).
		With(`pos`, q.Pos(id.Pos())))
}

func (q *Querier) NestingFunc(obj types.Object) *types.Func {
	if obj == nil {
		return nil
	}

	var pkgScope *types.Scope
	if pkg := obj.Pkg(); pkg != nil {
		pkgScope = pkg.Scope()
	}

	scope := obj.Parent()
	if scope == pkgScope {
		return nil
	}

	if scope == nil && pkgScope != nil && obj.Pos().IsValid() {
		scope = pkgScope.Innermost(obj.Pos())
	}

	for ; scope != nil; scope = scope.Parent() {
		if fn, ok := q.fnScopes[scope]; ok {
			return fn
		}
	}
	return nil
}
