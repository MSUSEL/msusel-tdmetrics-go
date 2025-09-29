package tempReference

import (
	"github.com/Snow-Gremlin/goToolbox/collections"

	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs/kind"
)

type factoryImp struct {
	constructs.FactoryCore[constructs.TempReference]
}

var _ constructs.Factory = (*factoryImp)(nil)

func New() constructs.TempReferenceFactory {
	return &factoryImp{FactoryCore: *constructs.NewFactoryCore(kind.TempReference, Comparer())}
}

func (f *factoryImp) NewTempReference(args constructs.TempReferenceArgs) constructs.TempReference {
	return f.Add(newTempReference(args))
}

func (f *factoryImp) TempReferences() collections.ReadonlySortedSet[constructs.TempReference] {
	return f.Items().Readonly()
}

func (f *factoryImp) ClearAllTempReferences() {
	f.Items().Clear()
}
