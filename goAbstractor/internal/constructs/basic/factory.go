package basic

import (
	"github.com/Snow-Gremlin/goToolbox/collections"

	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs/kind"
)

type factoryImp struct {
	*constructs.FactoryCore[constructs.Basic]
}

var _ constructs.Factory = (*factoryImp)(nil)

func New() constructs.BasicFactory {
	return &factoryImp{FactoryCore: constructs.NewFactoryCore(kind.Basic, Comparer())}
}

func (f *factoryImp) NewBasic(args constructs.BasicArgs) constructs.Basic {
	return f.Add(newBasic(args))
}

func (f *factoryImp) Basics() collections.ReadonlySortedSet[constructs.Basic] {
	return f.Items().Readonly()
}
