package usage

import (
	"github.com/Snow-Gremlin/goToolbox/collections"
	"github.com/Snow-Gremlin/goToolbox/collections/sortedSet"

	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs"
)

type factoryImp struct {
	usages collections.SortedSet[constructs.Usage]
}

func New() constructs.UsageFactory {
	return &factoryImp{usages: sortedSet.New(Comparer())}
}

func (f *factoryImp) NewUsage(args constructs.UsageArgs) constructs.Usage {
	v, _ := f.usages.TryAdd(newUsage(args))
	return v
}

func (f *factoryImp) Usages() collections.ReadonlySortedSet[constructs.Usage] {
	return f.usages.Readonly()
}
