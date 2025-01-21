package tempTypeParamRef

import (
	"github.com/Snow-Gremlin/goToolbox/collections"
	"github.com/Snow-Gremlin/goToolbox/collections/sortedSet"

	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs"
)

type factoryImp struct {
	tempTypeParamRefs collections.SortedSet[constructs.TempTypeParamRef]
}

func New() constructs.TempTypeParamRefFactory {
	return &factoryImp{tempTypeParamRefs: sortedSet.New(Comparer())}
}

func (f *factoryImp) NewTempTypeParamRef(args constructs.TempTypeParamRefArgs) constructs.TempTypeParamRef {
	v, _ := f.tempTypeParamRefs.TryAdd(newTempTypeParamRef(args))
	return v
}

func (f *factoryImp) TempTypeParamRefs() collections.ReadonlySortedSet[constructs.TempTypeParamRef] {
	return f.tempTypeParamRefs.Readonly()
}

func (f *factoryImp) ClearAllTempTypeParamRefs() {
	f.tempTypeParamRefs.Clear()
}
