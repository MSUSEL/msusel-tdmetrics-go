package constructs

import (
	"errors"
	"go/token"
	"go/types"

	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/jsonify"
	"github.com/Snow-Gremlin/goToolbox/utils"
	"golang.org/x/tools/go/packages"
)

type Struct interface {
	TypeDesc
	_struct()
}

type StructArgs struct {
	RealType *types.Struct
	Fields   []Named

	// Package is only needed if the real type is nil
	// so that a Go interface type has to be created.
	Package *packages.Package
}

func NewStruct(reg Register, args StructArgs) Struct {
	if utils.IsNil(args.RealType) {
		if utils.IsNil(args.Package) {
			panic(errors.New(`must provide a package if the real type for a struct is nil`))
		}

		fields := make([]*types.Var, len(args.Fields))
		pkg := args.Package.Types
		for i, f := range args.Fields {
			fields[i] = types.NewVar(token.NoPos, pkg, f.Name(), f.GoType())
		}

		args.RealType = types.NewStruct(fields, nil)
	}

	return reg.RegisterStruct(&structImp{
		realType: args.RealType,
		fields:   args.Fields,
	})
}

type structImp struct {
	realType *types.Struct
	fields   []Named
	index    int
}

func (ts *structImp) _struct() {}

func (ts *structImp) Visit(v Visitor) {
	visitList(v, ts.fields)
}

func (ts *structImp) SetIndex(index int) {
	ts.index = index
}

func (ts *structImp) GoType() types.Type {
	return ts.realType
}

func (ts *structImp) Equal(other TypeDesc) bool {
	return equalTest(ts, other, func(a, b *structImp) bool {
		return equalList(a.fields, b.fields)
	})
}

func (ts *structImp) ToJson(ctx *jsonify.Context) jsonify.Datum {
	if ctx.IsShort() {
		return jsonify.New(ctx, ts.index)
	}

	ctx2 := ctx.HideKind().Short()
	return jsonify.NewMap().
		AddIf(ctx, ctx.IsKindShown(), `kind`, `struct`).
		AddNonZero(ctx2, `fields`, ts.fields)
}

func (ts *structImp) String() string {
	return jsonify.ToString(ts)
}
