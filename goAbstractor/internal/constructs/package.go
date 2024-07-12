package constructs

import (
	"fmt"
	"slices"

	"github.com/Snow-Gremlin/goToolbox/utils"
	"golang.org/x/tools/go/packages"

	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/jsonify"
)

type (
	Package interface {
		Construct
		Source() *packages.Package
		Path() string
		Name() string
		ImportPaths() []string
		Imports() []Package

		Empty() bool
		Types() []TypeDef
		Methods() []Method

		Prune(predicate func(f any) bool)
		FindTypeDef(name string) TypeDef
		SetImports(imports []Package)
		AppendTypes(typeDef ...TypeDef)
		AppendValues(value ...ValueDef)
		SetMethods(methods []Method)
		AppendMethods(methods ...Method)

		setIndices(pkgIndex, typeDefIndex int) int
	}

	PackageArgs struct {
		RealPkg     *packages.Package
		Path        string
		Name        string
		ImportPaths []string
	}

	packageImp struct {
		pkg *packages.Package

		path    string
		name    string
		imports []Package

		interfaces []Named
		variables  []Named
		constants  []Named
		classes    []Class
		methods    []Method

		index       int
		importPaths []string
	}
)

func newPackage(args PackageArgs) Package {
	if utils.IsNil(args.RealPkg) {
		panic(fmt.Errorf(`must provide a real package for %s`, args.Name))
	}
	return &packageImp{
		pkg:         args.RealPkg,
		path:        args.Path,
		name:        args.Name,
		importPaths: args.ImportPaths,
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
		AddNonZero(ctx2.Long(), `types`, p.structs).
		AddNonZero(ctx2.Long(), `values`, p.values).
		AddNonZero(ctx2.Long(), `methods`, p.methods)
}

func (p *packageImp) Visit(v Visitor) {
	visitList(v, p.imports)
	visitList(v, p.structs)
	visitList(v, p.values)
	visitList(v, p.methods)
}

func castPred[T any](predicate func(f any) bool) func(f T) bool {
	return func(f T) bool { return predicate(f) }
}

func (p *packageImp) Prune(predicate func(f any) bool) {
	p.structs = slices.DeleteFunc(p.structs, castPred[TypeDef](predicate))
	p.values = slices.DeleteFunc(p.values, castPred[ValueDef](predicate))
	p.methods = slices.DeleteFunc(p.methods, castPred[Method](predicate))
}

func (p *packageImp) String() string {
	return jsonify.ToString(p)
}

func (p *packageImp) FindTypeDef(name string) TypeDef {
	for _, t := range p.structs {
		if name == t.Name() {
			return t
		}
	}
	return nil
}

func (p *packageImp) setIndices(pkgIndex, typeDefIndex int) int {
	p.index = pkgIndex
	for _, td := range p.structs {
		td.SetIndex(typeDefIndex)
		typeDefIndex++
	}
	return typeDefIndex
}

func (p *packageImp) Empty() bool {
	return len(p.structs) <= 0 &&
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
	return p.structs
}

func (p *packageImp) AddInterface(it ...Interface) {
	p.interfaces = append(p.interfaces, it...)
}

func (p *packageImp) AppendTypes(typeDef ...TypeDef) {
	p.structs = append(p.structs, typeDef...)
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
