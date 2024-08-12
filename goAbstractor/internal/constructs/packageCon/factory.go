package packageCon

import (
	"github.com/Snow-Gremlin/goToolbox/collections"
	"github.com/Snow-Gremlin/goToolbox/collections/sortedSet"

	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs"
)

type factoryImp struct {
	packages collections.SortedSet[constructs.Package]
}

func New() constructs.PackageFactory {
	return &factoryImp{packages: sortedSet.New(Comparer())}
}

func (f *factoryImp) NewPackage(args constructs.PackageArgs) constructs.Package {
	v, _ := f.packages.TryAdd(newPackage(args))
	return v
}

func (f *factoryImp) Packages() collections.ReadonlySortedSet[constructs.Package] {
	return f.packages.Readonly()
}

func (p *factoryImp) FindPackageByPath(path string) constructs.Package {
	pkg, _ := p.packages.Enumerate().
		Where(func(pkg constructs.Package) bool { return pkg.Path() == path }).
		First()
	return pkg
}
