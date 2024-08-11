package structDesc

import (
	"github.com/Snow-Gremlin/goToolbox/collections"
	"github.com/Snow-Gremlin/goToolbox/collections/sortedSet"
)

type StructDescFactory interface {
	NewStructDesc(args Args) StructDesc
	StructDescs() collections.ReadonlySet[StructDesc]
}

type factoryImp struct {
	references collections.SortedSet[StructDesc]
}

func New() StructDescFactory {
	return &factoryImp{references: sortedSet.New(Comparer())}
}

func (f *factoryImp) NewStructDesc(args Args) StructDesc {
	v, _ := f.references.TryAdd(newStructDesc(args))
	return v
}

func (f *factoryImp) StructDescs() collections.ReadonlySet[StructDesc] {
	return f.references.Readonly()
}
