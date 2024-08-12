package typeParam

import (
	"github.com/Snow-Gremlin/goToolbox/collections"
	"github.com/Snow-Gremlin/goToolbox/collections/sortedSet"

	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs"
)

type factoryImp struct {
	references collections.SortedSet[constructs.TypeParam]
}

func New() constructs.TypeParamFactory {
	return &factoryImp{references: sortedSet.New(Comparer())}
}

func (f *factoryImp) NewTypeParam(args constructs.TypeParamArgs) constructs.TypeParam {
	v, _ := f.references.TryAdd(newTypeParam(args))
	return v
}

func (f *factoryImp) TypeParams() collections.ReadonlySortedSet[constructs.TypeParam] {
	return f.references.Readonly()
}
