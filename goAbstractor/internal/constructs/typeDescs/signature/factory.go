package signature

import (
	"github.com/Snow-Gremlin/goToolbox/collections"
	"github.com/Snow-Gremlin/goToolbox/collections/sortedSet"
)

type SignatureFactory interface {
	NewSignature(args Args) Signature
	Signatures() collections.ReadonlySet[Signature]
}

type factoryImp struct {
	signatures collections.SortedSet[Signature]
}

func New() SignatureFactory {
	return &factoryImp{signatures: sortedSet.New(Comparer())}
}

func (f *factoryImp) NewSignature(args Args) Signature {
	v, _ := f.signatures.TryAdd(newSignature(args))
	return v
}

func (f *factoryImp) Signatures() collections.ReadonlySet[Signature] {
	return f.signatures.Readonly()
}
