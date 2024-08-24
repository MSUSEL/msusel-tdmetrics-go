package constructs

import (
	"go/types"

	"github.com/Snow-Gremlin/goToolbox/collections"

	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/locs"
)

type Method interface {
	Declaration
	IsMethod()

	Signature() Signature
	Metrics() Metrics
	ReceiverName() string
	SetReceiver(recv Object)
	NeedsReceiver() bool
	Receiver() Object
	NoCopyRecv() bool

	IsInit() bool
	IsNamed() bool
	IsGeneric() bool
	HasReceiver() bool
}

type MethodArgs struct {
	RealType *types.Signature
	Package  Package
	Name     string
	Location locs.Loc

	TypeParams []TypeParam
	Signature  Signature
	Metrics    Metrics
	RecvName   string
	Receiver   Object

	// NoCopyRecv indicates the receiver is passed by pointer and therefore
	// should not be copied. This currently is not used in abstraction
	// because it doesn't matter if a copy is assigned to or not. If assigned
	// to at all, the receiver type is a mutable according to the abstraction.
	NoCopyRecv bool
}

type MethodFactory interface {
	NewMethod(args MethodArgs) Method
	Methods() collections.ReadonlySortedSet[Method]
}
