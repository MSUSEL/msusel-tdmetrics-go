package resolver

import (
	"github.com/Snow-Gremlin/goToolbox/collections"
	"github.com/Snow-Gremlin/goToolbox/collections/enumerator"
	"github.com/Snow-Gremlin/goToolbox/terrors/terror"
	"github.com/Snow-Gremlin/goToolbox/utils"

	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/abstractor/instantiator"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/abstractor/resolver/inheritance"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs/interfaceDesc"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs/kind"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/logger"
)

type resolverImp struct {
	log  *logger.Logger
	proj constructs.Project
}

func Resolve(proj constructs.Project, log *logger.Logger) {
	resolve := &resolverImp{
		log:  log,
		proj: proj,
	}
	resolve.Imports()
	resolve.Receivers()
	resolve.ExpandInstantiations()
	resolve.TempReferences()
	resolve.ObjectInterfaces()
	resolve.Inheritance()
	resolve.DeadCodeElimination()
	resolve.Locations()
	resolve.Indices()
}

func (r *resolverImp) Imports() {
	r.log.Log(`resolve imports`)
	packages := r.proj.Packages()
	for i := range packages.Count() {
		pkg := packages.Get(i)
		for _, importPath := range pkg.ImportPaths() {
			impPackage := r.proj.FindPackageByPath(importPath)
			if impPackage == nil {
				panic(terror.New(`import package not found`).
					With(`package path`, pkg.Path).
					With(`import path`, importPath))
			}
			pkg.AddImport(impPackage)
		}
	}
}

func (r *resolverImp) Receivers() {
	r.log.Log(`resolve receivers`)
	packages := r.proj.Packages()
	for i := range packages.Count() {
		packages.Get(i).ResolveReceivers()
	}
}

// ExpandInstantiations adds propagation of instances so that if an object
// has a method added after the instance, the method also gets instances created.
func (r *resolverImp) ExpandInstantiations() {
	r.log.Log(`expand instantiations`)
	objects := r.proj.Objects()
	for i := range objects.Count() {
		r.expandInstantiations(objects.Get(i))
	}
}

func (r *resolverImp) expandInstantiations(obj constructs.Object) {
	// Add the method instances to the object.
	methods := obj.Methods()
	for i := range methods.Count() {
		mIts := methods.Get(i).Instances()
		for j := range mIts.Count() {
			it := mIts.Get(j).InstanceTypes()
			instantiator.Object(r.proj, nil, obj, it...)
		}
	}

	its := obj.Instances()
	for i := range its.Count() {
		r.expandObjectInst(obj, its.Get(i))
	}
}

