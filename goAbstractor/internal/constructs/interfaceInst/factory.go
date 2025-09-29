package interfaceInst

import (
	"github.com/Snow-Gremlin/goToolbox/collections"

	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs/kind"
)

type factoryImp struct {
	*constructs.FactoryCore[constructs.InterfaceInst]
}

var _ constructs.Factory = (*factoryImp)(nil)

func New() constructs.InterfaceInstFactory {
	return &factoryImp{FactoryCore: constructs.NewFactoryCore(kind.InterfaceInst, Comparer())}
}

func (f *factoryImp) NewInterfaceInst(args constructs.InterfaceInstArgs) constructs.InterfaceInst {
	return f.Add(newInstance(args))
}

func (f *factoryImp) InterfaceInsts() collections.ReadonlySortedSet[constructs.InterfaceInst] {
	return f.Items().Readonly()
}
