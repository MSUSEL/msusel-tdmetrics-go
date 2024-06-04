package constructs

import "github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/jsonify"

type Project interface {
	Register

	ToJson(ctx *jsonify.Context) jsonify.Datum
	Packages() []Package
	AppendPackage(pkg ...Package)
	AllInterfaces() []Interface
	AllReferences() []TypeDefRef

	// UpdateIndices should be called after all types have been registered
	// and all packages have been processed. This will update all the index
	// fields that will be used as references in the output models.
	UpdateIndices()
}

type projectImp struct {
	allPackages []Package

	allBasics     []Basic
	allInterfaces []Interface
	allNamed      []Named
	allSignatures []Signature
	allSolids     []Solid
	allStructs    []Struct
	allTypeDefRef []TypeDefRef
	allUnions     []Union
}

func NewProject() Project {
	return &projectImp{}
}

func (p *projectImp) ToJson(ctx *jsonify.Context) jsonify.Datum {
	m := jsonify.NewMap().
		Add(ctx, `language`, `go`)

	ctx1 := ctx.HideKind()
	m.AddNonZero(ctx1, `basics`, p.allBasics).
		AddNonZero(ctx1, `interfaces`, p.allInterfaces).
		AddNonZero(ctx1, `named`, p.allNamed).
		AddNonZero(ctx1, `signatures`, p.allSignatures).
		AddNonZero(ctx1, `solids`, p.allSolids).
		AddNonZero(ctx1, `structs`, p.allStructs).
		// Don't output typeDefRef.
		AddNonZero(ctx1, `unions`, p.allUnions).
		AddNonZero(ctx1, `packages`, p.allPackages)
	return m
}

func (p *projectImp) String() string {
	return jsonify.ToString(p)
}

func (p *projectImp) Packages() []Package {
	return p.allPackages
}

func (p *projectImp) AppendPackage(pkg ...Package) {
	p.allPackages = append(p.allPackages, pkg...)
}

func (p *projectImp) AllInterfaces() []Interface {
	return p.allInterfaces
}

func (p *projectImp) AllReferences() []TypeDefRef {
	return p.allTypeDefRef
}

func (p *projectImp) UpdateIndices() {
	// Type indices compound so that each has a unique offset.
	// Don't index typeDefRefs since they aren't outputted.
	index := 1
	index = setIndices(index, p.allBasics)
	index = setIndices(index, p.allInterfaces)
	index = setIndices(index, p.allNamed)
	index = setIndices(index, p.allSignatures)
	index = setIndices(index, p.allSolids)
	index = setIndices(index, p.allStructs)
	setIndices(index, p.allUnions)

	// Package indices are independent.
	for i, pkg := range p.allPackages {
		pkg.SetIndex(i)
	}
}

func setIndices[T TypeDesc](index int, s []T) int {
	for _, t := range s {
		t.SetIndex(index)
		index++
	}
	return index
}

func (p *projectImp) RegisterBasic(t Basic) Basic {
	return registerType(t, &p.allBasics)
}

func (p *projectImp) RegisterInterface(t Interface) Interface {
	return registerType(t, &p.allInterfaces)
}

func (p *projectImp) RegisterNamed(t Named) Named {
	return registerType(t, &p.allNamed)
}

func (p *projectImp) RegisterSignature(t Signature) Signature {
	return registerType(t, &p.allSignatures)
}

func (p *projectImp) RegisterSolid(t Solid) Solid {
	return registerType(t, &p.allSolids)
}

func (p *projectImp) RegisterStruct(t Struct) Struct {
	return registerType(t, &p.allStructs)
}

func (p *projectImp) RegisterTypeDefRef(t TypeDefRef) TypeDefRef {
	return registerType(t, &p.allTypeDefRef)
}

func (p *projectImp) RegisterUnion(t Union) Union {
	return registerType(t, &p.allUnions)
}

func registerType[T TypeDesc](t T, s *[]T) T {
	for _, t2 := range *s {
		if t.Equal(t2) {
			return t2
		}
	}
	*s = append(*s, t)
	return t
}
