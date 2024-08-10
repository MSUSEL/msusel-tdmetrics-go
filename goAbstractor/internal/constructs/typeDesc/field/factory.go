package field

import (
	"github.com/Snow-Gremlin/goToolbox/collections"
	"github.com/Snow-Gremlin/goToolbox/collections/sortedSet"
)

type FieldFactory interface {
	NewField(args Args) Field
	Fields() collections.ReadonlySet[Field]
}

type factoryImp struct {
	fields collections.SortedSet[Field]
}

func New() FieldFactory {
	return &factoryImp{fields: sortedSet.New(Comparer())}
}

func (f *factoryImp) NewField(args Args) Field {
	v, _ := f.fields.TryAdd(newField(args))
	return v
}

func (f *factoryImp) Fields() collections.ReadonlySet[Field] {
	return f.fields.Readonly()
}
