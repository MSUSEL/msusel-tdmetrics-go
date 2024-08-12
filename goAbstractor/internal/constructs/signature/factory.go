package signature

import (
	"github.com/Snow-Gremlin/goToolbox/collections"
	"github.com/Snow-Gremlin/goToolbox/collections/sortedSet"

	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs"
)

type factoryImp struct {
	signatures collections.SortedSet[constructs.Signature]
}

func New() constructs.SignatureFactory {
	return &factoryImp{signatures: sortedSet.New(Comparer())}
}

func (f *factoryImp) NewSignature(args constructs.SignatureArgs) constructs.Signature {
	v, _ := f.signatures.TryAdd(newSignature(args))
	return v
}

func (f *factoryImp) Signatures() collections.ReadonlySortedSet[constructs.Signature] {
	return f.signatures.Readonly()
}
