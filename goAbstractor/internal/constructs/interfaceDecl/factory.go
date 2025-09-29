package interfaceDecl

import (
	"github.com/Snow-Gremlin/goToolbox/collections"

	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs/kind"
)

type factoryImp struct {
	*constructs.FactoryCore[constructs.InterfaceDecl]
}

var _ constructs.Factory = (*factoryImp)(nil)

func New() constructs.InterfaceDeclFactory {
	return &factoryImp{FactoryCore: constructs.NewFactoryCore(kind.InterfaceDecl, Comparer())}
}

func (f *factoryImp) NewInterfaceDecl(args constructs.InterfaceDeclArgs) constructs.InterfaceDecl {
	v := f.Add(newInterfaceDecl(args))
	return args.Package.AddInterfaceDecl(v)
}

func (f *factoryImp) InterfaceDecls() collections.ReadonlySortedSet[constructs.InterfaceDecl] {
	return f.Items().Readonly()
}
