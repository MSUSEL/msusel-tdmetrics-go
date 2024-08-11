package value

import (
	"github.com/Snow-Gremlin/goToolbox/collections"
	"github.com/Snow-Gremlin/goToolbox/collections/sortedSet"

	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs/declarations"
)

type ValueFactory interface {
	NewValue(args Args) declarations.Value
	Values() collections.ReadonlySet[declarations.Value]
}

type factoryImp struct {
	values collections.SortedSet[declarations.Value]
}

func NewFactory() ValueFactory {
	return &factoryImp{values: sortedSet.New(Comparer())}
}

func (f *factoryImp) NewValue(args Args) declarations.Value {
	v, _ := f.values.TryAdd(New(args))
	return v
}

func (f *factoryImp) Values() collections.ReadonlySet[declarations.Value] {
	return f.values.Readonly()
}
