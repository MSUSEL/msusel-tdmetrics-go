package resolver

import (
	"fmt"

	"github.com/Snow-Gremlin/goToolbox/collections"
	"github.com/Snow-Gremlin/goToolbox/collections/sortedSet"
	"github.com/Snow-Gremlin/goToolbox/terrors/terror"

	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs/interfaceDesc"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs/kind"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/logger"
)

type Args struct {
	Log     *logger.Logger
	Project constructs.Project

	UseGlobalIndices bool
}

type resolverImp struct {
	log  *logger.Logger
	proj constructs.Project

	useGlobalIndices bool
}

func Resolve(args Args) {
	resolve := &resolverImp{
		log:  args.Log,
		proj: args.Project,

		useGlobalIndices: args.UseGlobalIndices,
	}
	resolve.Imports()
	resolve.Receivers()
	resolve.ExpandInstantiations()
	resolve.TempReferences()
	resolve.ObjectInterfaces()
	resolve.Inheritance()
	resolve.EliminateDeadCode()
	resolve.Locations()
	resolve.Identifiers()
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

	// TODO: Implement
}

func (r *resolverImp) ObjectInterfaces() {
	r.log.Log(`resolve object interfaces`)
	objects := r.proj.Objects()
	for i := range objects.Count() {
		r.objectInter(objects.Get(i))
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

func (r *resolverImp) Inheritance() {
	r.log.Log(`resolve inheritance`)
	its := r.proj.InterfaceDescs()
	roots := sortedSet.New(interfaceDesc.Comparer())
	log2 := r.log.Group(`inheritance`)
	log3 := log2.Prefix(`  `)
	for i := range its.Count() {
		log2.Logf(`--(%d): %s`, i, its.Get(i))
		addInheritance(roots, its.Get(i), log3)
	}
	// throw away roots, they are no longer needed.
}

func addInheritance(siblings collections.SortedSet[constructs.InterfaceDesc], it constructs.InterfaceDesc, log *logger.Logger) {
	log2 := log.Prefix(` |  `)
	add := siblings.Count() <= 0
	for i := siblings.Count() - 1; i >= 0; i-- {
		a := siblings.Get(i)
		switch {
		case a.Implements(it):
			// Yi <: X
			log.Logf(` |--(%d) Yi<:X %s`, i, a)
			addInheritance(a.Inherits(), it, log2)
		case it.Implements(a):
			// Yi :> X
			log.Logf(` |--(%d) Yi:>X %s`, i, a)
			it.AddInherits(a)
			siblings.RemoveRange(i, 1)
			add = true
		default:
			// Possible overlap, check for super-types in subtree.
			log.Logf(` |--(%d) else  %s`, i, a)
			seekInherits(a.Inherits(), it, log2)
			add = true
		}
	}
	if add {
		log.Log(` '-- add`)
		siblings.Add(it)
	} else {
		log.Log(` '-- no-op`)
	}
}

func seekInherits(siblings collections.SortedSet[constructs.InterfaceDesc], it constructs.InterfaceDesc, log *logger.Logger) {
	for i := siblings.Count() - 1; i >= 0; i-- {
		a := siblings.Get(i)
		if it.Implements(a) {
			log.Log(` - `, a)
			it.AddInherits(a)
		} else {
			seekInherits(a.Inherits(), it, log)
		}
	}
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

	if _, typ, ok := r.proj.FindType(ref.PackagePath(), ref.Name(), true); ok {
		if len(ref.InstanceTypes()) > 0 {
			if inst, found := findInstance(typ, ref.InstanceTypes()); found {
				ref.SetType(inst)
				return
			}
			panic(terror.New(`failed to find temp referenced instance`).
				With(`package path`, ref.PackagePath()).
				With(`name`, ref.Name()).
				With(`instance types`, ref.InstanceTypes()))
		}
		ref.SetType(typ)
		return
	}
	panic(terror.New(`failed to find temp referenced object`).
		With(`package path`, ref.PackagePath()).
		With(`name`, ref.Name()).
		With(`instance types`, ref.InstanceTypes()))
}

func findInstance(decl constructs.TypeDecl, instanceTypes []constructs.TypeDesc) (constructs.TypeDesc, bool) {
	switch decl.Kind() {
	case kind.Object:
		return decl.(constructs.Object).FindInstance(instanceTypes)
	case kind.InterfaceDecl:
		return decl.(constructs.InterfaceDecl).FindInstance(instanceTypes)
	default:
		panic(terror.New(`unexpected declaration type`).
			With(`kind`, decl.Kind()).
			With(`decl`, decl))
	}
}

func (r *resolverImp) EliminateDeadCode() {
	// TODO: Improve prune to use metrics to create a dead code elimination prune.
	//ab.log.Logln(`prune`)
	//proj.PruneTypes()
	//proj.PrunePackages()
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

// Identifiers should be called after all types have been registered
// and all packages have been processed. This will update all the identifiers
// that will be used as references in the output models.
func (r *resolverImp) Identifiers() {
	if r.useGlobalIndices {
		r.globalIndices()
	} else {
		r.kindLocalIndices()
	}
}

func (r *resolverImp) globalIndices() {
	r.log.Log(`resolve identifiers - global indices`)
	index := 0
	r.proj.AllConstructs().Foreach(func(c constructs.Construct) {
		if i, has := c.(constructs.Identifiable); has {
			index++
			i.SetId(index)
		}
	})
}

func (r *resolverImp) kindLocalIndices() {
	r.log.Log(`resolve identifiers - kind local indices`)
	const kindLocalFormat = `%s%d`
	var index int
	var kind kind.Kind
	r.proj.AllConstructs().Foreach(func(c constructs.Construct) {
		if i, has := c.(constructs.Identifiable); has {
			if cKind := c.Kind(); kind != cKind {
				kind = cKind
				index = 0
			}
			index++
			i.SetId(fmt.Sprintf(kindLocalFormat, kind, index))
		}
	})
}
