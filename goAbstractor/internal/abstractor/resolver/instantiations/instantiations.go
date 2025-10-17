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
	for mi := range in.proj.MethodInsts().Enumerate().Seq() {
		if !in.doneMetrics[mi] {
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
	for obj := range in.proj.Objects().Enumerate().Seq() {
		if !in.doneExpandInst[obj] {
			in.doneExpandInst[obj] = true
			in.changed = true
			in.expandInstantiations(obj)
		}
	}
}

func (in *instantiationsImp) expandInstantiations(obj constructs.Object) {
	// Add the method instances to the object.
	for method := range obj.Methods().Enumerate().Seq() {
		for mIt := range method.Instances().Enumerate().Seq() {
			// Create an object instance using the type argument for this
			// method instance so that the method has a receiver for it.
			it := mIt.InstanceTypes()
			instantiator.Object(in.log, in.querier, in.proj, nil, obj, nil, it)
		}
	}

	// Now that all the instances were collected, expand the instances
	// by adding all the method instance for that object instance, creating
	// any method instance that is missing.
	for it := range obj.Instances().Enumerate().Seq() {
		in.expandObjectInst(obj, it)
	}
}

// expandObjectInst adds the given instance into each method if it doesn't
// exist in that method. Then updates methods and receivers for the instance.
func (in *instantiationsImp) expandObjectInst(obj constructs.Object, instance constructs.ObjectInst) {
	for method := range obj.Methods().Enumerate().Seq() {
		con := instantiator.Method(in.log, in.querier, in.proj, method, instance.InstanceTypes())
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
	for obj := range in.proj.Objects().Enumerate().Seq() {
		if !in.donePointerRec[obj] {
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
		p := instantiator.InterfaceDecl(in.log, in.querier, in.proj, rt, ptr, nil, []constructs.TypeDesc{obj})
		trySetInheritance(ptr, p)

		for oIt := range obj.Instances().Enumerate().Seq() {
			// create a pointer for the object interface.
			rt := types.NewPointer(oIt.GoType())
			c := instantiator.InterfaceDecl(in.log, in.querier, in.proj, rt, ptr, nil, []constructs.TypeDesc{oIt})
			trySetInheritance(p, c)
		}
	}
}

func hasPointerReceivers(obj constructs.Object) bool {
	return obj.Methods().Enumerate().Any(constructs.Method.PointerRecv)
}

func getInterfaceDesc(td constructs.TypeDesc) (constructs.InterfaceDesc, bool) {
	switch t := td.(type) {
	case constructs.InterfaceDecl:
		return t.Interface(), true
	case constructs.InterfaceInst:
		return t.Resolved(), true
	default:
		return nil, false
	}
}

func trySetInheritance(parent, child constructs.TypeDesc) {
	if dp, ok := getInterfaceDesc(parent); ok {
		if dc, ok := getInterfaceDesc(child); ok {
			dc.AddInherits(dp)
		}
	}
}

func (in *instantiationsImp) expandAllNestedTypes() {
	for m := range in.proj.Methods().Enumerate().Seq() {
		if !in.doneNestedType[m] {
			in.doneNestedType[m] = true
			in.changed = true
			in.expandNestedTypes(m)
		}
	}
}

func (in *instantiationsImp) expandNestedTypes(method constructs.Method) {
	nestedObjs := findNestedTypes(method, in.proj.Objects())
	nestedIts := findNestedTypes(method, in.proj.InterfaceDecls())

	for mIt := range method.Instances().Enumerate().Seq() {
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
	instantiator.Object(in.log, in.querier, in.proj, obj.GoType(), obj, implicitTypes, it)

	for inst := range obj.Instances().Enumerate().Seq() {
		instantiator.Object(in.log, in.querier, in.proj, inst.GoType(), obj, implicitTypes, inst.InstanceTypes())
	}
}

func (in *instantiationsImp) expandNestedInterface(implicitTypes []constructs.TypeDesc, it constructs.InterfaceDecl) {
	tp := constructs.Cast[constructs.TypeDesc](it.TypeParams())
	p := instantiator.InterfaceDecl(in.log, in.querier, in.proj, it.GoType(), it, implicitTypes, tp)
	trySetInheritance(it, p)

	for inst := range it.Instances().Enumerate().Seq() {
		c := instantiator.InterfaceDecl(in.log, in.querier, in.proj, inst.GoType(), it, implicitTypes, inst.InstanceTypes())
		trySetInheritance(p, c)
	}
}

func findNestedTypes[T constructs.Nestable](nest constructs.NestType, ts collections.ReadonlySortedSet[T]) []T {
	if utils.IsNil(nest) {
		return nil
	}
	nested := []T{}
	for t := range ts.Enumerate().Seq() {
		if constructs.ComparerPend(t.Nest(), nest)() == 0 {
			nested = append(nested, t)
		}
	}
	return nested
}
