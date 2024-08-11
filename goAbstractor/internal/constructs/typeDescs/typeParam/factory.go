package typeParam

import (
	"github.com/Snow-Gremlin/goToolbox/collections"
	"github.com/Snow-Gremlin/goToolbox/collections/sortedSet"
)

type TypeParamFactory interface {
	NewTypeParam(args Args) TypeParam
	TypeParams() collections.ReadonlySet[TypeParam]
}

type factoryImp struct {
	references collections.SortedSet[TypeParam]
}

func New() TypeParamFactory {
	return &factoryImp{references: sortedSet.New(Comparer())}
}

func (f *factoryImp) NewTypeParam(args Args) TypeParam {
	v, _ := f.references.TryAdd(newTypeParam(args))
	return v
}

func (f *factoryImp) TypeParams() collections.ReadonlySet[TypeParam] {
	return f.references.Readonly()
}
