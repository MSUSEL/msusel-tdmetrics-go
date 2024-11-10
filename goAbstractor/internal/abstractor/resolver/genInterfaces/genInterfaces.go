package genInterfaces

import (
	"github.com/Snow-Gremlin/goToolbox/collections/enumerator"
	"github.com/Snow-Gremlin/goToolbox/utils"

	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/logger"
)

func GenerateInterfaces(log *logger.Logger, proj constructs.Project) {
	log = log.Group(`generateInterfaces`).Prefix(`  `)
	log2 := log.Prefix(`  `)

	// Resolve all objects interfaces
	objects := proj.Objects()
	for i := range objects.Count() {
		obj := objects.Get(i)

		// If the object doesn't have an interface create one and set it.
		if utils.IsNil(obj.Interface()) {
			log.Logf(`%d) %s.%s`, i, obj.Package().Path(), obj.Name())
			objectInter(proj, obj)
		}

		// Resolve all instances for the object
		insts := obj.Instances()
		for j := range insts.Count() {
			it := insts.Get(j)

			// If the instance doesn't have an interface create one and set it.
			if utils.IsNil(it.ResolvedInterface()) {
				log2.Logf(`%d.%d) [%s]`, i, j, enumerator.Enumerate(it.InstanceTypes()...).Join(`, `))
				objectInstanceInter(proj, it)
			}
		}
	}

	// TODO: If pointer interface, pull abstracts foreword,
	// i.e. a pointer to an object can call the objects methods.

}

func objectInter(proj constructs.Project, obj constructs.Object) {
	methods := obj.Methods()
	abstracts := make([]constructs.Abstract, methods.Count())
	for i := range methods.Count() {
		method := methods.Get(i)
		abstracts[i] = proj.NewAbstract(constructs.AbstractArgs{
			Name:      method.Name(),
			Exported:  method.Exported(),
			Signature: method.Signature(),
		})
	}

	it := proj.NewInterfaceDesc(constructs.InterfaceDescArgs{
		Abstracts: abstracts,
		Package:   obj.Package().Source(),
	})
	obj.SetInterface(it)
}

func objectInstanceInter(proj constructs.Project, objInst constructs.ObjectInst) {
	methodInsts := objInst.Methods()
	abstracts := make([]constructs.Abstract, methodInsts.Count())
	for i := range methodInsts.Count() {
		mi := methodInsts.Get(i)
		abstracts[i] = proj.NewAbstract(constructs.AbstractArgs{
			Name:      mi.Generic().Name(),
			Exported:  mi.Generic().Exported(),
			Signature: mi.Resolved(),
		})
	}

	it := proj.NewInterfaceDesc(constructs.InterfaceDescArgs{
		Abstracts: abstracts,
		Package:   objInst.Generic().Package().Source(),
	})
	objInst.SetResolvedInterface(it)
}
