package constructs

import (
	"go/types"

	"github.com/Snow-Gremlin/goToolbox/collections"
)

type Signature interface {
	TypeDesc
	IsSignature()

	// IsVacant indicates there are no parameters and no results,
	// i.e. `func()()`.
	IsVacant() bool
}

type SignatureArgs struct {
	RealType *types.Signature

	Variadic bool
	Params   []Argument
	Results  []Argument

	Package *types.Package
}

type SignatureFactory interface {
	NewSignature(args SignatureArgs) Signature
	Signatures() collections.ReadonlySortedSet[Signature]
}
