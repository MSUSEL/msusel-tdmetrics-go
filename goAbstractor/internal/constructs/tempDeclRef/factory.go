package tempDeclRef

import (
	"github.com/Snow-Gremlin/goToolbox/collections"

	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs/kind"
)

type factoryImp struct {
	constructs.FactoryCore[constructs.TempDeclRef]
}

var _ constructs.Factory = (*factoryImp)(nil)

func New() constructs.TempDeclRefFactory {
	return &factoryImp{FactoryCore: *constructs.NewFactoryCore(kind.TempDeclRef, Comparer())}
}

func (f *factoryImp) NewTempDeclRef(args constructs.TempDeclRefArgs) constructs.TempDeclRef {
	return f.Add(newTempDeclRef(args))
}

func (f *factoryImp) TempDeclRefs() collections.ReadonlySortedSet[constructs.TempDeclRef] {
	return f.Items().Readonly()
}

func (f *factoryImp) ClearAllTempDeclRefs() {
	f.Items().Clear()
}
