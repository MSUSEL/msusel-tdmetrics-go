package interfaceDesc

import (
	"github.com/Snow-Gremlin/goToolbox/collections"

	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs/kind"
)

type factoryImp struct {
	*constructs.FactoryCore[constructs.InterfaceDesc]
}

var _ constructs.Factory = (*factoryImp)(nil)

func New() constructs.InterfaceDescFactory {
	return &factoryImp{FactoryCore: constructs.NewFactoryCore(kind.InterfaceDesc, Comparer())}
}

func (f *factoryImp) NewInterfaceDesc(args constructs.InterfaceDescArgs) constructs.InterfaceDesc {
	return f.Add(newInterfaceDesc(args))
}

func (f *factoryImp) InterfaceDescs() collections.ReadonlySortedSet[constructs.InterfaceDesc] {
	return f.Items().Readonly()
}
