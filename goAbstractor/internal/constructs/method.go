package constructs

import (
	"go/types"

	"github.com/Snow-Gremlin/goToolbox/collections"

	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/locs"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/metrics"
)

type Method interface {
	Declaration

	ReceiverName() string
	SetReceiver(recv Object)
	NeedsReceiver() bool

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
	Metrics    metrics.Metrics
	RecvName   string
	Receiver   Object

	// NoCopyRecv indicates the receiver is passed by pointer and therefore
	// should not be copied. This currently is not used in abstraction
	// because it doesn't matter if a copy is assigned to or not. If assigned
	// to at all, the receiver type is a mutable according to the abstraction.
	// (I just thought it was an interesting bit of information to collect.)
	NoCopyRecv bool
}

type MethodFactory interface {
	NewMethod(args MethodArgs) Method
	Methods() collections.ReadonlySortedSet[Method]
}
