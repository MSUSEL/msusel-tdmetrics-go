package selection

import (
	"github.com/Snow-Gremlin/goToolbox/collections"

	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs/kind"
)

type factoryImp struct {
	constructs.FactoryCore[constructs.Selection]
}

var _ constructs.Factory = (*factoryImp)(nil)

func New() constructs.SelectionFactory {
	return &factoryImp{FactoryCore: *constructs.NewFactoryCore(kind.Selection, Comparer())}
}

func (f *factoryImp) NewSelection(args constructs.SelectionArgs) constructs.Selection {
	return f.Add(newSelection(args))
}

func (f *factoryImp) Selections() collections.ReadonlySortedSet[constructs.Selection] {
	return f.Items().Readonly()
}