// expandObjectInst adds the given instance into each method if it doesn't
// exist in that method. Then update methods and receivers for the instance.
func (r *resolverImp) expandObjectInst(obj constructs.Object, instance constructs.ObjectInst) {
	methods := obj.Methods()
	for i := range methods.Count() {
		method := methods.Get(i)
		con := instantiator.Method(r.proj, method, instance.InstanceTypes()...)
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

func (r *resolverImp) ObjectInterfaces() {
	r.log.Log(`resolve object interfaces`)
	log2 := r.log.Group(`objectInterfaces`).Prefix(`  `)
	log3 := log2.Prefix(`  `)

	objects := r.proj.Objects()
	for i := range objects.Count() {
		obj := objects.Get(i)
		if utils.IsNil(obj.Interface()) {
			log2.Logf(`%d) %s.%s`, i, obj.Package().Path(), obj.Name())
			r.objectInter(obj)
		}

		insts := obj.Instances()
		for j := range insts.Count() {
			it := insts.Get(j)
			if utils.IsNil(it.Interface()) {
				log3.Logf(`%d.%d) [%s]`, i, j, enumerator.Enumerate(it.InstanceTypes()...).Join(`, `))
				r.objectInstanceInter(it)
			}
		}
	}
}

func (r *resolverImp) objectInter(obj constructs.Object) {
	methods := obj.Methods()
	abstracts := make([]constructs.Abstract, methods.Count())
	for i := range methods.Count() {
		method := methods.Get(i)
		abstracts[i] = r.proj.NewAbstract(constructs.AbstractArgs{
			Name:      method.Name(),
			Signature: method.Signature(),
		})
	}

	it := r.proj.NewInterfaceDesc(constructs.InterfaceDescArgs{
		Abstracts: abstracts,
		Package:   obj.Package().Source(),
	})
	obj.SetInterface(it)
}

func (r *resolverImp) objectInstanceInter(objInst constructs.ObjectInst) {
	methodInsts := objInst.Methods()
	abstracts := make([]constructs.Abstract, methodInsts.Count())
	for i := range methodInsts.Count() {
		mi := methodInsts.Get(i)
		abstracts[i] = r.proj.NewAbstract(constructs.AbstractArgs{
			Name:      mi.Generic().Name(),
			Signature: mi.Resolved(),
		})
	}

	it := r.proj.NewInterfaceDesc(constructs.InterfaceDescArgs{
		Abstracts: abstracts,
		Package:   objInst.Generic().Package().Source(),
	})
	objInst.SetInterface(it)
}

func (r *resolverImp) Inheritance() {
	r.log.Log(`resolve inheritance`)
	its := r.proj.InterfaceDescs()
	log2 := r.log.Group(`inheritance`).Prefix(`  `)
	in := inheritance.New(interfaceDesc.Comparer(), log2)
	for i := range its.Count() {
		in.Process(its.Get(i))
	}
	log2.Log()
}

func (r *resolverImp) TempReferences() {
	r.log.Log(`resolve references`)
	refs := r.proj.TempReferences()
	for i := range refs.Count() {
		r.resolveTempRef(refs.Get(i))
	}

	r.proj.AllConstructs().Foreach(func(c constructs.Construct) {
		if trc, has := c.(constructs.TempReferenceContainer); has {
			trc.RemoveTempReferences()
		}
	})
	r.proj.ClearAllTempReferences()
}

func (r *resolverImp) resolveTempRef(ref constructs.TempReference) {
	if ref.Resolved() {
		return
	}

	_, typ, ok := r.proj.FindType(ref.PackagePath(), ref.Name(), true)
	if !ok {
		panic(terror.New(`failed to find temp referenced object`).
			With(`package path`, ref.PackagePath()).
			With(`name`, ref.Name()).
			With(`instance types`, ref.InstanceTypes()))
	}
	if len(ref.InstanceTypes()) <= 0 {
		ref.SetType(typ)
		return
	}

	switch typ.Kind() {
	case kind.Object:
		ref.SetType(instantiator.Object(r.proj, ref.GoType(), typ.(constructs.Object), ref.InstanceTypes()...))
	case kind.InterfaceDecl:
		ref.SetType(instantiator.InterfaceDecl(r.proj, ref.GoType(), typ.(constructs.InterfaceDecl), ref.InstanceTypes()...))
	default:
		panic(terror.New(`unexpected declaration type`).
			With(`kind`, typ.Kind()).
			With(`decl`, typ))
	}
}

func (r *resolverImp) DeadCodeElimination() {
	r.log.Log(`dead-code elimination`)
	/*
		needsUpdate := []constructs.Construct{}

		itDecls := r.proj.InterfaceDecls()
		for i := itDecls.Count() - 1; i >= 0; i-- {
			it := itDecls.Get(i)

		}
	*/

	// TODO: Improve prune to use metrics to create a dead code elimination prune.
}

func (r *resolverImp) Locations() {
	r.log.Log(`resolve locations`)
	r.proj.Locs().Reset()
	flagList(r.proj.InterfaceDecls())
	flagList(r.proj.Methods())
	flagList(r.proj.Objects())
	flagList(r.proj.Values())
}

func flagList[T constructs.Declaration](c collections.ReadonlySortedSet[T]) {
	for i := range c.Count() {
		c.Get(i).Location().Flag()
	}
}

// Indices should be called after all types have been registered
// and all packages have been processed. This will update all the indices
// that will be used as references in the output models.
func (r *resolverImp) Indices() {
	r.log.Log(`resolve indices`)
	var index int
	var kind kind.Kind
	r.proj.AllConstructs().Foreach(func(c constructs.Construct) {
		if cKind := c.Kind(); kind != cKind {
			kind = cKind
			index = 0
		}
		index++
		c.SetIndex(index)
	})
}
