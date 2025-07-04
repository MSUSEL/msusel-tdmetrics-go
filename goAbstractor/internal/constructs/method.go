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
	NestType
	Declaration
	IsMethod()

	Signature() Signature
	Metrics() Metrics
	ReceiverName() string
	SetReceiver(recv Object)
	NeedsReceiver() bool
	Receiver() Object
	PointerRecv() bool

	// IsInit indicates this method is an init function,
	// i.e `func init() { ... }`.
	IsInit() bool

	// IsMain indicates this method is the main function,
	// i.e `func main() { ... }` in `main` package.
	IsMain() bool

	// IsTester indicates this method is a test, benchmark test, fuzzy test,
	// or example function.
	IsTester() bool

	// IsTest indicates this method is a test function,
	// i.e `func TestXxx(*testing.T) { ... }`.
	// See https://pkg.go.dev/testing
	IsTest() bool

	// IsBenchmark indicates this method is a benchmark test function,
	// i.e `func BenchmarkXxx(*testing.B) { ... }`.
	// See https://pkg.go.dev/testing
	IsBenchmark() bool

	// IsFuzz indicates this method is a fuzzy test function,
	// i.e `func FuzzXxx(*testing.F) { ... }`.
	// See https://pkg.go.dev/testing
	IsFuzz() bool

	// IsExample indicates this method is an example function,
	// i.e `func ExampleXxx() { ... }`.
	// See https://pkg.go.dev/testing
	IsExample() bool

	// IsConcreteFunc indicates this is a top-level function,
	// without a receiver, and the function is not generic.
	IsConcreteFunc() bool

	// IsNamed indicates this method is a method declared with a name,
	// otherwise it is a function literal without a name.
	IsNamed() bool

	// HasReceiver indicates this method has a receiver object and
	// is a member of that object.
	HasReceiver() bool

	// IsGeneric indicates this method is a generic method with zero or more
	// instances. If there are zero instances, then the generic function
	// hasn't been used and instantiated with a concrete type.
	IsGeneric() bool

	TypeParams() []TypeParam
	AddInstance(inst MethodInst) MethodInst
	Instances() collections.ReadonlySortedSet[MethodInst]
	FindInstance(instanceTypes []TypeDesc) (MethodInst, bool)
}

type MethodArgs struct {
	FuncType *types.Func
	SigType  *types.Signature
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
