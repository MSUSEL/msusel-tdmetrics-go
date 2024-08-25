package tempReference

import (
	"github.com/Snow-Gremlin/goToolbox/collections"
	"github.com/Snow-Gremlin/goToolbox/collections/sortedSet"

	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs"
)

type factoryImp struct {
	tempReferences collections.SortedSet[constructs.TempReference]
}

func New() constructs.TempReferenceFactory {
	return &factoryImp{tempReferences: sortedSet.New(Comparer())}
}

func (f *factoryImp) NewTempReference(args constructs.TempReferenceArgs) constructs.TempReference {
	v, _ := f.tempReferences.TryAdd(newTempReference(args))
	return v
}

func (f *factoryImp) TempReferences() collections.ReadonlySortedSet[constructs.TempReference] {
	return f.tempReferences.Readonly()
}

func (f *factoryImp) ClearAllTempReferences() {
	f.tempReferences.Clear()
}
