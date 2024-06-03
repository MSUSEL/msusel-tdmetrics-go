package constructs

import (
	"golang.org/x/tools/go/packages"

	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/jsonify"
)

type Package interface {
	Source() *packages.Package
	FindTypeDef(name string) TypeDef
	SetIndex(index int)
	Path() string
	ImportPaths() []string
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
	imports []Package
	types   []TypeDef
	values  []ValueDef
	methods []Method

	index       int
	importPaths []string
}

func NewPackage(pkg *packages.Package, path string, importPaths []string) Package {
	return &packageImp{
		pkg:         pkg,
		path:        path,
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

	ctx2 := ctx.ShowKind().Short()
	return jsonify.NewMap().
		AddNonZero(ctx2, `path`, p.path).
		AddNonZero(ctx2, `imports`, p.imports).
		AddNonZero(ctx2, `types`, p.types).
		AddNonZero(ctx2, `values`, p.values).
		AddNonZero(ctx2, `methods`, p.methods)
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

func (p *packageImp) SetIndex(index int) {
	p.index = index
}

func (p *packageImp) Path() string {
	return p.path
}

func (p *packageImp) ImportPaths() []string {
	return p.importPaths
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
