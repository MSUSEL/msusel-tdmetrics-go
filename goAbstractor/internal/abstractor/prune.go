package abstractor

import (
	"slices"

	"github.com/Snow-Gremlin/goToolbox/utils"

	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs"
)

func (ab *abstractor) prune() {
	ab.pruneTypes()
	ab.prunePackages()
}

func (ab *abstractor) pruneTypes() {
	touched := map[constructs.Visitable]bool{
		ab.bakeAny(): true,
	}

	visitor := func(value constructs.Visitable) bool {
		if touched[value] {
			return false
		}
		touched[value] = true
		return true
	}

	// Visit everything reachable from the packages.
	// Do not visit the registered types since they are being pruned.
	for _, pkg := range ab.proj.Packages() {
		if !utils.IsNil(pkg) && visitor(pkg) {
			pkg.Visit(visitor)
		}
	}

	ab.proj.Types().Remove(func(td constructs.TypeDesc) bool {
		return !touched[td]
	})

	for _, pkg := range ab.proj.Packages() {
		pkg.Prune(func(v any) bool {
			return !touched[v.(constructs.Visitable)]
		})
	}
}

func (ab *abstractor) prunePackages() {
	empty := map[constructs.Package]bool{}
	for _, p := range ab.proj.Packages() {
		if p.Empty() {
			empty[p] = true
		}
	}

	ab.proj.FilterPackage(func(pkg constructs.Package) bool {
		return empty[pkg]
	})

	for _, p := range ab.proj.Packages() {
		p.SetImports(slices.DeleteFunc(p.Imports(), func(pkg constructs.Package) bool {
			return empty[pkg]
		}))
	}
}
