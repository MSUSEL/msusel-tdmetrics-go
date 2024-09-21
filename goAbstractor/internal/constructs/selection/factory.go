package selection

import (
	"github.com/Snow-Gremlin/goToolbox/collections"
	"github.com/Snow-Gremlin/goToolbox/collections/sortedSet"

	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs"
)

type factoryImp struct {
	selections collections.SortedSet[constructs.Selection]
}

func New() constructs.SelectionFactory {
	return &factoryImp{selections: sortedSet.New(Comparer())}
}

func (f *factoryImp) NewSelection(args constructs.SelectionArgs) constructs.Selection {
	v, _ := f.selections.TryAdd(newSelection(args))
	return v
}

func (f *factoryImp) Selections() collections.ReadonlySortedSet[constructs.Selection] {
	return f.selections.Readonly()
}
