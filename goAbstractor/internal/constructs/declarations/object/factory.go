package object

import (
	"github.com/Snow-Gremlin/goToolbox/collections"
	"github.com/Snow-Gremlin/goToolbox/collections/sortedSet"
)

type ObjectFactory interface {
	NewObject(args Args) Object
	Objects() collections.ReadonlySet[Object]
}

type factoryImp struct {
	objects collections.SortedSet[Object]
}

func New() ObjectFactory {
	return &factoryImp{objects: sortedSet.New(Comparer())}
}

func (f *factoryImp) NewObject(args Args) Object {
	v, _ := f.objects.TryAdd(newObject(args))
	return v
}

func (f *factoryImp) Objects() collections.ReadonlySet[Object] {
	return f.objects.Readonly()
}
