package object

import (
	"github.com/Snow-Gremlin/goToolbox/collections"
	"github.com/Snow-Gremlin/goToolbox/collections/sortedSet"

	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs/declarations"
)

type ObjectFactory interface {
	NewObject(args Args) declarations.Object
	Objects() collections.ReadonlySet[declarations.Object]
}

type factoryImp struct {
	objects collections.SortedSet[declarations.Object]
}

func NewFactory() ObjectFactory {
	return &factoryImp{objects: sortedSet.New(Comparer())}
}

func (f *factoryImp) NewObject(args Args) declarations.Object {
	v, _ := f.objects.TryAdd(New(args))
	return v
}

func (f *factoryImp) Objects() collections.ReadonlySet[declarations.Object] {
	return f.objects.Readonly()
}
