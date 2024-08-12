package object

import (
	"github.com/Snow-Gremlin/goToolbox/collections"
	"github.com/Snow-Gremlin/goToolbox/collections/sortedSet"

	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs"
)

type factoryImp struct {
	objects collections.SortedSet[constructs.Object]
}

func New() constructs.ObjectFactory {
	return &factoryImp{objects: sortedSet.New(Comparer())}
}

func (f *factoryImp) NewObject(args constructs.ObjectArgs) constructs.Object {
	v, _ := f.objects.TryAdd(newObject(args))
	return args.Package.AddObject(v)
}

func (f *factoryImp) Objects() collections.ReadonlySortedSet[constructs.Object] {
	return f.objects.Readonly()
}
