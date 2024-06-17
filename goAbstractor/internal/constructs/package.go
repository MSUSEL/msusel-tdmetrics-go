package constructs

import (
	"fmt"

	"golang.org/x/tools/go/packages"

	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/jsonify"
	"github.com/Snow-Gremlin/goToolbox/utils"
)

type Package interface {
	Visitable
	Source() *packages.Package
	FindTypeDef(name string) TypeDef
	SetIndices(pkgIndex, typeDefIndex int) int
	Empty() bool
	Path() string
	Name() string
	ImportPaths() []string
	Imports() []Package
	SetImports(imports []Package)
	Types() []TypeDef
	AppendTypes(typeDef ...TypeDef)
	AppendValues(value ...ValueDef)
	Methods() []Method
	SetMethods(methods []Method)
	AppendMethods(methods ...Method)
}

type packageImp struct {
	pkg *packages.Package

	path    string
	name    string
	imports []Package
	types   []TypeDef
	values  []ValueDef
	methods []Method

	index       int
	importPaths []string
}

func NewPackage(pkg *packages.Package, path, name string, importPaths []string) Package {
	if utils.IsNil(pkg) {
		panic(fmt.Errorf(`must provide a real package for %s`, name))
	}
	return &packageImp{
		pkg:         pkg,
		path:        path,
		name:        name,
		importPaths: importPaths,
	}
}

func (p *packageImp) Source() *packages.Package {
	return p.pkg
}

func (p *packageImp) ToJson(ctx *jsonify.Context) jsonify.Datum {
	if ctx.IsShort() {
		return jsonify.New(ctx, p.index)
	}

	ctx2 := ctx.HideKind()
	return jsonify.NewMap().
		AddIf(ctx, ctx.IsKindShown(), `kind`, `package`).
		AddNonZero(ctx2, `path`, p.path).
		AddNonZero(ctx2, `name`, p.name).
		AddNonZero(ctx2.Short(), `imports`, p.imports).
		AddNonZero(ctx2.Long(), `types`, p.types).
		AddNonZero(ctx2.Long(), `values`, p.values).
		AddNonZero(ctx2.Long(), `methods`, p.methods)
}

func (p *packageImp) Visit(v Visitor) {
	visitList(v, p.imports)
	visitList(v, p.types)
	visitList(v, p.values)
	visitList(v, p.methods)
}

func (p *packageImp) String() string {
	return jsonify.ToString(p)
}

func (p *packageImp) FindTypeDef(name string) TypeDef {
	for _, t := range p.types {
		if name == t.Name() {
			return t
		}
	}
	return nil
}

func (p *packageImp) SetIndices(pkgIndex, typeDefIndex int) int {
	p.index = pkgIndex
	for _, td := range p.types {
		td.SetIndex(typeDefIndex)
		typeDefIndex++
	}
	return typeDefIndex
}

func (p *packageImp) Empty() bool {
	return len(p.types) <= 0 &&
		len(p.values) <= 0 &&
		len(p.methods) <= 0
}

func (p *packageImp) Path() string {
	return p.path
}

func (p *packageImp) Name() string {
	return p.name
}

func (p *packageImp) ImportPaths() []string {
	return p.importPaths
}

func (p *packageImp) Imports() []Package {
	return p.imports
}

func (p *packageImp) SetImports(imports []Package) {
	p.imports = imports
}

func (p *packageImp) Types() []TypeDef {
	return p.types
}

func (p *packageImp) AppendTypes(typeDef ...TypeDef) {
	p.types = append(p.types, typeDef...)
}

func (p *packageImp) AppendValues(value ...ValueDef) {
	p.values = append(p.values, value...)
}

func (p *packageImp) Methods() []Method {
	return p.methods
}

func (p *packageImp) SetMethods(methods []Method) {
	p.methods = methods
}

func (p *packageImp) AppendMethods(methods ...Method) {
	p.methods = append(p.methods, methods...)
}
