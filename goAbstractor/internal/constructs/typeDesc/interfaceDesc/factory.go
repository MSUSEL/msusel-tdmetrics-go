package interfaceDesc

import (
	"github.com/Snow-Gremlin/goToolbox/collections"
	"github.com/Snow-Gremlin/goToolbox/collections/sortedSet"
)

type InterfaceDescFactory interface {
	NewInterfaceDesc(args Args) InterfaceDesc
	InterfaceDescs() collections.ReadonlySet[InterfaceDesc]
}

type factoryImp struct {
	interfaceDescs collections.SortedSet[InterfaceDesc]
}

func New() InterfaceDescFactory {
	return &factoryImp{interfaceDescs: sortedSet.New(Comparer())}
}

func (f *factoryImp) NewInterfaceDesc(args Args) InterfaceDesc {
	v, _ := f.interfaceDescs.TryAdd(newInterfaceDesc(args))
	return v
}

func (f *factoryImp) InterfaceDescs() collections.ReadonlySet[InterfaceDesc] {
	return f.interfaceDescs.Readonly()
}
