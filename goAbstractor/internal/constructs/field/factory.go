package field

import (
	"github.com/Snow-Gremlin/goToolbox/collections"
	"github.com/Snow-Gremlin/goToolbox/collections/sortedSet"

	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs"
)

type factoryImp struct {
	fields collections.SortedSet[constructs.Field]
}

func New() constructs.FieldFactory {
	return &factoryImp{fields: sortedSet.New(Comparer())}
}

func (f *factoryImp) NewField(args constructs.FieldArgs) constructs.Field {
	v, _ := f.fields.TryAdd(newField(args))
	return v
}

func (f *factoryImp) Fields() collections.ReadonlySortedSet[constructs.Field] {
	return f.fields.Readonly()
}
