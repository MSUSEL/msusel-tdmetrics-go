package constructs

import (
	"errors"
	"go/token"
	"go/types"

	"github.com/Snow-Gremlin/goToolbox/utils"
	"golang.org/x/tools/go/packages"

	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/jsonify"
)

type Signature interface {
	TypeDesc
	_signature()
}

type SignatureArgs struct {
	RealType   *types.Signature
	Variadic   bool
	Params     []Named
	TypeParams []Named
	Return     TypeDesc

	// Package is only needed if the real type is nil
	// so that a Go signature type has to be created.
	Package *packages.Package
}

func NewSignature(reg Register, args SignatureArgs) Signature {
	if utils.IsNil(args.RealType) {
		if utils.IsNil(args.Package) {
			panic(errors.New(`must provide a package if the real type for a signature is nil`))
		}

		pkg := args.Package.Types
		tp := make([]*types.TypeParam, len(args.TypeParams))
		for i, t := range args.TypeParams {
			tName := types.NewTypeName(token.NoPos, pkg, ``, t.GoType())
			tp[i] = types.NewTypeParam(tName, t.GoType())
		}

		params := make([]*types.Var, len(args.Params))
		for i, p := range args.Params {
			params[i] = types.NewVar(token.NoPos, pkg, p.Name(), p.GoType())
		}

		var ret *types.Tuple
		if !utils.IsNil(args.Return) {
			v := types.NewVar(token.NoPos, pkg, ``, args.Return.GoType())
			ret = types.NewTuple(v)
		}

		args.RealType = types.NewSignatureType(nil, nil,
			tp, types.NewTuple(params...), ret, args.Variadic)
	}

	return reg.RegisterSignature(&signatureImp{
		realType:   args.RealType,
		variadic:   args.Variadic,
		params:     args.Params,
		typeParams: args.TypeParams,
		returnType: args.Return,
	})
}

type signatureImp struct {
	realType *types.Signature

	variadic   bool
	params     []Named
	typeParams []Named
	returnType TypeDesc

	index int
}

func (sig *signatureImp) _signature() {}

func (sig *signatureImp) Visit(v Visitor) {
	visitList(v, sig.params)
	visitList(v, sig.typeParams)
	visitTest(v, sig.returnType)
}

func (sig *signatureImp) SetIndex(index int) {
	sig.index = index
}

func (sig *signatureImp) GoType() types.Type {
	return sig.realType
}

func (sig *signatureImp) Equal(other TypeDesc) bool {
	return equalTest(sig, other, func(a, b *signatureImp) bool {
		return a.variadic == b.variadic &&
			equal(a.returnType, b.returnType) &&
			equalList(a.params, b.params) &&
			equalList(a.typeParams, b.typeParams)
	})
}

func (sig *signatureImp) ToJson(ctx *jsonify.Context) jsonify.Datum {
	if ctx.IsShort() {
		return jsonify.New(ctx, sig.index)
	}

	ctx2 := ctx.HideKind().Short()
	return jsonify.NewMap().
		AddIf(ctx, ctx.IsKindShown(), `kind`, `signature`).
		AddNonZero(ctx2, `variadic`, sig.variadic).
		AddNonZero(ctx2, `params`, sig.params).
		AddNonZero(ctx2, `typeParams`, sig.typeParams).
		AddNonZero(ctx2, `return`, sig.returnType)
}

func (sig *signatureImp) String() string {
	return jsonify.ToString(sig)
}
