package typeParam

import (
	"github.com/Snow-Gremlin/goToolbox/collections"
	"github.com/Snow-Gremlin/goToolbox/collections/sortedSet"

	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs/typeDescs"
)

type TypeParamFactory interface {
	NewTypeParam(args Args) typeDescs.TypeParam
	TypeParams() collections.ReadonlySet[typeDescs.TypeParam]
}

type factoryImp struct {
	references collections.SortedSet[typeDescs.TypeParam]
}

func NewFactory() TypeParamFactory {
	return &factoryImp{references: sortedSet.New(Comparer())}
}

func (f *factoryImp) NewTypeParam(args Args) typeDescs.TypeParam {
	v, _ := f.references.TryAdd(New(args))
	return v
}

func (f *factoryImp) TypeParams() collections.ReadonlySet[typeDescs.TypeParam] {
	return f.references.Readonly()
}
