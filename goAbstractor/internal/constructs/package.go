package constructs

import (
	"strings"

	"github.com/Snow-Gremlin/goToolbox/collections"
	"github.com/Snow-Gremlin/goToolbox/collections/enumerator"
	"github.com/Snow-Gremlin/goToolbox/terrors/terror"
	"golang.org/x/tools/go/packages"

	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/assert"
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

		addImport(p Package) Package
		addClass(c Class) Class
		addInterDef(id InterDef) InterDef
		addMethod(m Method) Method
		addValue(v Value) Value

		empty() bool
		findType(name string) Definition
		allTypes() collections.Enumerator[Definition]

		resolveReceivers()
		resolveClassInterfaces(proj Project)
		removeImports(predicate func(Construct) bool)
	}

	PackageArgs struct {
		RealPkg     *packages.Package
		Path        string
		Name        string
		ImportPaths []string
	}

	packageImp struct {
		pkg *packages.Package

		path        string
		name        string
		index       int
		importPaths []string
		imports     Set[Package]

		classes   Set[Class]
		interDefs Set[InterDef]
		methods   Set[Method]
		values    Set[Value]
	}
)

func newPackage(args PackageArgs) Package {
	assert.ArgNotNil(`real type`, args.RealPkg)
	assert.ArgNotEmpty(`path`, args.Path)
	assert.ArgValidId(`name`, args.Name)

	return &packageImp{
		pkg:         args.RealPkg,
		path:        args.Path,
		name:        args.Name,
		classes:     NewSet[Class](),
		imports:     NewSet[Package](),
		interDefs:   NewSet[InterDef](),
		methods:     NewSet[Method](),
		values:      NewSet[Value](),
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

func (p *packageImp) addImport(i Package) Package {
	return p.imports.Insert(i)
}

func (p *packageImp) addClass(c Class) Class {
	return p.classes.Insert(c)
}

func (p *packageImp) addInterDef(id InterDef) InterDef {
	return p.interDefs.Insert(id)
}

func (p *packageImp) addMethod(m Method) Method {
	return p.methods.Insert(m)
}

func (p *packageImp) addValue(v Value) Value {
	return p.values.Insert(v)
}

func (p *packageImp) empty() bool {
	return p.allTypes().Empty()
}

func (p *packageImp) findType(name string) Definition {
	def, _ := p.allTypes().Where(func(t Definition) bool { return t.Name() == name }).First()
	return def
}

func (p *packageImp) allTypes() collections.Enumerator[Definition] {
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

	ctx2 := ctx.HideKind().Short()
	return jsonify.NewMap().
		AddIf(ctx, ctx.IsKindShown(), `kind`, p.Kind()).
		AddNonZero(ctx2, `path`, p.path).
		AddNonZero(ctx2, `name`, p.name).
		AddNonZero(ctx2, `classes`, p.classes).
		AddNonZero(ctx2, `imports`, p.imports).
		AddNonZero(ctx2, `interDefs`, p.interDefs).
		AddNonZero(ctx2, `methods`, p.methods).
		AddNonZero(ctx2, `values`, p.values)
}

func (p *packageImp) Visit(v visitor.Visitor) {
	visitor.VisitList(v, p.imports.Values())
	visitor.VisitList(v, p.classes.Values())
	visitor.VisitList(v, p.interDefs.Values())
	visitor.VisitList(v, p.methods.Values())
	visitor.VisitList(v, p.values.Values())
}

func (p *packageImp) resolveReceivers() {
	methods := p.methods.Values()
	for i := range methods.Count() {
		m := methods.Get(i)
		if rec := m.ReceiverName(); len(rec) > 0 {
			t := p.findType(rec)
			if t == nil {
				panic(terror.New(`failed to find receiver`).
					With(`name`, rec))
			}
			c, ok := t.(Class)
			if !ok {
				panic(terror.New(`receiver was not a class`).
					With(`name`, rec).
					WithType(`gotten type`, t).
					With(`gotten value`, t))
			}
			c.addMethod(m)
		}
	}
}

func (p *packageImp) resolveClassInterfaces(proj Project) {
	classes := p.classes.Values()
	for i := range classes.Count() {
		classes.Get(i).resolveInterface(proj, p)
	}
}

func (p *packageImp) removeImports(predicate func(Construct) bool) {
	p.imports.Remove(predicate)
}
