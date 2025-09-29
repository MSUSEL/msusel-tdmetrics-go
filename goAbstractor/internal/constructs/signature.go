package constructs

import (
	"go/types"

	"golang.org/x/tools/go/packages"

	"github.com/Snow-Gremlin/goToolbox/collections"
)

type Signature interface {
	TypeDesc
	IsSignature()

	Variadic() bool
	Params() []Argument
	Results() []Argument

	// IsVacant indicates there are no parameters and no results,
	// i.e. `func()()`.
	IsVacant() bool
}

type SignatureArgs struct {
	RealType *types.Signature

	Variadic bool
	Params   []Argument
	Results  []Argument

	// Package is needed when the real type isn't given.
	// The package is used to help create the real type.
	Package *packages.Package
}

type SignatureFactory interface {
	Factory
	NewSignature(args SignatureArgs) Signature
	Signatures() collections.ReadonlySortedSet[Signature]
}
