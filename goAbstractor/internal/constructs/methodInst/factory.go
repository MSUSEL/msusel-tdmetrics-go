package methodInst

import (
	"github.com/Snow-Gremlin/goToolbox/collections"

	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs/kind"
)

type factoryImp struct {
	*constructs.FactoryCore[constructs.MethodInst]
}

var _ constructs.Factory = (*factoryImp)(nil)

func New() constructs.MethodInstFactory {
	return &factoryImp{FactoryCore: constructs.NewFactoryCore(kind.MethodInst, Comparer())}
}

func (f *factoryImp) NewMethodInst(args constructs.MethodInstArgs) constructs.MethodInst {
	return f.Add(newInstance(args))
}

func (f *factoryImp) MethodInsts() collections.ReadonlySortedSet[constructs.MethodInst] {
	return f.Items().Readonly()
}
