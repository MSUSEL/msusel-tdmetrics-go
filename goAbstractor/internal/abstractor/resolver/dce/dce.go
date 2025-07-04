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

// TODO: Check if imported packages that aren't used still have inits called.

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

func (d *dce) primeAlive() {
	entryPkg := d.proj.EntryPoint()
	assert.ArgNotNil(`entry point package`, entryPkg)
	d.pend(entryPkg)

	// Check if the entry package has the main method.
	// This abstractor will not include tests so the packages
	// are either for a main method or for a library.
	if entryPkg.Name() == `main` {
		main, ok := entryPkg.Methods().Enumerate().
			Where(func(m constructs.Method) bool { return m.IsMain() }).
			First()
		if ok {
			d.primeAliveWithMain(entryPkg, main)
			return
		}
	}
	d.primeAliveWithLibrary(entryPkg)
}

func (d *dce) primeAliveWithMain(entryPkg constructs.Package, main constructs.Method) {
	d.pend(main)

	entryPkg.Methods().Enumerate().
		Where(func(m constructs.Method) bool { return m.IsInit() }).
		Foreach(func(m constructs.Method) { d.pend(m) })
}

func (d *dce) primeAliveWithLibrary(entryPkg constructs.Package) {
	entryPkg.InterfaceDecls().Enumerate().
		Where(func(it constructs.InterfaceDecl) bool { return it.Exported() }).
		Foreach(func(it constructs.InterfaceDecl) { d.pend(it) })

	entryPkg.Methods().Enumerate().
		Where(func(m constructs.Method) bool {
			return (!m.HasReceiver() && m.Exported()) || m.IsInit()
		}).Foreach(func(m constructs.Method) { d.pend(m) })

	entryPkg.Objects().Enumerate().
		Where(func(obj constructs.Object) bool { return obj.Exported() }).
		Foreach(func(obj constructs.Object) { d.pend(obj) })

	entryPkg.Values().Enumerate().
		Where(func(v constructs.Value) bool {
			return v.Exported() || v.HasSideEffect()
		}).Foreach(func(v constructs.Value) { d.pend(v) })
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

func (d *dce) updateTempDeclRef(c constructs.TempDeclRef)           { d.pend(c.ResolvedType()) }
func (d *dce) updateTempReference(c constructs.TempReference)       { d.pend(c.ResolvedType()) }
func (d *dce) updateTempTypeParamRef(c constructs.TempTypeParamRef) { d.pend(c.ResolvedType()) }

func (d *dce) updateInterfaceDecl(c constructs.InterfaceDecl) {
	d.pend(c.Package())
	d.pend(c.Interface())
	d.pend(c.Nest())
	pendSlice(d, c.TypeParams())
	pendSlice(d, c.ImplicitTypeParams())
	// Do not make instances alive for the alive generics.
}

func (d *dce) updateMethod(c constructs.Method) {
	d.pend(c.Package())
	d.pend(c.Receiver())
	d.pend(c.Signature())
	pendSlice(d, c.TypeParams())
	// Do not make instances alive for the alive generics.
}

func (d *dce) updateObject(c constructs.Object) {
	d.pend(c.Package())
	d.pend(c.Interface())
	d.pend(c.Nest())
	pendSlice(d, c.TypeParams())
	pendSlice(d, c.ImplicitTypeParams())
	// Do not make instances alive for the alive generics.
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
	// TODO: Implement
}

func (d *dce) updateMetrics(c constructs.Metrics) {
	// TODO: Implement
}

func (d *dce) updateSelection(c constructs.Selection) {
	// TODO: Implement
}

func (d *dce) updateSignature(c constructs.Signature) {
	// TODO: Implement
}

func (d *dce) updateStructDesc(c constructs.StructDesc) {
	// TODO: Implement
}

func (d *dce) updateValue(c constructs.Value) {
	// TODO: Implement
}
