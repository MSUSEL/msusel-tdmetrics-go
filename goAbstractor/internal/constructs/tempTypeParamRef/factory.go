package tempTypeParamRef

import (
	"github.com/Snow-Gremlin/goToolbox/collections"

	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs/kind"
)

type factoryImp struct {
	constructs.FactoryCore[constructs.TempTypeParamRef]
}

func New() constructs.TempTypeParamRefFactory {
	return &factoryImp{FactoryCore: *constructs.NewFactoryCore(kind.TempTypeParamRef, Comparer())}
}

func (f *factoryImp) NewTempTypeParamRef(args constructs.TempTypeParamRefArgs) constructs.TempTypeParamRef {
	return f.Add(newTempTypeParamRef(args))
}

func (f *factoryImp) TempTypeParamRefs() collections.ReadonlySortedSet[constructs.TempTypeParamRef] {
	return f.Items().Readonly()
}

func (f *factoryImp) ClearAllTempTypeParamRefs() {
	f.Items().Clear()
}
