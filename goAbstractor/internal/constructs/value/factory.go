package value

import (
	"github.com/Snow-Gremlin/goToolbox/collections"
	"github.com/Snow-Gremlin/goToolbox/collections/sortedSet"

	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs"
)

type factoryImp struct {
	values collections.SortedSet[constructs.Value]
}

func New() constructs.ValueFactory {
	return &factoryImp{values: sortedSet.New(Comparer())}
}

func (f *factoryImp) NewValue(args constructs.ValueArgs) constructs.Value {
	v, _ := f.values.TryAdd(newValue(args))
	return args.Package.AddValue(v)
}

func (f *factoryImp) Values() collections.ReadonlySortedSet[constructs.Value] {
	return f.values.Readonly()
}
