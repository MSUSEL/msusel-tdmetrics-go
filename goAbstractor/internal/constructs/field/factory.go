package field

import (
	"github.com/Snow-Gremlin/goToolbox/collections"

	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs/kind"
)

type factoryImp struct {
	*constructs.FactoryCore[constructs.Field]
}

var _ constructs.Factory = (*factoryImp)(nil)

func New() constructs.FieldFactory {
	return &factoryImp{FactoryCore: constructs.NewFactoryCore(kind.Field, Comparer())}
}

func (f *factoryImp) NewField(args constructs.FieldArgs) constructs.Field {
	return f.Add(newField(args))
}

func (f *factoryImp) Fields() collections.ReadonlySortedSet[constructs.Field] {
	return f.Items().Readonly()
}
