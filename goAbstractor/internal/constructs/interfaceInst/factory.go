package interfaceInst

import (
	"github.com/Snow-Gremlin/goToolbox/collections"
	"github.com/Snow-Gremlin/goToolbox/collections/sortedSet"

	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs"
)

type factoryImp struct {
	instances collections.SortedSet[constructs.InterfaceInst]
}

func New() constructs.InterfaceInstFactory {
	return &factoryImp{instances: sortedSet.New(Comparer())}
}

func (f *factoryImp) NewInterfaceInst(args constructs.InterfaceInstArgs) constructs.InterfaceInst {
	v, _ := f.instances.TryAdd(newInstance(args))
	return v
}

func (f *factoryImp) InterfaceInsts() collections.ReadonlySortedSet[constructs.InterfaceInst] {
	return f.instances.Readonly()
}
