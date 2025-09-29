package abstract

import (
	"github.com/Snow-Gremlin/goToolbox/collections"

	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs/kind"
)

type factoryImp struct {
	*constructs.FactoryCore[constructs.Abstract]
}

var _ constructs.Factory = (*factoryImp)(nil)

func New() constructs.AbstractFactory {
	return &factoryImp{FactoryCore: constructs.NewFactoryCore(kind.Abstract, Comparer())}
}

func (f *factoryImp) NewAbstract(args constructs.AbstractArgs) constructs.Abstract {
	return f.Add(newAbstract(args))
}

func (f *factoryImp) Abstracts() collections.ReadonlySortedSet[constructs.Abstract] {
	return f.Items().Readonly()
}
