package dce

import (
	"github.com/Snow-Gremlin/goToolbox/collections"
	"github.com/Snow-Gremlin/goToolbox/collections/sortedSet"
	"github.com/Snow-Gremlin/goToolbox/comp"
	"github.com/Snow-Gremlin/goToolbox/utils"

	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/assert"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs/kind"
)

func DeadCodeElimination(proj constructs.Project) {
	d := &dce{
		proj:    proj,
		pending: sortedSet.New(comp.ComparableComparer[constructs.Construct]()),
	}

	d.primeAlive()

	for !d.pending.Empty() {
		if c := d.pending.TakeFirst(); !c.Alive() {
			c.SetAlive(true)
			d.updateAlive(c)
		}
	}
}

type dce struct {
	proj    constructs.Project
	pending collections.SortedSet[constructs.Construct]
}

func (d *dce) forcePend(c constructs.Construct) {
	if !utils.IsNil(c) {
		c.SetAlive(true)
		d.pending.Add(c)
	}
}

func (d *dce) pend(c constructs.Construct) {
	if !utils.IsNil(c) && !c.Alive() {
		c.SetAlive(true)
		d.pending.Add(c)
	}
}

func pendSlice[T constructs.Construct](d *dce, cs []T) {
	for _, c := range cs {
		d.pend(c)
	}
}

func pendSet[T constructs.Construct](d *dce, c collections.ReadonlySortedSet[T]) {
	for i := range c.Count() {
		d.pend(c.Get(i))
	}
}

func (d *dce) primeAlive() {
	entryPkg := d.proj.EntryPoint()
	assert.ArgNotNil(`entry point package`, entryPkg)
	d.forcePend(entryPkg)
	d.primeAliveGeneral()

	// Since this abstractor doesn't currently specify for tests,
	// if there are any tester functions, just make them alive.
	d.primeAliveForTests()

	// Check if the entry package has the main method.
	if entryPkg.Name() == `main` {
		main, ok := entryPkg.Methods().Enumerate().
			Where(func(m constructs.Method) bool { return m.IsMain() }).
			First()
		if ok {
			d.primeAliveForMain(main)
			return
		}
	}

	// If no main has been found, treat this like a library.
	d.primeAliveForLibrary(entryPkg)
}

func (d *dce) primeAliveGeneral() {
	d.proj.Methods().Enumerate().
		Where(func(m constructs.Method) bool { return m.IsInit() }).
		Foreach(func(m constructs.Method) { d.forcePend(m) })

	d.proj.Values().Enumerate().
		Where(func(v constructs.Value) bool { return v.HasSideEffect() }).
		Foreach(func(v constructs.Value) { d.forcePend(v) })
}

func (d *dce) primeAliveForTests() {
	d.proj.Methods().Enumerate().
		Where(func(m constructs.Method) bool { return m.IsTester() }).
		Foreach(func(m constructs.Method) { d.forcePend(m) })
}

func (d *dce) primeAliveForMain(main constructs.Method) {
	d.forcePend(main)
}

func (d *dce) primeAliveForLibrary(entryPkg constructs.Package) {
	d.primeAliveGeneral()
	entryPkg.InterfaceDecls().Enumerate().
		Where(func(it constructs.InterfaceDecl) bool { return it.Exported() }).
		Foreach(func(it constructs.InterfaceDecl) { d.forcePend(it) })

	entryPkg.Methods().Enumerate().
		Where(func(m constructs.Method) bool { return !m.HasReceiver() && m.Exported() }).
		Foreach(func(m constructs.Method) { d.forcePend(m) })

	entryPkg.Objects().Enumerate().
		Where(func(obj constructs.Object) bool { return obj.Exported() }).
		Foreach(func(obj constructs.Object) { d.forcePend(obj) })

	entryPkg.Values().Enumerate().
		Where(func(v constructs.Value) bool { return v.Exported() }).
		Foreach(func(v constructs.Value) { d.forcePend(v) })
}

