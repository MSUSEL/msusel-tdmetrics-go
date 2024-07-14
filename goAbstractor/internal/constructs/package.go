package constructs

import (
	"fmt"

	"github.com/Snow-Gremlin/goToolbox/collections"
	"github.com/Snow-Gremlin/goToolbox/utils"
	"golang.org/x/tools/go/packages"

	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs/kind"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/jsonify"
)

type (
	Package interface {
		Construct
		_package()

		Source() *packages.Package
		Path() string
		Name() string
		ImportPaths() []string
		Imports() collections.ReadonlyList[Package]

		FindType(name string) Definition
		AllTypes() collections.Enumerator[Definition]
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
		imports Set[Package]

		interfaces Set[InterDef]
		values     Set[Value]
		classes    Set[Class]
		methods    Set[Method]

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

func (p *packageImp) _package()                 {}
func (p *packageImp) Kind() kind.Kind           { return kind.Package }
func (p *packageImp) SetIndex(index int)        { p.index = index }
func (p *packageImp) Source() *packages.Package { return p.pkg }
func (p *packageImp) Path() string              { return p.path }
func (p *packageImp) Name() string              { return p.name }
func (p *packageImp) ImportPaths() []string     { return p.importPaths }

func (p *packageImp) Imports() collections.ReadonlyList[Package] {
	return p.imports.Values()
}

func (p *packageImp) FindType(name string) Definition {

}

func (p *packageImp) AllTypes() collections.Enumerator[Definition] {

}

func (p *packageImp) CompareTo(other Construct) int {

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
