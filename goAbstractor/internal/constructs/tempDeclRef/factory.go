package tempDeclRef

import (
	"github.com/Snow-Gremlin/goToolbox/collections"
	"github.com/Snow-Gremlin/goToolbox/collections/sortedSet"

	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs"
)

type factoryImp struct {
	tempDeclRefs collections.SortedSet[constructs.TempDeclRef]
}

func New() constructs.TempDeclRefFactory {
	return &factoryImp{tempDeclRefs: sortedSet.New(Comparer())}
}

func (f *factoryImp) NewTempDeclRef(args constructs.TempDeclRefArgs) constructs.TempDeclRef {
	v, _ := f.tempDeclRefs.TryAdd(newTempDeclRef(args))
	return v
}

func (f *factoryImp) TempDeclRefs() collections.ReadonlySortedSet[constructs.TempDeclRef] {
	return f.tempDeclRefs.Readonly()
}

func (f *factoryImp) ClearAllTempDeclRefs() {
	f.tempDeclRefs.Clear()
}
