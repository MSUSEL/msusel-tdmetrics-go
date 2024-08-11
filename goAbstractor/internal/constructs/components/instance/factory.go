package instance

import (
	"github.com/Snow-Gremlin/goToolbox/collections"
	"github.com/Snow-Gremlin/goToolbox/collections/sortedSet"
)

type InstanceFactory interface {
	NewInstance(args Args) Instance
	Instances() collections.ReadonlySet[Instance]
}

type factoryImp struct {
	instances collections.SortedSet[Instance]
}

func New() InstanceFactory {
	return &factoryImp{instances: sortedSet.New(Comparer())}
}

func (f *factoryImp) NewInstance(args Args) Instance {
	v, _ := f.instances.TryAdd(newInstance(args))
	return v
}

func (f *factoryImp) Instances() collections.ReadonlySet[Instance] {
	return f.instances.Readonly()
}
