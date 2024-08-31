package abstract

import (
	"github.com/Snow-Gremlin/goToolbox/comp"

	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/assert"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs/kind"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/jsonify"
)

type abstractImp struct {
	name      string
	signature constructs.Signature
	index     int
}

func newAbstract(args constructs.AbstractArgs) constructs.Abstract {
	assert.ArgValidId(`name`, args.Name)
	assert.ArgNotNil(`signature`, args.Signature)
	return &abstractImp{
		name:      args.Name,
		signature: args.Signature,
	}
}

func (a *abstractImp) IsAbstract() {}

func (a *abstractImp) Kind() kind.Kind    { return kind.Abstract }
func (a *abstractImp) Index() int         { return a.index }
func (a *abstractImp) SetIndex(index int) { a.index = index }

func (a *abstractImp) Name() string                    { return a.name }
func (a *abstractImp) Signature() constructs.Signature { return a.signature }

func (a *abstractImp) CompareTo(other constructs.Construct) int {
	return constructs.CompareTo[constructs.Abstract](a, other, Comparer())
}

func Comparer() comp.Comparer[constructs.Abstract] {
	return func(a, b constructs.Abstract) int {
		aImp, bImp := a.(*abstractImp), b.(*abstractImp)
		if aImp == bImp {
			return 0
		}
		return comp.Or(
			comp.DefaultPend(aImp.name, bImp.name),
			constructs.ComparerPend(aImp.signature, bImp.signature),
		)
	}
}

func (a *abstractImp) ToJson(ctx *jsonify.Context) jsonify.Datum {
	if ctx.IsOnlyIndex() {
		return jsonify.New(ctx, a.index)
	}
	if ctx.IsShort() {
		return jsonify.NewSprintf(`%s%d`, a.Kind(), a.index)
	}
	return jsonify.NewMap().
		AddIf(ctx, ctx.IsDebugKindIncluded(), `kind`, a.Kind()).
		AddIf(ctx, ctx.IsDebugIndexIncluded(), `index`, a.index).
		Add(ctx, `name`, a.name).
		Add(ctx.OnlyIndex(), `signature`, a.signature)
}

func (a *abstractImp) String() string {
	return a.name + ` ` + a.signature.String()
}
