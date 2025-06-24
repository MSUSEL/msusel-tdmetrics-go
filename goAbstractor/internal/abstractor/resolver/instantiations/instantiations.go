package instantiations

import (
	"go/types"

	"github.com/Snow-Gremlin/goToolbox/collections"
	"github.com/Snow-Gremlin/goToolbox/terrors/terror"
	"github.com/Snow-Gremlin/goToolbox/utils"

	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/abstractor/analyzer"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/abstractor/baker"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/abstractor/converter"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/abstractor/instantiator"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/abstractor/querier"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/logger"
)

type Instantiations interface {

	// ExpandInstantiations adds propagation of instances so that if an object
	// has a method added after the instance, the method also gets instances created.
	ExpandInstantiations() bool
}

func New(log *logger.Logger, querier *querier.Querier, proj constructs.Project) Instantiations {
	return &instantiationsImp{
		log:            log,
		querier:        querier,
		proj:           proj,
		bk:             baker.New(proj),
		typeCache:      map[any]any{},
		changed:        true,
		doneMetrics:    map[constructs.MethodInst]bool{},
		doneExpandInst: map[constructs.Object]bool{},
		donePointerRec: map[constructs.Object]bool{},
		doneNestedType: map[constructs.Method]bool{},
	}
}

type instantiationsImp struct {
	log       *logger.Logger
	querier   *querier.Querier
	proj      constructs.Project
	bk        baker.Baker
	typeCache map[any]any
	changed   bool

	doneMetrics    map[constructs.MethodInst]bool
	doneExpandInst map[constructs.Object]bool
	donePointerRec map[constructs.Object]bool
	doneNestedType map[constructs.Method]bool
}

func (in *instantiationsImp) ExpandInstantiations() bool {
	in.changed = false
	in.fillOutAllMetrics()
	in.expandAllInstantiations()
	in.fillOutAllPointerReceivers()
	in.expandAllNestedTypes()
	return in.changed
}

func (in *instantiationsImp) fillOutAllMetrics() {
	methodInsts := in.proj.MethodInsts()
	for i := range methodInsts.Count() {
		if mi := methodInsts.Get(i); !in.doneMetrics[mi] {
			in.doneMetrics[mi] = true
			in.changed = true
			in.fillOutMetrics(mi)
		}
	}
}

func (in *instantiationsImp) fillOutMetrics(mi constructs.MethodInst) {
	m := mi.Generic()
	curPkg := m.Package()
	node := m.Metrics().Node()
	tpReplacer := m.Metrics().TpReplacer()
	conv := converter.New(in.log, in.querier, in.bk, in.proj, curPkg, m, mi.InstanceTypes(), tpReplacer, in.typeCache)
	metrics := analyzer.Analyze(in.log, in.querier, in.proj, curPkg, in.bk, conv, node)
	mi.SetMetrics(metrics)
}

func (in *instantiationsImp) expandAllInstantiations() {
	objects := in.proj.Objects()
	for i := range objects.Count() {
		if obj := objects.Get(i); !in.doneExpandInst[obj] {
			in.doneExpandInst[obj] = true
			in.changed = true
			in.expandInstantiations(obj)
		}
	}
}

func (in *instantiationsImp) expandInstantiations(obj constructs.Object) {
	// Add the method instances to the object.
	methods := obj.Methods()
	for i := range methods.Count() {
		mIts := methods.Get(i).Instances()
		for j := range mIts.Count() {
			// Create an object instance using the type argument for this
			// method instance so that the method has a receiver for it.
			it := mIts.Get(j).InstanceTypes()
			instantiator.Object(in.log, in.proj, nil, obj, nil, it)
		}
	}

	// Now that all the instances were collected, expand the instances
	// by adding all the method instance for that object instance, creating
	// any method instance that is missing.
	its := obj.Instances()
	for i := range its.Count() {
		in.expandObjectInst(obj, its.Get(i))
	}
}

