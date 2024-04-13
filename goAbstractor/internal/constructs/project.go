package constructs

import (
	"encoding/json"

	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs/typeDesc"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/jsonify"
)

type Project struct {
	Packages []*Package

	AllInterfaces []*typeDesc.Interface
	AllSignatures []*typeDesc.Signature
	AllStructs    []*typeDesc.Struct
	AllTypeParams []*typeDesc.TypeParam
}

func (p *Project) ToJson(ctx *jsonify.Context) jsonify.Datum {
	m := jsonify.NewMap().
		Add(ctx, `language`, `go`)

	ctx1 := ctx.Copy().Set(`noKind`, true)
	m.AddNonZero(ctx1, `interfaces`, p.AllInterfaces).
		AddNonZero(ctx1, `signatures`, p.AllSignatures).
		AddNonZero(ctx1, `structs`, p.AllStructs).
		AddNonZero(ctx1, `typeParams`, p.AllTypeParams)

	ctx2 := ctx.Copy().Set(`onlyIndex`, true)
	m.AddNonZero(ctx2, `packages`, p.Packages)
	return m
}

func (p *Project) MarshalJSON() ([]byte, error) {
	ctx := jsonify.NewContext()
	return json.Marshal(p.ToJson(ctx))
}
