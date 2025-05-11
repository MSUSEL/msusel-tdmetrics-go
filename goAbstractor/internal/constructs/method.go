package constructs

import (
	"go/types"

	"github.com/Snow-Gremlin/goToolbox/collections"

	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/locs"
)

// Method is a named method typically explicitly defined at the given location
// in the source code. The method may have a receiver for another project.
//
// If type parameters are given then the method is generic.
// Instances with realized versions of the method,
// are added for each used instance in the source code.
// If there are no instances then the generic method isn't used.
type Method interface {
	Declaration
	IsMethod()

	Signature() Signature
	Metrics() Metrics
	ReceiverName() string
	SetReceiver(recv Object)
	NeedsReceiver() bool
	Receiver() Object
	PointerRecv() bool

	IsInit() bool
	IsMain() bool
	IsNamed() bool
	HasReceiver() bool
	IsGeneric() bool
	TypeParams() []TypeParam
	AddInstance(inst MethodInst) MethodInst
	Instances() collections.ReadonlySortedSet[MethodInst]
	FindInstance(instanceTypes []TypeDesc) (MethodInst, bool)
}

type MethodArgs struct {
	RealType *types.Signature
	Package  Package
	Name     string
	Exported bool
	Location locs.Loc

	TypeParams []TypeParam
	Signature  Signature
	Metrics    Metrics
	RecvName   string
	Receiver   Object

	// PointerRecv indicates the receiver is passed by pointer
	// and therefore should not be copied. This also changes how
	// the method is accessible via an interface.
	PointerRecv bool
}

type MethodFactory interface {
	NewMethod(args MethodArgs) Method
	Methods() collections.ReadonlySortedSet[Method]
}
