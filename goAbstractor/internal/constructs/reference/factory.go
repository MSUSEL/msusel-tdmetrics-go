package reference

import (
	"github.com/Snow-Gremlin/goToolbox/collections"
	"github.com/Snow-Gremlin/goToolbox/collections/sortedSet"

	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs/typeDescs"
)

type ReferenceFactory interface {
	NewReference(args Args) typeDescs.Reference
	References() collections.ReadonlySet[typeDescs.Reference]
}

type factoryImp struct {
	references collections.SortedSet[typeDescs.Reference]
}

func NewFactory() ReferenceFactory {
	return &factoryImp{references: sortedSet.New(Comparer())}
}

func (f *factoryImp) NewReference(args Args) typeDescs.Reference {
	v, _ := f.references.TryAdd(New(args))
	return v
}

func (f *factoryImp) References() collections.ReadonlySet[typeDescs.Reference] {
	return f.references.Readonly()
}
