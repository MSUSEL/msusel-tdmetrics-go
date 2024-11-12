package genInterfaces

import (
	"github.com/Snow-Gremlin/goToolbox/collections/enumerator"
	"github.com/Snow-Gremlin/goToolbox/utils"

	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs/hint"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs/innate"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/logger"
)

func GenerateInterfaces(log *logger.Logger, proj constructs.Project) {
	log = log.Group(`generateInterfaces`).Prefix(`  `)
	objectInterfaces(log, proj)
	extendPointers(proj)
}

// objectInterfaces resolves all object interfaces
// and the interfaces for object instances.
func objectInterfaces(log *logger.Logger, proj constructs.Project) {
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
}

// objectInter creates the basic non-pointer interface for the given object.
// Not all callable methods will appear in this interface since pointer
// methods are left out. However this is the minimal required interface
// that the whole object can be cast to as a non-pointer instance.
// See extendingPointer.md in the docs folder.
func objectInter(proj constructs.Project, obj constructs.Object) {
	methods := obj.Methods()
	abstracts := make([]constructs.Abstract, 0, methods.Count())

	// Add all non-pointer methods.
	for i := range methods.Count() {
		method := methods.Get(i)
		if !method.PointerRecv() {
			abstract := proj.NewAbstract(constructs.AbstractArgs{
				Name:      method.Name(),
				Exported:  method.Exported(),
				Signature: method.Signature(),
			})
			abstracts = append(abstracts, abstract)
		}
	}

	// Add any methods for base data.
	synthAbs := getSyntheticAbstracts(obj.Data())
	abstracts = append(abstracts, synthAbs...)

	// Create an set the interfaces.
	it := proj.NewInterfaceDesc(constructs.InterfaceDescArgs{
		Abstracts: utils.RemoveZeros(abstracts),
		Package:   obj.Package().Source(),
	})
	obj.SetInterface(it)
}

// objectInstanceInter creates the basic non-pointer interface for
// the given instance of an object. This is similar to the interface
// created via `objectInter`.
func objectInstanceInter(proj constructs.Project, objInst constructs.ObjectInst) {
	methodInsts := objInst.Methods()
	abstracts := make([]constructs.Abstract, 0, methodInsts.Count())

	// Add all non-pointer methods.
	for i := range methodInsts.Count() {
		mi := methodInsts.Get(i)
		method := mi.Generic()
		if !method.PointerRecv() {
			abstract := proj.NewAbstract(constructs.AbstractArgs{
				Name:      method.Name(),
				Exported:  method.Exported(),
				Signature: mi.Resolved(),
			})
			abstracts = append(abstracts, abstract)
		}
	}

	// Add any methods for base data.
	synthAbs := getSyntheticAbstracts(objInst.ResolvedData())
	abstracts = append(abstracts, synthAbs...)

	// Create an set the interfaces.
	it := proj.NewInterfaceDesc(constructs.InterfaceDescArgs{
		Abstracts: utils.RemoveZeros(abstracts),
		Package:   objInst.Generic().Package().Source(),
	})
	objInst.SetResolvedInterface(it)
}

// getSyntheticAbstracts will determine if the structure containing the
// data for an object is synthetic and needs to pull some function into
// the interface, e.g. `type A []int` can call `len(A)` so `$len` from `[]int`
// needs to be pulled forward into `A`.
func getSyntheticAbstracts(st constructs.StructDesc) []constructs.Abstract {
	if !st.Synthetic() {
		return nil
	}

	baseData := st.Fields()[0].Type()
	var it constructs.InterfaceDesc
	switch bd := baseData.(type) {
	case constructs.InterfaceDesc:
		it = bd
	case constructs.InterfaceDecl:
		it = bd.Interface()
	case constructs.InterfaceInst:
		it = bd.Resolved()
	}

	if !utils.IsNil(it) && it.Hint() != hint.None {
		abs := []constructs.Abstract{}
		// There should only be innate methods (with `$`) created by the baker
		// or abstractor but filter to only those methods to ensure no methods
		// that were pulled forward into a pointer are picked up.
		for _, ab := range it.Abstracts() {
			if innate.Is(ab.Name()) {
				abs = append(abs, ab)
			}
		}
		return abs
	}
	return nil
}

func extendPointers(proj constructs.Project) {
	its := proj.InterfaceDescs()
	for i := range its.Count() {
		if it := its.Get(i); it.Hint() == hint.Pointer {
			derefSig := constructs.FindSigByName(it.Abstracts(), innate.Deref)
			if !utils.IsNil(derefSig) && len(derefSig.Results()) == 1 {
				extendPointer(proj, it, derefSig.Results()[0].Type())
			}
		}
	}
}

func extendPointer(proj constructs.Project, it constructs.InterfaceDesc, refType constructs.TypeDesc) {
	switch rt := refType.(type) {
	case constructs.Object:
		methods := rt.Methods()
		abstracts := make([]constructs.Abstract, methods.Count())
		for i := range methods.Count() {
			method := methods.Get(i)
			abstracts[i] = proj.NewAbstract(constructs.AbstractArgs{
				Name:      method.Name(),
				Exported:  method.Exported(),
				Signature: method.Signature(),
			})
		}
		it.SetAdditionalAbstracts(abstracts)

	case constructs.ObjectInst:
		methodInsts := rt.Methods()
		abstracts := make([]constructs.Abstract, methodInsts.Count())
		for i := range methodInsts.Count() {
			mi := methodInsts.Get(i)
			method := mi.Generic()
			abstracts[i] = proj.NewAbstract(constructs.AbstractArgs{
				Name:      method.Name(),
				Exported:  method.Exported(),
				Signature: mi.Resolved(),
			})
		}
		it.SetAdditionalAbstracts(abstracts)
	}
}
