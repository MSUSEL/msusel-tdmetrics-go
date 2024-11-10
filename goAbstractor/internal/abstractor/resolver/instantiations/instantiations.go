package instantiations

import (
	"github.com/Snow-Gremlin/goToolbox/terrors/terror"
	"github.com/Snow-Gremlin/goToolbox/utils"

	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/abstractor/instantiator"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs"
)

// ExpandInstantiations adds propagation of instances so that if an object
// has a method added after the instance, the method also gets instances created.
func ExpandInstantiations(proj constructs.Project) {
	objects := proj.Objects()
	for i := range objects.Count() {
		expandInstantiations(proj, objects.Get(i))
	}
}

func expandInstantiations(proj constructs.Project, obj constructs.Object) {
	// Add the method instances to the object.
	methods := obj.Methods()
	for i := range methods.Count() {
		mIts := methods.Get(i).Instances()
		for j := range mIts.Count() {
			it := mIts.Get(j).InstanceTypes()
			instantiator.Object(proj, nil, obj, it...)
		}
	}

	// Now that all the instances were collected, expand the instances.
	its := obj.Instances()
	for i := range its.Count() {
		expandObjectInst(proj, obj, its.Get(i))
	}
}

// expandObjectInst adds the given instance into each method if it doesn't
// exist in that method. Then update methods and receivers for the instance.
func expandObjectInst(proj constructs.Project, obj constructs.Object, instance constructs.ObjectInst) {
	methods := obj.Methods()
	for i := range methods.Count() {
		method := methods.Get(i)
		con := instantiator.Method(proj, method, instance.InstanceTypes()...)
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
