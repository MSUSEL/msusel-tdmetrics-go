package interfaceDesc

import (
	"github.com/Snow-Gremlin/goToolbox/collections"
	"github.com/Snow-Gremlin/goToolbox/collections/sortedSet"

	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs"
)

type factoryImp struct {
	interfaceDescs collections.SortedSet[constructs.InterfaceDesc]
}

func New() constructs.InterfaceDescFactory {
	return &factoryImp{interfaceDescs: sortedSet.New(Comparer())}
}

func (f *factoryImp) NewInterfaceDesc(args constructs.InterfaceDescArgs) constructs.InterfaceDesc {
	v, _ := f.interfaceDescs.TryAdd(newInterfaceDesc(args))
	return v
}

func (f *factoryImp) InterfaceDescs() collections.ReadonlySortedSet[constructs.InterfaceDesc] {
	return f.interfaceDescs.Readonly()
}
