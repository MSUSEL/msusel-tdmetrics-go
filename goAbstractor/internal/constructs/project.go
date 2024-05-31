package constructs

import (
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs/typeDesc"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/jsonify"
)

type Project interface {
	Packages() []*Package
	AppendPackage(pkg ...*Package)
	AllInterfaces() []typeDesc.Interface

	// UpdateIndices should be called after all types have been registered
	// and all packages have been processed. This will update all the index
	// fields that will be used as references in the output models.
	UpdateIndices()
	RegisterInterface(t typeDesc.Interface) typeDesc.Interface
	RegisterSignature(t typeDesc.Signature) typeDesc.Signature
	RegisterStruct(t typeDesc.Struct) typeDesc.Struct
}

type projectImp struct {
	allPackages []*Package

	allInterfaces []typeDesc.Interface
	allSignatures []typeDesc.Signature
	allStructs    []typeDesc.Struct
}

func NewProject() Project {
	return &projectImp{}
}

func (p *projectImp) ToJson(ctx *jsonify.Context) jsonify.Datum {
	m := jsonify.NewMap().
		Add(ctx, `language`, `go`)

	ctx1 := ctx.HideKind()
	m.AddNonZero(ctx1, `interfaces`, p.allInterfaces).
		AddNonZero(ctx1, `signatures`, p.allSignatures).
		AddNonZero(ctx1, `structs`, p.allStructs).
		AddNonZero(ctx1, `packages`, p.allPackages)
	return m
}

func (p *projectImp) String() string {
	return jsonify.ToString(p)
}

func (p *projectImp) Packages() []*Package {
	return p.allPackages
}

func (p *projectImp) AppendPackage(pkg ...*Package) {
	p.allPackages = append(p.allPackages, pkg...)
}

func (p *projectImp) AllInterfaces() []typeDesc.Interface {
	return p.allInterfaces
}

func (p *projectImp) UpdateIndices() {
	// Type indices compound so that each has a unique offset.
	index := 0
	index = setIndices(index, p.allInterfaces)
	index = setIndices(index, p.allSignatures)
	setIndices(index, p.allStructs)

	// Package indices are independent.
	for i, pkg := range p.allPackages {
		pkg.Index = i
	}
}

func setIndices[T typeDesc.TypeDesc](index int, s []T) int {
	for _, t := range s {
		t.SetIndex(index)
		index++
	}
	return index
}

func (p *projectImp) RegisterInterface(t typeDesc.Interface) typeDesc.Interface {
	return registerType(t, &p.allInterfaces)
}

func (p *projectImp) RegisterSignature(t typeDesc.Signature) typeDesc.Signature {
	return registerType(t, &p.allSignatures)
}

func (p *projectImp) RegisterStruct(t typeDesc.Struct) typeDesc.Struct {
	return registerType(t, &p.allStructs)
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
