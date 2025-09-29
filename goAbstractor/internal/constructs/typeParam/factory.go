package typeParam

import (
	"github.com/Snow-Gremlin/goToolbox/collections"

	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs/kind"
)

type factoryImp struct {
	constructs.FactoryCore[constructs.TypeParam]
}

var _ constructs.Factory = (*factoryImp)(nil)

func New() constructs.TypeParamFactory {
	return &factoryImp{FactoryCore: *constructs.NewFactoryCore(kind.TypeParam, Comparer())}
}

func (f *factoryImp) NewTypeParam(args constructs.TypeParamArgs) constructs.TypeParam {
	return f.Add(newTypeParam(args))
}

func (f *factoryImp) TypeParams() collections.ReadonlySortedSet[constructs.TypeParam] {
	return f.Items().Readonly()
}
