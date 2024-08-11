package constructs

import (
	"go/types"

	"github.com/Snow-Gremlin/goToolbox/collections"

	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/locs"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/metrics"
)

type Method interface {
	Declaration

	AddInstance(inst Instance) Instance
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
}

type MethodFactory interface {
	NewMethod(args MethodArgs) Method
	Methods() collections.ReadonlySet[Method]
}
