package constructs

import (
	"slices"

	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/jsonify"
)

type Project interface {
	Types() Register
	ToJson(ctx *jsonify.Context) jsonify.Datum
	Packages() []Package
	AppendPackage(pkg ...Package)
	FilterPackage(predicate func(pkg Package) bool)

	// UpdateIndices should be called after all types have been registered
	// and all packages have been processed. This will update all the index
	// fields that will be used as references in the output models.
	UpdateIndices()
	Prune(keep ...TypeDesc)
}

type projectImp struct {
	allPackages []Package
	allTypes    Register
}

func NewProject() Project {
	return &projectImp{
		allTypes: NewRegister(),
	}
}

func (p *projectImp) Types() Register {
	return p.allTypes
}

func (p *projectImp) ToJson(ctx *jsonify.Context) jsonify.Datum {
	ctx2 := ctx.HideKind()
	return jsonify.NewMap().
		Add(ctx2, `language`, `go`).
		AddNonZero(ctx2, `types`, p.allTypes).
		AddNonZero(ctx2, `packages`, p.allPackages)
}

func (p *projectImp) Prune(keep ...TypeDesc) {
	p.pruneTypes(keep)
	p.prunePackages()
}

func (p *projectImp) pruneTypes(keep []TypeDesc) {
	touched := map[Visitable]bool{}
	for _, k := range keep {
		touched[k] = true
	}

	// Visit everything reachable from the packages.
	// Do not visit the registered types since they are being pruned.
	visitList(func(value Visitable) bool {
		if touched[value] {
			return false
		}
		touched[value] = true
		return true
	}, p.allPackages)

	p.Types().Remove(func(td TypeDesc) bool {
		return !touched[td]
	})
}

func (p *projectImp) prunePackages() {
	empty := map[Package]bool{}
	for _, p := range p.allPackages {
		if p.Empty() {
			empty[p] = true
		}
	}

	p.allPackages = slices.DeleteFunc(p.allPackages, func(pkg Package) bool {
		return empty[pkg]
	})

	/*
		for _, p := range p.allPackages {
			p.SetImports(slices.DeleteFunc(p.Imports(), func(pkg Package) bool {
				return empty[pkg]
			}))
		}
	*/
}

func (p *projectImp) String() string {
	return jsonify.ToString(p)
}

func (p *projectImp) Packages() []Package {
	return p.allPackages
}

func (p *projectImp) AppendPackage(pkg ...Package) {
	p.allPackages = append(p.allPackages, pkg...)
}

func (p *projectImp) FilterPackage(predicate func(pkg Package) bool) {
	p.allPackages = slices.DeleteFunc(p.allPackages, predicate)
}

func (p *projectImp) UpdateIndices() {
	// Type indices compound so that each has a unique offset.
	// The typeDefs in each package are also uniquely offset.
	index := 1
	index = p.allTypes.UpdateIndices(index)
	for i, pkg := range p.allPackages {
		index = pkg.SetIndices(i+1, index)
	}
}
