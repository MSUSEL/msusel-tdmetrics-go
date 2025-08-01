package packageCon

import (
	"fmt"

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
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/stringer"
)

type packageImp struct {
	constructs.ConstructCore
	pkg *packages.Package

	path        string
	name        string
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

func (p *packageImp) Kind() kind.Kind { return kind.Package }
func (p *packageImp) Path() string    { return p.path }
func (p *packageImp) Name() string    { return p.name }

func (p *packageImp) Source() *packages.Package { return p.pkg }

func (p *packageImp) EntryPoint() bool {
	return p.pkg.PkgPath == `command-line-arguments`
}

func (p *packageImp) ImportPaths() []string { return p.importPaths }

func (p *packageImp) Imports() collections.ReadonlySortedSet[constructs.Package] {
	return p.imports.Readonly()
}

func (p *packageImp) InterfaceDecls() collections.ReadonlySortedSet[constructs.InterfaceDecl] {
	return p.interfaces.Readonly()
}

func (p *packageImp) Methods() collections.ReadonlySortedSet[constructs.Method] {
	return p.methods.Readonly()
}

func (p *packageImp) Objects() collections.ReadonlySortedSet[constructs.Object] {
	return p.objects.Readonly()
}

func (p *packageImp) Values() collections.ReadonlySortedSet[constructs.Value] {
	return p.values.Readonly()
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

func (p *packageImp) findInterfaceDecl(name string, nest constructs.NestType) (constructs.InterfaceDecl, bool) {
	return p.interfaces.Enumerate().
		Where(func(t constructs.InterfaceDecl) bool {
			return comp.Or(
				comp.DefaultPend(t.Name(), name),
				constructs.ComparerPend(t.Nest(), nest),
			) == 0
		}).
		First()
}

func (p *packageImp) findObject(name string, nest constructs.NestType) (constructs.Object, bool) {
	return p.objects.Enumerate().
		Where(func(t constructs.Object) bool {
			return comp.Or(
				comp.DefaultPend(t.Name(), name),
				constructs.ComparerPend(t.Nest(), nest),
			) == 0
		}).
		First()
}

func (p *packageImp) findMethod(name string) (constructs.Method, bool) {
	return p.methods.Enumerate().
		Where(func(t constructs.Method) bool {
			return !t.HasReceiver() && t.Name() == name
		}).
		First()
}

func (p *packageImp) findValue(name string) (constructs.Value, bool) {
	return p.values.Enumerate().
		Where(func(t constructs.Value) bool {
			return t.Name() == name
		}).
		First()
}

func (p *packageImp) FindTypeDecl(name string, nest constructs.NestType) constructs.TypeDecl {
	if v, ok := p.findInterfaceDecl(name, nest); ok {
		return v
	}
	if v, ok := p.findObject(name, nest); ok {
		return v
	}
	return nil
}

func (p *packageImp) FindDecl(name string, nest constructs.NestType) constructs.Declaration {
	if v, ok := p.findInterfaceDecl(name, nest); ok {
		return v
	}
	if v, ok := p.findObject(name, nest); ok {
		return v
	}
	if v, ok := p.findMethod(name); ok {
		//assert.ArgIsNil(`nest`, nest)
		return v
	}
	if v, ok := p.findValue(name); ok {
		//assert.ArgIsNil(`nest`, nest)
		return v
	}
	return nil
}

func (p *packageImp) DebugPrintDecl() {
	if !p.InterfaceDecls().Empty() {
		fmt.Printf("Interface:\n\t%s\n", p.InterfaceDecls().Enumerate().Join("\n\t"))
	}
	if !p.Objects().Empty() {
		fmt.Printf("Objects:\n\t%s\n", p.Objects().Enumerate().Join("\n\t"))
	}
	if !p.Methods().Empty() {
		fmt.Printf("Methods:\n\t%s\n", p.Methods().Enumerate().Join("\n\t"))
	}
	if !p.Values().Empty() {
		fmt.Printf("Values:\n\t%s\n", p.Values().Enumerate().Join("\n\t"))
	}
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
	if ctx.IsOnlyIndex() {
		return jsonify.New(ctx, p.Index())
	}
	if ctx.IsShort() {
		return jsonify.NewSprintf(`%s%d`, p.Kind(), p.Index())
	}
	if ctx.SkipDead() && !p.Alive() {
		return nil
	}
	if !ctx.KeepDuplicates() && p.Duplicate() {
		return nil
	}
	return jsonify.NewMap().
		AddIf(ctx, ctx.IsDebugKindIncluded(), `kind`, p.Kind()).
		AddIf(ctx, ctx.IsDebugIndexIncluded(), `index`, p.Index()).
		AddIf(ctx, ctx.IsDebugAliveIncluded(), `alive`, p.Alive()).
		Add(ctx, `path`, p.path).
		Add(ctx, `name`, p.name).
		AddNonZero(ctx.OnlyIndex(), `imports`, constructs.JsonSet(ctx.OnlyIndex(), p.imports.ToSlice())).
		AddNonZero(ctx.OnlyIndex(), `interfaces`, constructs.JsonSet(ctx.OnlyIndex(), p.interfaces.ToSlice())).
		AddNonZero(ctx.OnlyIndex(), `methods`, constructs.JsonSet(ctx.OnlyIndex(), p.methods.ToSlice())).
		AddNonZero(ctx.OnlyIndex(), `objects`, constructs.JsonSet(ctx.OnlyIndex(), p.objects.ToSlice())).
		AddNonZero(ctx.OnlyIndex(), `values`, constructs.JsonSet(ctx.OnlyIndex(), p.values.ToSlice()))
}

func (p *packageImp) ResolveReceivers() {
	methods := p.methods
	for i := range methods.Count() {
		m := methods.Get(i)
		if !m.NeedsReceiver() {
			continue
		}

		if rec := m.ReceiverName(); len(rec) > 0 {
			o, ok := p.findObject(rec, nil)
			if !ok {
				panic(terror.New(`failed to find receiver object`).
					With(`name`, rec))
			}
			o.AddMethod(m)
			m.SetReceiver(o)
		}
	}
}

func (p *packageImp) ToStringer(s stringer.Stringer) { s.WriteString(p.name) }

func (p *packageImp) String() string { return p.name }
