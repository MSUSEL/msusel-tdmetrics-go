package structDesc

import (
	"github.com/Snow-Gremlin/goToolbox/collections"
	"github.com/Snow-Gremlin/goToolbox/collections/sortedSet"

	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs/typeDescs"
)

type StructDescFactory interface {
	NewStructDesc(args Args) typeDescs.StructDesc
	StructDescs() collections.ReadonlySet[typeDescs.StructDesc]
}

type factoryImp struct {
	references collections.SortedSet[typeDescs.StructDesc]
}

func NewFactory() StructDescFactory {
	return &factoryImp{references: sortedSet.New(Comparer())}
}

func (f *factoryImp) NewStructDesc(args Args) typeDescs.StructDesc {
	v, _ := f.references.TryAdd(New(args))
	return v
}

func (f *factoryImp) StructDescs() collections.ReadonlySet[typeDescs.StructDesc] {
	return f.references.Readonly()
}
