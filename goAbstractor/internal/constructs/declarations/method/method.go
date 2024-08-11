package method

import (
	"go/types"
	"strings"

	"github.com/Snow-Gremlin/goToolbox/terrors/terror"
	"github.com/Snow-Gremlin/goToolbox/utils"

	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/assert"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs/declarations"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/jsonify"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/locs"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/metrics"
)

type Method interface {
	declarations.Declaration
	_method()

	addInstance(inst Instance) Instance
	receiverName() string
	setReceiver(recv Object)
	needsReceiver() bool

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
	Variadic   bool
	Params     []Argument
	Results    []Argument

	Metrics metrics.Metrics

	NoCopyRecv bool
	RecvName   string
	Receiver   Object
}

type methodImp struct {
	realType *types.Signature
	pkg      Package
	name     string
	loc      locs.Loc

	typeParams []TypeParam
	variadic   bool
	params     []Argument
	results    []Argument

	metrics   metrics.Metrics
	instances Set[Instance]

	noCopyRecv bool
	recvName   string
	receiver   Object

	index int
}

func newMethod(args MethodArgs) Method {
	assert.ArgNotNil(`package`, args.Package)
	assert.ArgNoNils(`params`, args.Params)
	assert.ArgNoNils(`results`, args.Results)
	assert.ArgNotNil(`type params`, args.TypeParams)
	assert.ArgNotNil(`location`, args.Location)

	if !utils.IsNil(args.Receiver) {
		rName := args.Receiver.Name()
		if len(args.RecvName) > 0 && args.RecvName != rName {
			panic(terror.New(`name of receiver and a receiver were both given but the name didn't match the receiver`).
				With(`receiver name`, args.RecvName).
				With(`receiver`, rName))
		}
		args.RecvName = rName
	}

	if len(args.RecvName) > 0 && len(args.TypeParams) > 0 {
		panic(terror.New(`don't provide type params on a method with a receiver`))
	}

	if utils.IsNil(args.RealType) {
		if len(args.TypeParams) > 0 {
			panic(terror.New(`unsupported: cannot create a real type using type parameters`))
		}

		pkg := args.Package.Source().Types
		pos := args.Location.Pos()
		createTuple := func(args []Argument) *types.Tuple {
			vars := make([]*types.Var, len(args))
			for i, p := range args {
				vars[i] = types.NewVar(pos, pkg, p.Name(), p.Type().GoType())
			}
			return types.NewTuple(vars...)
		}
		par := createTuple(args.Params)
		ret := createTuple(args.Results)

		var recv *types.Var
		if !utils.IsNil(args.Receiver) {
			recv = types.NewVar(pos, pkg, args.Receiver.Name(), args.Receiver.GoType())
		} else if len(args.RecvName) > 0 {
			panic(terror.New(`cannot create a real type using only a receiver name`))
		}

		args.RealType = types.NewSignatureType(recv, nil, nil, par, ret, args.Variadic)
	}

	met := &methodImp{
		pkg:        args.Package,
		name:       args.Name,
		loc:        args.Location,
		typeParams: args.TypeParams,
		variadic:   args.Variadic,
		params:     args.Params,
		results:    args.Results,
		metrics:    args.Metrics,
		noCopyRecv: args.NoCopyRecv,
		recvName:   args.RecvName,
		receiver:   args.Receiver,
		instances:  NewSet[Instance](),
	}

	if !utils.IsNil(met.receiver) {
		return met.receiver.addMethod(met)
	}
	return met
}

func (m *methodImp) _method()           {}
func (m *methodImp) Kind() kind.Kind    { return kind.Method }
func (m *methodImp) setIndex(index int) { m.index = index }
func (m *methodImp) GoType() types.Type { return m.realType }

func (m *methodImp) Package() Package   { return m.pkg }
func (m *methodImp) Name() string       { return m.name }
func (m *methodImp) Location() locs.Loc { return m.loc }

func (m *methodImp) addInstance(inst Instance) Instance {
	return m.instances.Insert(inst)
}

func (m *methodImp) receiverName() string    { return m.recvName }
func (m *methodImp) setReceiver(recv Object) { m.receiver = recv }

func (m *methodImp) needsReceiver() bool {
	return utils.IsNil(m.receiver) && len(m.recvName) > 0
}

func (m *methodImp) IsInit() bool {
	return strings.HasPrefix(m.name, `init#`) &&
		len(m.recvName) <= 0 &&
		len(m.typeParams) <= 0 &&
		len(m.params) <= 0 &&
		len(m.results) <= 0
}

func (m *methodImp) IsNamed() bool {
	return len(m.name) > 0
}

func (m *methodImp) IsGeneric() bool {
	return len(m.typeParams) > 0
}

func (m *methodImp) HasReceiver() bool {
	return !utils.IsNil(m.receiver) || len(m.recvName) > 0
}

func (m *methodImp) compareTo(other Construct) int {
	b := other.(*methodImp)
	return or(
		func() int { return Compare(m.pkg, b.pkg) },
		func() int { return strings.Compare(m.name, b.name) },
		func() int { return compareSlice(m.typeParams, b.typeParams) },
		func() int { return compareSlice(m.params, b.params) },
		func() int { return compareSlice(m.results, b.results) },
		func() int { return boolCompare(m.variadic, b.variadic) },
		func() int { return strings.Compare(m.recvName, b.recvName) },
	)
}

func (m *methodImp) ToJson(ctx *jsonify.Context) jsonify.Datum {
	if ctx.IsShort() {
		return jsonify.New(ctx, m.index)
	}

	ctx2 := ctx.HideKind().Short()
	return jsonify.NewMap().
		AddIf(ctx, ctx.IsKindShown(), `kind`, m.Kind()).
		AddIf(ctx, ctx.IsIndexShown(), `index`, m.index).
		AddNonZero(ctx2, `package`, m.pkg).
		AddNonZero(ctx2, `name`, m.name).
		AddNonZero(ctx2, `loc`, m.loc).
		AddNonZero(ctx2, `typeParams`, m.typeParams).
		AddNonZero(ctx2, `variadic`, m.variadic).
		AddNonZero(ctx2, `params`, m.params).
		AddNonZero(ctx2, `results`, m.results).
		AddNonZero(ctx2, `metrics`, m.metrics).
		AddNonZero(ctx2, `instances`, m.instances).
		AddNonZero(ctx2, `receiver`, m.receiver).
		AddNonZeroIf(ctx2, ctx.IsReceiverShown(), `noCopyRecv`, m.noCopyRecv).
		AddNonZeroIf(ctx2, ctx.IsReceiverShown(), `recvName`, m.recvName)
}
