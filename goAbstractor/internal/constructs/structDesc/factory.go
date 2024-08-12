package structDesc

import (
	"github.com/Snow-Gremlin/goToolbox/collections"
	"github.com/Snow-Gremlin/goToolbox/collections/sortedSet"

	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs"
)

type factoryImp struct {
	references collections.SortedSet[constructs.StructDesc]
}

func New() constructs.StructDescFactory {
	return &factoryImp{references: sortedSet.New(Comparer())}
}

func (f *factoryImp) NewStructDesc(args constructs.StructDescArgs) constructs.StructDesc {
	v, _ := f.references.TryAdd(newStructDesc(args))
	return v
}

func (f *factoryImp) StructDescs() collections.ReadonlySortedSet[constructs.StructDesc] {
	return f.references.Readonly()
}
