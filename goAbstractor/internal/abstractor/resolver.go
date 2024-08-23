package abstractor

import (
	"github.com/Snow-Gremlin/goToolbox/collections"
	"github.com/Snow-Gremlin/goToolbox/collections/sortedSet"
	"github.com/Snow-Gremlin/goToolbox/terrors/terror"

	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs/interfaceDesc"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/logger"
)

func (ab *abstractor) resolveImports() {
	ab.log.Log(`resolve imports`)
	packages := ab.proj.Packages()
	for i := range packages.Count() {
		pkg := packages.Get(i)
		for _, importPath := range pkg.ImportPaths() {
			impPackage := ab.proj.FindPackageByPath(importPath)
			if impPackage == nil {
				panic(terror.New(`import package not found`).
					With(`package path`, pkg.Path).
					With(`import path`, importPath))
			}
			pkg.AddImport(impPackage)
		}
	}
}

func (ab *abstractor) resolveReceivers() {
	ab.log.Log(`resolve receivers`)
	packages := ab.proj.Packages()
	for i := range packages.Count() {
		packages.Get(i).ResolveReceivers()
	}
}

func (ab *abstractor) resolveObjectInterfaces() {
	ab.log.Log(`resolve object interfaces`)
	objects := ab.proj.Objects()
	for i := range objects.Count() {
		ab.resolveObjectInter(objects.Get(i))
	}
}

func (ab *abstractor) resolveObjectInter(obj constructs.Object) {
	methods := obj.Methods()
	abstracts := make([]constructs.Abstract, methods.Count())
	for i := range methods.Count() {
		method := methods.Get(i)
		abstracts[i] = ab.proj.NewAbstract(constructs.AbstractArgs{
			Name:      method.Name(),
			Signature: method.Signature(),
		})
	}
	it := ab.proj.NewInterfaceDesc(constructs.InterfaceDescArgs{
		Abstracts: abstracts,
		Package:   obj.Package().Source(),
	})
	obj.SetInterface(it)
}

func (ab *abstractor) resolveInheritance() {
	ab.log.Log(`resolve inheritance`)
	its := ab.proj.InterfaceDescs()
	roots := sortedSet.New(interfaceDesc.Comparer())
	log2 := ab.log.Group(`inheritance`)
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

func (ab *abstractor) resolveReferences() {
	ab.log.Log(`resolve references`)
	refs := ab.proj.References()
	for i := range refs.Count() {
		ab.resolveReference(refs.Get(i))
	}
}

func (ab *abstractor) resolveReference(ref constructs.Reference) {
	if ref.Resolved() {
		return
	}

	if _, typ, ok := ab.proj.FindType(ref.PackagePath(), ref.Name(), true); ok {

		// TODO: Handle type parameters to find instance

		ref.SetType(typ)
	}
}

func (ab *abstractor) eliminateDeadCode() {
	// TODO: Improve prune to use metrics to create a dead code elimination prune.
	//ab.log.Logln(`prune`)
	//proj.PruneTypes()
	//proj.PrunePackages()
}

func (ab *abstractor) resolveLocations() {
	ab.log.Log(`resolve locations`)
	ab.locs.Reset()
	flagList(ab.proj.InterfaceDecls())
	flagList(ab.proj.Methods())
	flagList(ab.proj.Objects())
	flagList(ab.proj.Values())
}

func flagList[T constructs.Declaration](c collections.ReadonlySortedSet[T]) {
	for i := range c.Count() {
		c.Get(i).Location().Flag()
	}
}

// UpdateIndices should be called after all types have been registered
// and all packages have been processed. This will update all the index
// fields that will be used as references in the output models.
func (ab *abstractor) updateIndices() {
	ab.log.Log(`update indices`)
	// Type indices compound so that each has a unique offset.
	index := 1
	index = updateIndices(ab.proj.Abstracts(), index)
	index = updateIndices(ab.proj.Arguments(), index)
	index = updateIndices(ab.proj.Basics(), index)
	index = updateIndices(ab.proj.Fields(), index)
	index = updateIndices(ab.proj.Instances(), index)
	index = updateIndices(ab.proj.InterfaceDecls(), index)
	index = updateIndices(ab.proj.InterfaceDescs(), index)
	index = updateIndices(ab.proj.Methods(), index)
	index = updateIndices(ab.proj.Objects(), index)
	index = updateIndices(ab.proj.Packages(), index)
	// Don't index the p.References()
	index = updateIndices(ab.proj.Signatures(), index)
	index = updateIndices(ab.proj.StructDescs(), index)
	index = updateIndices(ab.proj.TypeParams(), index)
	updateIndices(ab.proj.Values(), index)
}

func updateIndices[T constructs.Construct](col collections.ReadonlySortedSet[T], index int) int {
	for i := range col.Count() {
		col.Get(i).SetIndex(index)
		index++
	}
	return index
}
