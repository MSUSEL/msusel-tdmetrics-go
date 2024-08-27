package methodInst

import (
	"github.com/Snow-Gremlin/goToolbox/collections"
	"github.com/Snow-Gremlin/goToolbox/collections/sortedSet"

	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs"
)

type factoryImp struct {
	instances collections.SortedSet[constructs.MethodInst]
}

func New() constructs.MethodInstFactory {
	return &factoryImp{instances: sortedSet.New(Comparer())}
}

func (f *factoryImp) NewMethodInst(args constructs.MethodInstArgs) constructs.MethodInst {
	v, _ := f.instances.TryAdd(newInstance(args))
	return v
}

func (f *factoryImp) MethodInsts() collections.ReadonlySortedSet[constructs.MethodInst] {
	return f.instances.Readonly()
}
