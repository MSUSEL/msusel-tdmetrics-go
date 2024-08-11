package packageCon

import (
	"strings"

	"github.com/Snow-Gremlin/goToolbox/collections"
	"github.com/Snow-Gremlin/goToolbox/collections/enumerator"
	"github.com/Snow-Gremlin/goToolbox/terrors/terror"
	"golang.org/x/tools/go/packages"

	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/assert"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs/declaration"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/jsonify"
)

const Kind = `package`

type PackageCon interface {
	constructs.Construct
	_package()

	Source() *packages.Package
	Path() string
	Name() string
	ImportPaths() []string
	Imports() collections.ReadonlyList[PackageCon]
	InitCount() int

	addImport(p Package) Package
	addInterface(it Interface) Interface
	addMethod(m Method) Method
	addObject(id Object) Object
	addValue(v Value) Value

	empty() bool
	findDeclaration(name string) declaration.Declaration
	allDeclarations() collections.Enumerator[declaration.Declaration]

	resolveReceivers()
	removeImports(predicate func(constructs.Construct) bool)
}

type Args struct {
	RealPkg     *packages.Package
	Path        string
	Name        string
	ImportPaths []string
}

type packageConImp struct {
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

func newPackage(args Args) PackageCon {
	assert.ArgNotNil(`real type`, args.RealPkg)
	assert.ArgNotEmpty(`path`, args.Path)
	assert.ArgValidId(`name`, args.Name)

	return &packageConImp{
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

func (p *packageConImp) _package()          {}
func (p *packageConImp) Kind() string       { return Kind }
func (p *packageConImp) SetIndex(index int) { p.index = index }

func (p *packageConImp) Source() *packages.Package { return p.pkg }
func (p *packageConImp) Path() string              { return p.path }
func (p *packageConImp) Name() string              { return p.name }
func (p *packageConImp) ImportPaths() []string     { return p.importPaths }

func (p *packageConImp) Imports() collections.ReadonlyList[PackageCon] {
	return p.imports.Values()
}

func (p *packageConImp) InitCount() int {
	count := 0
	methods := p.methods.Values()
	for i := range methods.Count() {
		if methods.Get(i).IsInit() {
			count++
		}
	}
	return count
}

func (p *packageConImp) addImport(i Package) Package {
	return p.imports.Insert(i)
}

func (p *packageConImp) addInterface(it Interface) Interface {
	return p.interfaces.Insert(it)
}

func (p *packageConImp) addMethod(m Method) Method {
	return p.methods.Insert(m)
}

func (p *packageConImp) addObject(d Object) Object {
	return p.objects.Insert(d)
}

func (p *packageConImp) addValue(v Value) Value {
	return p.values.Insert(v)
}

func (p *packageConImp) empty() bool {
	return p.allDeclarations().Empty()
}

func (p *packageConImp) findDeclaration(name string) declaration.Declaration {
	def, _ := p.allDeclarations().
		Where(func(t declaration.Declaration) bool { return t.Name() == name }).
		First()
	return def
}

func (p *packageConImp) allDeclarations() collections.Enumerator[declaration.Declaration] {
	i := enumerator.Cast[declaration.Declaration](p.interfaces.Values().Enumerate())
	m := enumerator.Cast[declaration.Declaration](p.methods.Values().Enumerate())
	o := enumerator.Cast[declaration.Declaration](p.objects.Values().Enumerate())
	v := enumerator.Cast[declaration.Declaration](p.values.Values().Enumerate())
	return i.Concat(m).Concat(o).Concat(v)
}

func (p *packageConImp) compareTo(other constructs.Construct) int {
	b := other.(*packageConImp)
	return or(
		func() int { return strings.Compare(p.path, b.path) },
		func() int { return strings.Compare(p.name, b.name) },
	)
}

func (p *packageConImp) ToJson(ctx *jsonify.Context) jsonify.Datum {
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

func (p *packageConImp) resolveReceivers() {
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

func (p *packageConImp) removeImports(predicate func(constructs.Construct) bool) {
	p.imports.Remove(predicate)
}
