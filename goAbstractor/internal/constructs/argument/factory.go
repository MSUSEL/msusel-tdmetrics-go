package argument

import (
	"github.com/Snow-Gremlin/goToolbox/collections"

	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs/kind"
)

type factoryImp struct {
	*constructs.FactoryCore[constructs.Argument]
}

var _ constructs.Factory = (*factoryImp)(nil)

func New() constructs.ArgumentFactory {
	return &factoryImp{FactoryCore: constructs.NewFactoryCore(kind.Argument, Comparer())}
}

func (f *factoryImp) NewArgument(args constructs.ArgumentArgs) constructs.Argument {
	return f.Add(newArgument(args))
}

func (f *factoryImp) Arguments() collections.ReadonlySortedSet[constructs.Argument] {
	return f.Items().Readonly()
}
