package objectInst

import (
	"github.com/Snow-Gremlin/goToolbox/collections"

	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs/kind"
)

type factoryImp struct {
	constructs.FactoryCore[constructs.ObjectInst]
}

var _ constructs.Factory = (*factoryImp)(nil)

func New() constructs.ObjectInstFactory {
	return &factoryImp{FactoryCore: *constructs.NewFactoryCore(kind.ObjectInst, Comparer())}
}

func (f *factoryImp) NewObjectInst(args constructs.ObjectInstArgs) constructs.ObjectInst {
	return f.Add(newInstance(args))
}

func (f *factoryImp) ObjectInsts() collections.ReadonlySortedSet[constructs.ObjectInst] {
	return f.Items().Readonly()
}
