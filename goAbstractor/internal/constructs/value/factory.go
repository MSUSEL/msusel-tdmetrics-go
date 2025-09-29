package value

import (
	"github.com/Snow-Gremlin/goToolbox/collections"

	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs/kind"
)

type factoryImp struct {
	constructs.FactoryCore[constructs.Value]
}

var _ constructs.Factory = (*factoryImp)(nil)

func New() constructs.ValueFactory {
	return &factoryImp{FactoryCore: *constructs.NewFactoryCore(kind.Value, Comparer())}
}

func (f *factoryImp) NewValue(args constructs.ValueArgs) constructs.Value {
	v := f.Add(newValue(args))
	return args.Package.AddValue(v)
}

func (f *factoryImp) Values() collections.ReadonlySortedSet[constructs.Value] {
	return f.Items().Readonly()
}