// expandObjectInst adds the given instance into each method if it doesn't
// exist in that method. Then updates methods and receivers for the instance.
func (in *instantiationsImp) expandObjectInst(obj constructs.Object, instance constructs.ObjectInst) {
	methods := obj.Methods()
	for i := range methods.Count() {
		method := methods.Get(i)
		con := instantiator.Method(in.log, in.proj, method, instance.InstanceTypes())
		if utils.IsNil(con) {
			panic(terror.New(`unable to instantiate method while expanding object`).
				With(`method`, method).
				With(`object`, obj).
				With(`instance`, instance))
		}
		methodInst := con.(constructs.MethodInst)
		methodInst.SetReceiver(instance)
		instance.AddMethod(methodInst)
	}
}

func (in *instantiationsImp) fillOutAllPointerReceivers() {
	objects := in.proj.Objects()
	for i := range objects.Count() {
		if obj := objects.Get(i); !in.donePointerRec[obj] {
			in.donePointerRec[obj] = true
			in.changed = true
			in.fillOutPointerReceivers(obj)
		}
	}
}

func (in *instantiationsImp) fillOutPointerReceivers(obj constructs.Object) {
	if hasPointerReceivers(obj) {
		ptr := in.bk.BakePointer()
		// create a pointer for the generic object.
		rt := types.NewPointer(obj.GoType())
		instantiator.InterfaceDecl(in.log, in.proj, rt, ptr, nil, []constructs.TypeDesc{obj})

		oIts := obj.Instances()
		for i := range oIts.Count() {
			oIt := oIts.Get(i)
			// create a pointer for the object interface.
			rt := types.NewPointer(oIt.GoType())
			instantiator.InterfaceDecl(in.log, in.proj, rt, ptr, nil, []constructs.TypeDesc{oIt})
		}
	}
}

func hasPointerReceivers(obj constructs.Object) bool {
	return obj.Methods().Enumerate().Any(constructs.Method.PointerRecv)
}

func (in *instantiationsImp) expandAllNestedTypes() {
	methods := in.proj.Methods()
	for i := range methods.Count() {
		if m := methods.Get(i); !in.doneNestedType[m] {
			in.doneNestedType[m] = true
			in.changed = true
			in.expandNestedTypes(m)
		}
	}
}

func (in *instantiationsImp) expandNestedTypes(method constructs.Method) {
	nestedObjs := findNestedTypes(method, in.proj.Objects())
	nestedIts := findNestedTypes(method, in.proj.InterfaceDecls())

	mIts := method.Instances()
	for i := range mIts.Count() {
		mIt := mIts.Get(i)
		implicitTypes := mIt.InstanceTypes()
		for _, obj := range nestedObjs {
			in.expandNestedObject(implicitTypes, obj)
		}
		for _, it := range nestedIts {
			in.expandNestedInterface(implicitTypes, it)
		}
	}
}

func (in *instantiationsImp) expandNestedObject(implicitTypes []constructs.TypeDesc, obj constructs.Object) {
	it := constructs.Cast[constructs.TypeDesc](obj.TypeParams())
	instantiator.Object(in.log, in.proj, obj.GoType(), obj, implicitTypes, it)

	instances := obj.Instances()
	for j := range instances.Count() {
		inst := instances.Get(j)
		instantiator.Object(in.log, in.proj, inst.GoType(), obj, implicitTypes, inst.InstanceTypes())
	}
}

func (in *instantiationsImp) expandNestedInterface(implicitTypes []constructs.TypeDesc, it constructs.InterfaceDecl) {
	tp := constructs.Cast[constructs.TypeDesc](it.TypeParams())
	instantiator.InterfaceDecl(in.log, in.proj, it.GoType(), it, implicitTypes, tp)

	instances := it.Instances()
	for j := range instances.Count() {
		inst := instances.Get(j)
		instantiator.InterfaceDecl(in.log, in.proj, inst.GoType(), it, implicitTypes, inst.InstanceTypes())
	}
}

func findNestedTypes[T constructs.Nestable](nest constructs.NestType, ts collections.ReadonlySortedSet[T]) []T {
	if utils.IsNil(nest) {
		return nil
	}
	nested := []T{}
	for i := range ts.Count() {
		t := ts.Get(i)
		if constructs.ComparerPend(t.Nest(), nest)() == 0 {
			nested = append(nested, t)
		}
	}
	return nested
}
