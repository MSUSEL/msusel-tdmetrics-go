package constructs

import (
	"fmt"
	"strings"

	"github.com/Snow-Gremlin/goToolbox/collections"
	"github.com/Snow-Gremlin/goToolbox/collections/enumerator"
	"github.com/Snow-Gremlin/goToolbox/utils"
	"golang.org/x/tools/go/packages"

	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs/kind"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/jsonify"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/visitor"
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

		path      string
		name      string
		imports   Set[Package]
		interDefs Set[InterDef]
		classes   Set[Class]
		values    Set[Value]
		methods   Set[Method]

		index       int
		importPaths []string
	}
)

func newPackage(args PackageArgs) Package {

	// TODO: FIX
	if utils.IsNil(args.RealPkg) {
		panic(fmt.Errorf(`must provide a real package for %s`, args.Name))
	}
	return &packageImp{
		pkg:         args.RealPkg,
		path:        args.Path,
		name:        args.Name,
		imports:     NewSet[Package](),
		interDefs:   NewSet[InterDef](),
		classes:     NewSet[Class](),
		values:      NewSet[Value](),
		methods:     NewSet[Method](),
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
	def, _ := p.AllTypes().Where(func(t Definition) bool { return t.Name() == name }).First()
	return def
}

func (p *packageImp) AllTypes() collections.Enumerator[Definition] {
	i := enumerator.Cast[Definition](p.interDefs.Values().Enumerate())
	c := enumerator.Cast[Definition](p.classes.Values().Enumerate())
	v := enumerator.Cast[Definition](p.values.Values().Enumerate())
	m := enumerator.Cast[Definition](p.methods.Values().Enumerate())
	return i.Concat(c).Concat(v).Concat(m)
}

func (p *packageImp) CompareTo(other Construct) int {
	b := other.(*packageImp)
	if cmp := strings.Compare(p.path, b.path); cmp != 0 {
		return cmp
	}
	return strings.Compare(p.name, b.name)
}

func (p *packageImp) ToJson(ctx *jsonify.Context) jsonify.Datum {
	if ctx.IsShort() {
		return jsonify.New(ctx, p.index)
	}

	ctx2 := ctx.HideKind()
	return jsonify.NewMap().
		AddIf(ctx, ctx.IsKindShown(), `kind`, p.Kind()).
		AddNonZero(ctx2, `path`, p.path).
		AddNonZero(ctx2, `name`, p.name).
		AddNonZero(ctx2.Short(), `imports`, p.imports).
		AddNonZero(ctx2.Long(), `interDefs`, p.interDefs).
		AddNonZero(ctx2.Long(), `classes`, p.classes).
		AddNonZero(ctx2.Long(), `values`, p.values).
		AddNonZero(ctx2.Long(), `methods`, p.methods)
}

func (p *packageImp) Visit(v visitor.Visitor) bool {
	return visitor.VisitList(v, p.imports.Values()) &&
		visitor.VisitList(v, p.interDefs.Values()) &&
		visitor.VisitList(v, p.classes.Values()) &&
		visitor.VisitList(v, p.values.Values()) &&
		visitor.VisitList(v, p.methods.Values())
}
