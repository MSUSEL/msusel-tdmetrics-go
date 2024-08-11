package argument

import (
	"github.com/Snow-Gremlin/goToolbox/collections"
	"github.com/Snow-Gremlin/goToolbox/collections/sortedSet"

	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs"
)

type factoryImp struct {
	args collections.SortedSet[constructs.Argument]
}

func New() constructs.ArgumentFactory {
	return &factoryImp{args: sortedSet.New(Comparer())}
}

func (f *factoryImp) NewArgument(args constructs.ArgumentArgs) constructs.Argument {
	v, _ := f.args.TryAdd(newArgument(args))
	return v
}

func (f *factoryImp) Arguments() collections.ReadonlySet[constructs.Argument] {
	return f.args.Readonly()
}
