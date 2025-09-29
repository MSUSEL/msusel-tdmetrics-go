package method

import (
	"github.com/Snow-Gremlin/goToolbox/collections"

	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs/kind"
)

type factoryImp struct {
	*constructs.FactoryCore[constructs.Method]
}

var _ constructs.Factory = (*factoryImp)(nil)

func New() constructs.MethodFactory {
	return &factoryImp{FactoryCore: constructs.NewFactoryCore(kind.Method, Comparer())}
}

func (f *factoryImp) NewMethod(args constructs.MethodArgs) constructs.Method {
	v := f.Add(newMethod(args))
	return args.Package.AddMethod(v)
}

func (f *factoryImp) Methods() collections.ReadonlySortedSet[constructs.Method] {
	return f.Items().Readonly()
}
