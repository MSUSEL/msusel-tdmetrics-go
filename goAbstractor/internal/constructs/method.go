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
