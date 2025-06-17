package instantiations

import (
	"go/types"

	"github.com/Snow-Gremlin/goToolbox/collections"
	"github.com/Snow-Gremlin/goToolbox/terrors/terror"
	"github.com/Snow-Gremlin/goToolbox/utils"

	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/abstractor/baker"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/abstractor/instantiator"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/logger"
)

// ExpandInstantiations adds propagation of instances so that if an object
// has a method added after the instance, the method also gets instances created.
func ExpandInstantiations(log *logger.Logger, proj constructs.Project) {
	objects := proj.Objects()
	for i := range objects.Count() {
		expandInstantiations(log, proj, objects.Get(i))
	}

	bk := baker.New(proj)
	for i := range objects.Count() {
		fillOutPointerReceivers(log, bk, proj, objects.Get(i))
	}

	methods := proj.Methods()
	for i := range methods.Count() {
		expandNestedTypes(log, proj, methods.Get(i))
	}
}

func expandInstantiations(log *logger.Logger, proj constructs.Project, obj constructs.Object) {
	// Add the method instances to the object.
	methods := obj.Methods()
	for i := range methods.Count() {
		mIts := methods.Get(i).Instances()
		for j := range mIts.Count() {
			// Create an object instance using the type argument for this
			// method instance so that the method has a receiver for it.
			it := mIts.Get(j).InstanceTypes()
			instantiator.Object(log, proj, nil, obj, nil, it)
		}
	}

	// Now that all the instances were collected, expand the instances
	// by adding all the method instance for that object instance, creating
	// any method instance that is missing.
	its := obj.Instances()
	for i := range its.Count() {
		expandObjectInst(log, proj, obj, its.Get(i))
	}
}

// expandObjectInst adds the given instance into each method if it doesn't
// exist in that method. Then updates methods and receivers for the instance.
func expandObjectInst(log *logger.Logger, proj constructs.Project, obj constructs.Object, instance constructs.ObjectInst) {
	methods := obj.Methods()
	for i := range methods.Count() {
		method := methods.Get(i)
		con := instantiator.Method(log, proj, method, instance.InstanceTypes())
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

func fillOutPointerReceivers(log *logger.Logger, bk baker.Baker, proj constructs.Project, obj constructs.Object) {
	if hasPointerReceivers(obj) {
		ptr := bk.BakePointer()
		// create a pointer for the generic object.
		rt := types.NewPointer(obj.GoType())
		instantiator.InterfaceDecl(log, proj, rt, ptr, nil, []constructs.TypeDesc{obj})

		oIts := obj.Instances()
		for i := range oIts.Count() {
			oIt := oIts.Get(i)
			// create a pointer for the object interface.
			rt := types.NewPointer(oIt.GoType())
			instantiator.InterfaceDecl(log, proj, rt, ptr, nil, []constructs.TypeDesc{oIt})
		}
	}
}

func hasPointerReceivers(obj constructs.Object) bool {
	return obj.Methods().Enumerate().Any(constructs.Method.PointerRecv)
}

func expandNestedTypes(log *logger.Logger, proj constructs.Project, method constructs.Method) {
	nestedObjs := findNestedTypes(method, proj.Objects())
	nestedIts := findNestedTypes(method, proj.InterfaceDecls())

	mIts := method.Instances()
	for i := range mIts.Count() {
		mIt := mIts.Get(i)
		implicitTypes := mIt.InstanceTypes()

		for _, obj := range nestedObjs {
			it := constructs.Cast[constructs.TypeDesc](obj.TypeParams())
			instantiator.Object(log, proj, obj.GoType(), obj, implicitTypes, it)

			instances := obj.Instances()
			for j := range instances.Count() {
				inst := instances.Get(j)
				instantiator.Object(log, proj, inst.GoType(), obj, implicitTypes, inst.InstanceTypes())
			}
		}

		for _, it := range nestedIts {
			tp := constructs.Cast[constructs.TypeDesc](it.TypeParams())
			instantiator.InterfaceDecl(log, proj, it.GoType(), it, implicitTypes, tp)

			instances := it.Instances()
			for j := range instances.Count() {
				inst := instances.Get(j)
				instantiator.InterfaceDecl(log, proj, inst.GoType(), it, implicitTypes, inst.InstanceTypes())
			}
		}
	}
}

func findNestedTypes[T constructs.TypeDecl](nest constructs.NestType, ts collections.ReadonlySortedSet[T]) []T {
	if utils.IsNil(nest) {
		return nil
	}
	nested := []T{}
	for i := range ts.Count() {
		t := ts.Get(i)
		if n, ok := any(t).(interface{ Nest() constructs.NestType }); ok {
			if constructs.ComparerPend(n.Nest(), nest)() == 0 {
				nested = append(nested, t)
			}
		}
	}
	return nested
}
