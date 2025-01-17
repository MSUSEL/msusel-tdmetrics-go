package resolver

import (
	"github.com/Snow-Gremlin/goToolbox/collections"
	"github.com/Snow-Gremlin/goToolbox/terrors/terror"

	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/abstractor/resolver/dce"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/abstractor/resolver/genInterfaces"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/abstractor/resolver/inheritance"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/abstractor/resolver/instantiations"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/abstractor/resolver/references"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs/interfaceDesc"
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

	// Resolve imports of packages and receivers in methods.
	resolve.Imports()
	resolve.Receivers()

	// First pass of removing references.
	// This includes creating instances that were referenced in the metrics.
	resolve.References()

	// Fill out all instantiations of generic object, interface, and methods.
	// Also fill out all pointer receivers that are still not defined.
	resolve.ExpandInstantiations()

	// Second pass of removing references.
	// This takes care of any references that the instantiation had to make.
	// There should be none but doesn't hurt to check.
	resolve.References()

	// Determine interfaces for objects and object instances.
	// Also extend the interfaces for pointers.
	resolve.GenerateInterfaces()

	// Determine inheritance hierarchy to solidify duck-typing.
	resolve.Inheritance()

	// Remove anything that isn't needed.
	resolve.DeadCodeElimination()

	// Update the locations and indices to prepare for outputting.
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

func (r *resolverImp) ExpandInstantiations() {
	r.log.Log(`expand instantiations`)
	instantiations.ExpandInstantiations(r.log, r.proj)
}

func (r *resolverImp) GenerateInterfaces() {
	r.log.Log(`generate interfaces`)
	genInterfaces.GenerateInterfaces(r.log, r.proj)
}

func (r *resolverImp) Inheritance() {
	r.log.Log(`resolve inheritance`)
	inheritance.Resolve(r.log, interfaceDesc.Comparer(), r.proj.InterfaceDescs())
}

func (r *resolverImp) References() {
	r.log.Log(`resolve references`)
	references.References(r.log, r.proj)
}

func (r *resolverImp) DeadCodeElimination() {
	r.log.Log(`dead-code elimination`)
	dce.DeadCodeElimination(r.proj)
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
	r.proj.UpdateIndices()
}