func (d *dce) updateAlive(c constructs.Construct) {
	switch c.Kind() {
	case kind.Abstract:
		d.updateAbstract(c.(constructs.Abstract))
	case kind.Argument:
		d.updateArgument(c.(constructs.Argument))
	case kind.Basic:
		d.updateBasic(c.(constructs.Basic))
	case kind.Field:
		d.updateField(c.(constructs.Field))
	case kind.InterfaceDecl:
		d.updateInterfaceDecl(c.(constructs.InterfaceDecl))
	case kind.InterfaceDesc:
		d.updateInterfaceDesc(c.(constructs.InterfaceDesc))
	case kind.InterfaceInst:
		d.updateInterfaceInst(c.(constructs.InterfaceInst))
	case kind.Method:
		d.updateMethod(c.(constructs.Method))
	case kind.MethodInst:
		d.updateMethodInst(c.(constructs.MethodInst))
	case kind.Metrics:
		d.updateMetrics(c.(constructs.Metrics))
	case kind.Object:
		d.updateObject(c.(constructs.Object))
	case kind.ObjectInst:
		d.updateObjectInst(c.(constructs.ObjectInst))
	case kind.Package:
		d.updatePackage(c.(constructs.Package))
	case kind.Selection:
		d.updateSelection(c.(constructs.Selection))
	case kind.Signature:
		d.updateSignature(c.(constructs.Signature))
	case kind.StructDesc:
		d.updateStructDesc(c.(constructs.StructDesc))
	case kind.TempDeclRef:
		d.updateTempDeclRef(c.(constructs.TempDeclRef))
	case kind.TempReference:
		d.updateTempReference(c.(constructs.TempReference))
	case kind.TempTypeParamRef:
		d.updateTempTypeParamRef(c.(constructs.TempTypeParamRef))
	case kind.TypeParam:
		d.updateTypeParam(c.(constructs.TypeParam))
	case kind.Value:
		d.updateValue(c.(constructs.Value))
	}
}

func (d *dce) updateBasic(c constructs.Basic)     {}
func (d *dce) updatePackage(c constructs.Package) {}

func (d *dce) updateAbstract(c constructs.Abstract)   { d.pend(c.Signature()) }
func (d *dce) updateArgument(c constructs.Argument)   { d.pend(c.Type()) }
func (d *dce) updateField(c constructs.Field)         { d.pend(c.Type()) }
func (d *dce) updateTypeParam(c constructs.TypeParam) { d.pend(c.Type()) }

func (d *dce) updateStructDesc(c constructs.StructDesc) { pendSlice(d, c.Fields()) }

func (d *dce) updateTempDeclRef(c constructs.TempDeclRef)           { d.pend(c.ResolvedType()) }
func (d *dce) updateTempReference(c constructs.TempReference)       { d.pend(c.ResolvedType()) }
func (d *dce) updateTempTypeParamRef(c constructs.TempTypeParamRef) { d.pend(c.ResolvedType()) }

func (d *dce) updateInterfaceDecl(c constructs.InterfaceDecl) {
	d.pend(c.Package())
	d.pend(c.Interface())
	d.pend(c.Nest())
	pendSlice(d, c.TypeParams())
	pendSlice(d, c.ImplicitTypeParams())
	// Do not automatically make instances alive for the alive generics.
}

func (d *dce) updateMethod(c constructs.Method) {
	d.pend(c.Package())
	d.pend(c.Receiver())
	d.pend(c.Signature())
	pendSlice(d, c.TypeParams())
	// Do not automatically make instances alive for the alive generics.
}

func (d *dce) updateObject(c constructs.Object) {
	d.pend(c.Package())
	d.pend(c.Interface())
	d.pend(c.Nest())
	pendSlice(d, c.TypeParams())
	pendSlice(d, c.ImplicitTypeParams())
	// Do not automatically make instances alive for the alive generics.
}

func (d *dce) updateInterfaceInst(c constructs.InterfaceInst) {
	d.pend(c.Generic())
	d.pend(c.Resolved())
	pendSlice(d, c.InstanceTypes())
	pendSlice(d, c.ImplicitTypes())
}

func (d *dce) updateMethodInst(c constructs.MethodInst) {
	d.pend(c.Generic())
	d.pend(c.Resolved())
	pendSlice(d, c.InstanceTypes())
}

func (d *dce) updateObjectInst(c constructs.ObjectInst) {
	d.pend(c.Generic())
	d.pend(c.ResolvedData())
	d.pend(c.ResolvedInterface())
	pendSlice(d, c.InstanceTypes())
	pendSlice(d, c.ImplicitTypes())
}

func (d *dce) updateInterfaceDesc(c constructs.InterfaceDesc) {
	d.pend(c.PinnedPackage())
	pendSlice(d, c.Abstracts())
	pendSlice(d, c.Approx())
	pendSlice(d, c.Exact())
}

func (d *dce) updateMetrics(c constructs.Metrics) {
	pendSet(d, c.Invokes())
	pendSet(d, c.Reads())
	pendSet(d, c.Writes())
}

func (d *dce) updateSelection(c constructs.Selection) {
	d.pend(c.Origin())
	d.pend(c.Target())
}

func (d *dce) updateSignature(c constructs.Signature) {
	pendSlice(d, c.Params())
	pendSlice(d, c.Results())
}

func (d *dce) updateValue(c constructs.Value) {
	d.pend(c.Package())
	d.pend(c.Metrics())
	d.pend(c.Type())
}
