package structDesc

import (
	"github.com/Snow-Gremlin/goToolbox/collections"

	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs/kind"
)

type factoryImp struct {
	constructs.FactoryCore[constructs.StructDesc]
}

var _ constructs.Factory = (*factoryImp)(nil)

func New() constructs.StructDescFactory {
	return &factoryImp{FactoryCore: *constructs.NewFactoryCore(kind.StructDesc, Comparer())}
}

func (f *factoryImp) NewStructDesc(args constructs.StructDescArgs) constructs.StructDesc {
	return f.Add(newStructDesc(args))
}

func (f *factoryImp) StructDescs() collections.ReadonlySortedSet[constructs.StructDesc] {
	return f.Items().Readonly()
}
