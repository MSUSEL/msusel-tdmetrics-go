package packageCon

import (
	"github.com/Snow-Gremlin/goToolbox/collections"

	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs/kind"
)

type factoryImp struct {
	constructs.FactoryCore[constructs.Package]
}

var _ constructs.Factory = (*factoryImp)(nil)

func New() constructs.PackageFactory {
	return &factoryImp{FactoryCore: *constructs.NewFactoryCore(kind.Package, Comparer())}
}

func (f *factoryImp) NewPackage(args constructs.PackageArgs) constructs.Package {
	return f.Add(newPackage(args))
}

func (f *factoryImp) Packages() collections.ReadonlySortedSet[constructs.Package] {
	return f.Items().Readonly()
}

func (p *factoryImp) FindPackageByPath(path string) constructs.Package {
	pkg, _ := p.Items().Enumerate().
		Where(func(pkg constructs.Package) bool { return pkg.Path() == path }).
		First()
	return pkg
}
