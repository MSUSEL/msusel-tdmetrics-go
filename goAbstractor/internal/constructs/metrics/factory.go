package metrics

import (
	"github.com/Snow-Gremlin/goToolbox/collections"

	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs/kind"
)

type factoryImp struct {
	*constructs.FactoryCore[constructs.Metrics]
}

var _ constructs.Factory = (*factoryImp)(nil)

func New() constructs.MetricsFactory {
	return &factoryImp{FactoryCore: constructs.NewFactoryCore(kind.Metrics, Comparer())}
}

func (f *factoryImp) NewMetrics(args constructs.MetricsArgs) constructs.Metrics {
	return f.Add(newMetrics(args))
}

func (f *factoryImp) Metrics() collections.ReadonlySortedSet[constructs.Metrics] {
	return f.Items().Readonly()
}
