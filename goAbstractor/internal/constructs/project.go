package constructs

import (
	"encoding/json"

	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs/typeDesc"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/jsonify"
)

type Project struct {
	Packages []*Package

	AllStructs    []*typeDesc.Struct
	AllInterfaces []*typeDesc.Interface
	AllSignatures []*typeDesc.Signature
	AllTypeParams []*typeDesc.TypeParam
}

func (p *Project) ToJson(ctx *jsonify.Context) jsonify.Datum {
	m := jsonify.NewMap().
		Add(ctx, `language`, `go`).
		AddNonZero(ctx, `structs`, p.AllStructs).
		AddNonZero(ctx, `interfaces`, p.AllInterfaces).
		AddNonZero(ctx, `signatures`, p.AllSignatures).
		AddNonZero(ctx, `typeParams`, p.AllTypeParams)

	ctx = ctx.Copy().Set(`onlyIndex`, true)
	m.AddNonZero(ctx, `packages`, p.Packages)
	return m
}

func (p *Project) MarshalJSON() ([]byte, error) {
	ctx := jsonify.NewContext()
	return json.Marshal(p.ToJson(ctx))
}
