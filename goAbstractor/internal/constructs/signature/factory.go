package signature

import (
	"github.com/Snow-Gremlin/goToolbox/collections"
	"github.com/Snow-Gremlin/goToolbox/collections/sortedSet"

	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs/typeDescs"
)

type SignatureFactory interface {
	NewSignature(args Args) typeDescs.Signature
	Signatures() collections.ReadonlySet[typeDescs.Signature]
}

type factoryImp struct {
	signatures collections.SortedSet[typeDescs.Signature]
}

func NewFactory() SignatureFactory {
	return &factoryImp{signatures: sortedSet.New(Comparer())}
}

func (f *factoryImp) NewSignature(args Args) typeDescs.Signature {
	v, _ := f.signatures.TryAdd(New(args))
	return v
}

func (f *factoryImp) Signatures() collections.ReadonlySet[typeDescs.Signature] {
	return f.signatures.Readonly()
}
