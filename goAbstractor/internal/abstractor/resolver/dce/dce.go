package dce

import (
	"github.com/Snow-Gremlin/goToolbox/collections"
	"github.com/Snow-Gremlin/goToolbox/collections/sortedSet"

	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/assert"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs"
)

// TODO: Check if imported packages that aren't used still have inits called.

func DeadCodeElimination(proj constructs.Project) {
	d := &dce{
		proj:    proj,
		pending: sortedSet.New[constructs.Construct](),
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
	if !c.Alive() {
		d.pending.Add(c)
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

	// TODO: Finnish
}
