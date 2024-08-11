package value

import (
	"github.com/Snow-Gremlin/goToolbox/collections"
	"github.com/Snow-Gremlin/goToolbox/collections/sortedSet"
)

type ValueFactory interface {
	NewValue(args Args) Value
	Values() collections.ReadonlySet[Value]
}

type factoryImp struct {
	values collections.SortedSet[Value]
}

func New() ValueFactory {
	return &factoryImp{values: sortedSet.New(Comparer())}
}

func (f *factoryImp) NewValue(args Args) Value {
	v, _ := f.values.TryAdd(newValue(args))
	return v
}

func (f *factoryImp) Values() collections.ReadonlySet[Value] {
	return f.values.Readonly()
}
