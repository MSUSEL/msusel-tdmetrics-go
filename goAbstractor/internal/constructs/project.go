package constructs

import "github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/jsonify"

type Project interface {
	Types() Register
	ToJson(ctx *jsonify.Context) jsonify.Datum
	Packages() []Package
	AppendPackage(pkg ...Package)

	// UpdateIndices should be called after all types have been registered
	// and all packages have been processed. This will update all the index
	// fields that will be used as references in the output models.
	UpdateIndices()
	Visit(v Visitor)
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

func (p *projectImp) Visit(v Visitor) {
	visitList(v, p.allPackages)
	// Do not visit the registered types.
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

func (p *projectImp) UpdateIndices() {
	// Type indices compound so that each has a unique offset.
	// The typeDefs in each package are also uniquely offset.
	index := 1
	index = p.allTypes.UpdateIndices(index)
	for i, pkg := range p.allPackages {
		index = pkg.SetIndices(i+1, index)
	}
}
