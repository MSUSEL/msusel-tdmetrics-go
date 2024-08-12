package reference

import (
	"github.com/Snow-Gremlin/goToolbox/collections"
	"github.com/Snow-Gremlin/goToolbox/collections/sortedSet"

	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs"
)

type factoryImp struct {
	references collections.SortedSet[constructs.Reference]
}

func New() constructs.ReferenceFactory {
	return &factoryImp{references: sortedSet.New(Comparer())}
}

func (f *factoryImp) NewReference(args constructs.ReferenceArgs) constructs.Reference {
	v, _ := f.references.TryAdd(newReference(args))
	return v
}

func (f *factoryImp) References() collections.ReadonlySortedSet[constructs.Reference] {
	return f.references.Readonly()
}
