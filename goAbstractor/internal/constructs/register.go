package constructs

import "github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/jsonify"

type Register interface {
	AllInterfaces() []Interface
	AllReferences() []TypeDefRef
	UpdateIndices(index int) int

	RegisterBasic(t Basic) Basic
	RegisterInterface(t Interface) Interface
	RegisterNamed(t Named) Named
	RegisterSignature(t Signature) Signature
	RegisterSolid(t Solid) Solid
	RegisterStruct(t Struct) Struct
	RegisterTypeDefRef(t TypeDefRef) TypeDefRef
	RegisterUnion(t Union) Union
}

func NewRegister() Register {
	return &registerImp{}
}

type registerImp struct {
	allBasics      []Basic
	allInterfaces  []Interface
	allNamed       []Named
	allSignatures  []Signature
	allSolids      []Solid
	allStructs     []Struct
	allTypeDefRefs []TypeDefRef
	allUnions      []Union
}

func (r *registerImp) ToJson(ctx *jsonify.Context) jsonify.Datum {
	ctx2 := ctx.HideKind()
	return jsonify.NewMap().
		AddNonZero(ctx2, `basics`, r.allBasics).
		AddNonZero(ctx2, `interfaces`, r.allInterfaces).
		AddNonZero(ctx2, `named`, r.allNamed).
		AddNonZero(ctx2, `signatures`, r.allSignatures).
		AddNonZero(ctx2, `solids`, r.allSolids).
		AddNonZero(ctx2, `structs`, r.allStructs).
		// Don't output r.allTypeDefRefs
		AddNonZero(ctx2, `unions`, r.allUnions)
}

func (r *registerImp) AllInterfaces() []Interface {
	return r.allInterfaces
}

func (r *registerImp) AllReferences() []TypeDefRef {
	return r.allTypeDefRefs
}

func (r *registerImp) UpdateIndices(index int) int {
	// Type indices compound so that each has a unique offset.
	// The typeDefs in each package are also uniquely offset.
	index = setIndices(index, r.allBasics)
	index = setIndices(index, r.allInterfaces)
	index = setIndices(index, r.allNamed)
	index = setIndices(index, r.allSignatures)
	index = setIndices(index, r.allSolids)
	index = setIndices(index, r.allStructs)
	// Don't index r.allTypeDefRefs
	index = setIndices(index, r.allUnions)
	return index
}

func setIndices[T TypeDesc](index int, s []T) int {
	for _, t := range s {
		t.SetIndex(index)
		index++
	}
	return index
}

func (r *registerImp) RegisterBasic(t Basic) Basic {
	return registerType(t, &r.allBasics)
}

func (r *registerImp) RegisterInterface(t Interface) Interface {
	return registerType(t, &r.allInterfaces)
}

func (r *registerImp) RegisterNamed(t Named) Named {
	return registerType(t, &r.allNamed)
}

func (r *registerImp) RegisterSignature(t Signature) Signature {
	return registerType(t, &r.allSignatures)
}

func (r *registerImp) RegisterSolid(t Solid) Solid {
	return registerType(t, &r.allSolids)
}

func (r *registerImp) RegisterStruct(t Struct) Struct {
	return registerType(t, &r.allStructs)
}

func (r *registerImp) RegisterTypeDefRef(t TypeDefRef) TypeDefRef {
	return registerType(t, &r.allTypeDefRefs)
}

func (r *registerImp) RegisterUnion(t Union) Union {
	return registerType(t, &r.allUnions)
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
