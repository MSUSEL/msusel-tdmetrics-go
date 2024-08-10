package packageCon

import (
	"strings"

	"github.com/Snow-Gremlin/goToolbox/collections"
	"github.com/Snow-Gremlin/goToolbox/collections/enumerator"
	"github.com/Snow-Gremlin/goToolbox/terrors/terror"
	"golang.org/x/tools/go/packages"

	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/assert"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/jsonify"
)

const Kind = `package`

type Package interface {
	constructs.Construct
	_package()

	Source() *packages.Package
	Path() string
	Name() string
	ImportPaths() []string
	Imports() collections.ReadonlyList[Package]
	InitCount() int

	addImport(p Package) Package
	addInterface(it Interface) Interface
	addMethod(m Method) Method
	addObject(id Object) Object
	addValue(v Value) Value

	empty() bool
	findDeclaration(name string) Declaration
	allDeclarations() collections.Enumerator[Declaration]

	resolveReceivers()
	removeImports(predicate func(Construct) bool)
}

type Args struct {
	RealPkg     *packages.Package
	Path        string
	Name        string
	ImportPaths []string
}

type packageImp struct {
	pkg *packages.Package

	path        string
	name        string
	index       int
	importPaths []string

	imports    Set[Package]
	interfaces Set[Interface]
	methods    Set[Method]
	objects    Set[Object]
	values     Set[Value]
}

func newPackage(args PackageArgs) Package {
	assert.ArgNotNil(`real type`, args.RealPkg)
	assert.ArgNotEmpty(`path`, args.Path)
	assert.ArgValidId(`name`, args.Name)

	return &packageImp{
		pkg:         args.RealPkg,
		path:        args.Path,
		name:        args.Name,
		importPaths: args.ImportPaths,
		imports:     NewSet[Package](),
		interfaces:  NewSet[Interface](),
		methods:     NewSet[Method](),
		objects:     NewSet[Object](),
		values:      NewSet[Value](),
	}
}

func (p *packageImp) _package()          {}
func (p *packageImp) Kind() string       { return Kind }
func (p *packageImp) setIndex(index int) { p.index = index }

func (p *packageImp) Source() *packages.Package { return p.pkg }
func (p *packageImp) Path() string              { return p.path }
func (p *packageImp) Name() string              { return p.name }
func (p *packageImp) ImportPaths() []string     { return p.importPaths }

func (p *packageImp) Imports() collections.ReadonlyList[Package] {
	return p.imports.Values()
}

func (p *packageImp) InitCount() int {
	count := 0
	methods := p.methods.Values()
	for i := range methods.Count() {
		if methods.Get(i).IsInit() {
			count++
		}
	}
	return count
}

func (p *packageImp) addImport(i Package) Package {
	return p.imports.Insert(i)
}

func (p *packageImp) addInterface(it Interface) Interface {
	return p.interfaces.Insert(it)
}

func (p *packageImp) addMethod(m Method) Method {
	return p.methods.Insert(m)
}

func (p *packageImp) addObject(d Object) Object {
	return p.objects.Insert(d)
}

func (p *packageImp) addValue(v Value) Value {
	return p.values.Insert(v)
}

func (p *packageImp) empty() bool {
	return p.allDeclarations().Empty()
}

func (p *packageImp) findDeclaration(name string) Declaration {
	def, _ := p.allDeclarations().
		Where(func(t Declaration) bool { return t.Name() == name }).
		First()
	return def
}

func (p *packageImp) allDeclarations() collections.Enumerator[Declaration] {
	i := enumerator.Cast[Declaration](p.interfaces.Values().Enumerate())
	m := enumerator.Cast[Declaration](p.methods.Values().Enumerate())
	o := enumerator.Cast[Declaration](p.objects.Values().Enumerate())
	v := enumerator.Cast[Declaration](p.values.Values().Enumerate())
	return i.Concat(m).Concat(o).Concat(v)
}

func (p *packageImp) compareTo(other Construct) int {
	b := other.(*packageImp)
	return or(
		func() int { return strings.Compare(p.path, b.path) },
		func() int { return strings.Compare(p.name, b.name) },
	)
}

func (p *packageImp) ToJson(ctx *jsonify.Context) jsonify.Datum {
	if ctx.IsShort() {
		return jsonify.New(ctx, p.index)
	}

	ctx2 := ctx.HideKind().Short()
	return jsonify.NewMap().
		AddIf(ctx, ctx.IsKindShown(), `kind`, Kind).
		AddIf(ctx, ctx.IsIndexShown(), `index`, p.index).
		AddNonZero(ctx2, `path`, p.path).
		AddNonZero(ctx2, `name`, p.name).
		AddNonZero(ctx2, `imports`, p.imports).
		AddNonZero(ctx2, `interfaces`, p.interfaces).
		AddNonZero(ctx2, `methods`, p.methods).
		AddNonZero(ctx2, `objects`, p.objects).
		AddNonZero(ctx2, `values`, p.values)
}

func (p *packageImp) resolveReceivers() {
	methods := p.methods.Values()
	for i := range methods.Count() {
		m := methods.Get(i)
		if !m.needsReceiver() {
			continue
		}

		if rec := m.receiverName(); len(rec) > 0 {
			d := p.findDeclaration(rec)
			if d == nil {
				panic(terror.New(`failed to find receiver`).
					With(`name`, rec))
			}
			o, ok := d.(Object)
			if !ok {
				panic(terror.New(`receiver was not an object`).
					With(`name`, rec).
					WithType(`gotten type`, d).
					With(`gotten value`, d))
			}
			o.addMethod(m)
			m.setReceiver(o)
		}
	}
}

func (p *packageImp) removeImports(predicate func(Construct) bool) {
	p.imports.Remove(predicate)
}
