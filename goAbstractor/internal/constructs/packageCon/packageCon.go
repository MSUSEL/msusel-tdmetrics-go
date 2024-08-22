package packageCon

import (
	"github.com/Snow-Gremlin/goToolbox/collections"
	"github.com/Snow-Gremlin/goToolbox/collections/sortedSet"
	"github.com/Snow-Gremlin/goToolbox/comp"
	"github.com/Snow-Gremlin/goToolbox/terrors/terror"
	"golang.org/x/tools/go/packages"

	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/assert"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs/interfaceDecl"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs/kind"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs/method"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs/object"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs/value"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/jsonify"
)

type packageImp struct {
	pkg *packages.Package

	path        string
	name        string
	index       int
	importPaths []string

	imports    collections.SortedSet[constructs.Package]
	interfaces collections.SortedSet[constructs.InterfaceDecl]
	methods    collections.SortedSet[constructs.Method]
	objects    collections.SortedSet[constructs.Object]
	values     collections.SortedSet[constructs.Value]
}

func newPackage(args constructs.PackageArgs) constructs.Package {
	assert.ArgNotNil(`real type`, args.RealPkg)
	assert.ArgNotEmpty(`path`, args.Path)
	assert.ArgValidId(`name`, args.Name)

	return &packageImp{
		pkg:         args.RealPkg,
		path:        args.Path,
		name:        args.Name,
		importPaths: args.ImportPaths,
		imports:     sortedSet.New(Comparer()),
		interfaces:  sortedSet.New(interfaceDecl.Comparer()),
		methods:     sortedSet.New(method.Comparer()),
		objects:     sortedSet.New(object.Comparer()),
		values:      sortedSet.New(value.Comparer()),
	}
}

func (p *packageImp) IsPackage() {}

func (p *packageImp) Kind() kind.Kind    { return kind.Package }
func (p *packageImp) SetIndex(index int) { p.index = index }

func (p *packageImp) Source() *packages.Package { return p.pkg }
func (p *packageImp) Path() string              { return p.path }
func (p *packageImp) Name() string              { return p.name }
func (p *packageImp) ImportPaths() []string     { return p.importPaths }
func (p *packageImp) String() string            { return p.name }

func (p *packageImp) Imports() collections.ReadonlySortedSet[constructs.Package] {
	return p.imports.Readonly()
}

func (p *packageImp) InitCount() int {
	return p.methods.Enumerate().
		Where(func(m constructs.Method) bool { return m.IsInit() }).
		Count()
}

func (p *packageImp) AddImport(i constructs.Package) constructs.Package {
	v, _ := p.imports.TryAdd(i)
	return v
}

func (p *packageImp) AddInterfaceDecl(it constructs.InterfaceDecl) constructs.InterfaceDecl {
	v, _ := p.interfaces.TryAdd(it)
	return v
}

func (p *packageImp) AddMethod(m constructs.Method) constructs.Method {
	v, _ := p.methods.TryAdd(m)
	return v
}

func (p *packageImp) AddObject(d constructs.Object) constructs.Object {
	v, _ := p.objects.TryAdd(d)
	return v
}

func (p *packageImp) AddValue(v constructs.Value) constructs.Value {
	v2, _ := p.values.TryAdd(v)
	return v2
}

func (p *packageImp) Empty() bool {
	return p.interfaces.Enumerate().Empty() &&
		p.methods.Enumerate().Empty() &&
		p.objects.Enumerate().Empty() &&
		p.values.Enumerate().Empty()
}

func (p *packageImp) findInterfaceDecl(name string) (constructs.InterfaceDecl, bool) {
	return p.interfaces.Enumerate().
		Where(func(t constructs.InterfaceDecl) bool { return t.Name() == name }).
		First()
}

func (p *packageImp) findObject(name string) (constructs.Object, bool) {
	return p.objects.Enumerate().
		Where(func(t constructs.Object) bool { return t.Name() == name }).
		First()
}

func (p *packageImp) FindTypeDecl(name string) constructs.TypeDecl {
	if v, ok := p.findInterfaceDecl(name); ok {
		return v
	}
	if v, ok := p.findObject(name); ok {
		return v
	}
	return nil
}

func (d *packageImp) CompareTo(other constructs.Construct) int {
	return constructs.CompareTo[constructs.Package](d, other, Comparer())
}

func Comparer() comp.Comparer[constructs.Package] {
	return func(a, b constructs.Package) int {
		aImp, bImp := a.(*packageImp), b.(*packageImp)
		return comp.Or(
			comp.DefaultPend(aImp.path, bImp.path),
			comp.DefaultPend(aImp.name, bImp.name),
		)
	}
}

func (p *packageImp) ToJson(ctx *jsonify.Context) jsonify.Datum {
	if ctx.IsShort() {
		return jsonify.New(ctx, p.index)
	}

	ctx2 := ctx.HideKind().Short()
	return jsonify.NewMap().
		AddIf(ctx, ctx.IsKindShown(), `kind`, p.Kind()).
		AddIf(ctx, ctx.IsIndexShown(), `index`, p.index).
		AddNonZero(ctx2, `path`, p.path).
		AddNonZero(ctx2, `name`, p.name).
		AddNonZero(ctx2, `imports`, p.imports.ToSlice()).
		AddNonZero(ctx2, `interfaces`, p.interfaces.ToSlice()).
		AddNonZero(ctx2, `methods`, p.methods.ToSlice()).
		AddNonZero(ctx2, `objects`, p.objects.ToSlice()).
		AddNonZero(ctx2, `values`, p.values.ToSlice())
}

func (p *packageImp) ResolveReceivers() {
	methods := p.methods
	for i := range methods.Count() {
		m := methods.Get(i)
		if !m.NeedsReceiver() {
			continue
		}

		if rec := m.ReceiverName(); len(rec) > 0 {
			o, ok := p.findObject(rec)
			if !ok {
				panic(terror.New(`failed to find receiver object`).
					With(`name`, rec))
			}
			o.AddMethod(m)
			m.SetReceiver(o)
		}
	}
}
