package signature

import (
	"github.com/Snow-Gremlin/goToolbox/collections"

	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs/kind"
)

type factoryImp struct {
	constructs.FactoryCore[constructs.Signature]
}

var _ constructs.Factory = (*factoryImp)(nil)

func New() constructs.SignatureFactory {
	return &factoryImp{FactoryCore: *constructs.NewFactoryCore(kind.Signature, Comparer())}
}

func (f *factoryImp) NewSignature(args constructs.SignatureArgs) constructs.Signature {
	return f.Add(newSignature(args))
}

func (f *factoryImp) Signatures() collections.ReadonlySortedSet[constructs.Signature] {
	return f.Items().Readonly()
}
