package instance

import (
	"github.com/Snow-Gremlin/goToolbox/collections"
	"github.com/Snow-Gremlin/goToolbox/collections/sortedSet"

	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs"
)

type factoryImp struct {
	instances collections.SortedSet[constructs.Instance]
}

func New() constructs.InstanceFactory {
	return &factoryImp{instances: sortedSet.New(Comparer())}
}

func (f *factoryImp) NewInstance(args constructs.InstanceArgs) constructs.Instance {
	v, _ := f.instances.TryAdd(newInstance(args))
	return v
}

func (f *factoryImp) Instances() collections.ReadonlySet[constructs.Instance] {
	return f.instances.Readonly()
}
