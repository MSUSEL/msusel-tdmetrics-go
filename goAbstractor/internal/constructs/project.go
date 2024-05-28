package constructs

import (
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs/typeDesc"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/jsonify"
)

type Project struct {
	Packages []*Package

	AllInterfaces []*typeDesc.Interface
	AllSignatures []*typeDesc.Signature
	AllStructs    []typeDesc.Struct
}

func NewProject() *Project {
	return &Project{}
}

func (p *Project) ToJson(ctx *jsonify.Context) jsonify.Datum {
	m := jsonify.NewMap().
		Add(ctx, `language`, `go`)

	ctx1 := ctx.HideKind()
	m.AddNonZero(ctx1, `interfaces`, p.AllInterfaces).
		AddNonZero(ctx1, `signatures`, p.AllSignatures).
		AddNonZero(ctx1, `structs`, p.AllStructs).
		AddNonZero(ctx1, `packages`, p.Packages)
	return m
}

func (p *Project) String() string {
	return jsonify.ToString(p)
}

// UpdateIndices should only be called after all types have been registered
// and all packages have been processed. This will update all the index
// fields that will be used as references in the output models.
func (p *Project) UpdateIndices() {
	// Type indices compound so that each has a unique offset.
	index := 0
	for _, t := range p.AllInterfaces {
		t.SetIndex(index)
		index++
	}
	for _, t := range p.AllSignatures {
		t.SetIndex(index)
		index++
	}
	for _, t := range p.AllStructs {
		t.SetIndex(index)
		index++
	}

	// Package indices are independent.
	for i, pkg := range p.Packages {
		pkg.Index = i
	}
}

func (p *Project) RegisterInterface(t *typeDesc.Interface) *typeDesc.Interface {
	return registerType(t, &p.AllInterfaces)
}

func (p *Project) RegisterSignature(t *typeDesc.Signature) *typeDesc.Signature {
	return registerType(t, &p.AllSignatures)
}

func (p *Project) RegisterStruct(t typeDesc.Struct) typeDesc.Struct {
	return registerType(t, &p.AllStructs)
}

func registerType[T typeDesc.TypeDesc](t T, s *[]T) T {
	for _, t2 := range *s {
		if t.Equal(t2) {
			return t2
		}
	}
	*s = append(*s, t)
	return t
}
