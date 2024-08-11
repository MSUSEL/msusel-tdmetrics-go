package reference

import (
	"github.com/Snow-Gremlin/goToolbox/collections"
	"github.com/Snow-Gremlin/goToolbox/collections/sortedSet"
)

type ReferenceFactory interface {
	NewReference(args Args) Reference
	References() collections.ReadonlySet[Reference]
}

type factoryImp struct {
	references collections.SortedSet[Reference]
}

func New() ReferenceFactory {
	return &factoryImp{references: sortedSet.New(Comparer())}
}

func (f *factoryImp) NewReference(args Args) Reference {
	v, _ := f.references.TryAdd(newReference(args))
	return v
}

func (f *factoryImp) References() collections.ReadonlySet[Reference] {
	return f.references.Readonly()
}
