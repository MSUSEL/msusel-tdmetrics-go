package resolver

import (
	"github.com/Snow-Gremlin/goToolbox/collections"
	"github.com/Snow-Gremlin/goToolbox/collections/sortedSet"
	"github.com/Snow-Gremlin/goToolbox/terrors/terror"

	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs/interfaceDesc"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/logger"
)

type resolverImp struct {
	log  *logger.Logger
	proj constructs.Project
}

func Resolve(log *logger.Logger, proj constructs.Project) {
	resolve := &resolverImp{log: log, proj: proj}
	resolve.Imports()
	resolve.Receivers()
	resolve.ObjectInterfaces()
	resolve.Inheritance()
	resolve.References()
	resolve.EliminateDeadCode()
	resolve.Locations()
	resolve.UpdateIndices()
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

func (r *resolverImp) References() {
	r.log.Log(`resolve references`)
	refs := r.proj.References()
	for i := range refs.Count() {
		r.resolveRef(refs.Get(i))
	}
}

func (r *resolverImp) resolveRef(ref constructs.Reference) {
	if ref.Resolved() {
		return
	}

	if _, typ, ok := r.proj.FindType(ref.PackagePath(), ref.Name(), true); ok {

		// TODO: Handle type parameters to find instance

		ref.SetType(typ)
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

// UpdateIndices should be called after all types have been registered
// and all packages have been processed. This will update all the index
// fields that will be used as references in the output models.
func (r *resolverImp) UpdateIndices() {
	r.log.Log(`update indices`)
	// Type indices compound so that each has a unique offset.
	index := 1
	index = updateIndices(r.proj.Abstracts(), index)
	index = updateIndices(r.proj.Arguments(), index)
	index = updateIndices(r.proj.Basics(), index)
	index = updateIndices(r.proj.Fields(), index)
	index = updateIndices(r.proj.Instances(), index)
	index = updateIndices(r.proj.InterfaceDecls(), index)
	index = updateIndices(r.proj.InterfaceDescs(), index)
	index = updateIndices(r.proj.Methods(), index)
	index = updateIndices(r.proj.Objects(), index)
	index = updateIndices(r.proj.Packages(), index)
	// Don't index the p.References()
	index = updateIndices(r.proj.Signatures(), index)
	index = updateIndices(r.proj.StructDescs(), index)
	index = updateIndices(r.proj.TypeParams(), index)
	updateIndices(r.proj.Values(), index)
}

func updateIndices[T constructs.Construct](col collections.ReadonlySortedSet[T], index int) int {
	for i := range col.Count() {
		col.Get(i).SetIndex(index)
		index++
	}
	return index
}
