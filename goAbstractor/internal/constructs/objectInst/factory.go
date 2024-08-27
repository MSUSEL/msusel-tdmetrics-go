package objectInst

import (
	"github.com/Snow-Gremlin/goToolbox/collections"
	"github.com/Snow-Gremlin/goToolbox/collections/sortedSet"

	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs"
)

type factoryImp struct {
	instances collections.SortedSet[constructs.ObjectInst]
}

func New() constructs.ObjectInstFactory {
	return &factoryImp{instances: sortedSet.New(Comparer())}
}

func (f *factoryImp) NewObjectInst(args constructs.ObjectInstArgs) constructs.ObjectInst {
	v, _ := f.instances.TryAdd(newInstance(args))
	return v
}

func (f *factoryImp) ObjectInsts() collections.ReadonlySortedSet[constructs.ObjectInst] {
	return f.instances.Readonly()
}
