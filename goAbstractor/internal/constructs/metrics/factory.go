package metrics

import (
	"github.com/Snow-Gremlin/goToolbox/collections"
	"github.com/Snow-Gremlin/goToolbox/collections/sortedSet"

	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs"
)

type factoryImp struct {
	metrics collections.SortedSet[constructs.Metrics]
}

func New() constructs.MetricsFactory {
	return &factoryImp{metrics: sortedSet.New(Comparer())}
}

func (f *factoryImp) NewMetrics(args constructs.MetricsArgs) constructs.Metrics {
	v, _ := f.metrics.TryAdd(newMetrics(args))
	return v
}

func (f *factoryImp) Metrics() collections.ReadonlySortedSet[constructs.Metrics] {
	return f.metrics.Readonly()
}
