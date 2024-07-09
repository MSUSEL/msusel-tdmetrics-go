package constructs

import (
	"github.com/Snow-Gremlin/goToolbox/utils"

	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/jsonify"
)

type Types interface {
	AllInterfaces() []Interface
	AllReferences() []Reference
	UpdateIndices(index int) int
	Remove(predict func(TypeDesc) bool)

	RegisterBasic(t Basic) Basic
	RegisterInterface(t Interface) Interface
	RegisterNamed(t Named) Named
	RegisterSignature(t Signature) Signature
	RegisterSolid(t Solid) Solid
	RegisterStruct(t Struct) Struct
	RegisterTypeDefRef(t Reference) Reference
	RegisterUnion(t Union) Union
}

func NewTypes() Types {
	return &typesImp{}
}

type typesImp struct {
	allBasics     []Basic
	allInterfaces []Interface
	allNamed      []Named
	allSignatures []Signature
	allSolids     []Solid
	allStructs    []Struct
	allReferences []Reference
	allUnions     []Union
}

func (r *typesImp) ToJson(ctx *jsonify.Context) jsonify.Datum {
	ctx2 := ctx.HideKind()
	return jsonify.NewMap().
		AddNonZero(ctx2, `basics`, r.allBasics).
		AddNonZero(ctx2, `interfaces`, r.allInterfaces).
		AddNonZero(ctx2, `named`, r.allNamed).
		AddNonZero(ctx2, `signatures`, r.allSignatures).
		AddNonZero(ctx2, `solids`, r.allSolids).
		AddNonZero(ctx2, `structs`, r.allStructs).
		// Don't output r.allReferences
		AddNonZero(ctx2, `unions`, r.allUnions)
}

func (r *typesImp) AllInterfaces() []Interface {
	return r.allInterfaces
}

func (r *typesImp) AllReferences() []Reference {
	return r.allReferences
}

func (r *typesImp) UpdateIndices(index int) int {
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
	for _, td := range s {
		td.SetIndex(index)
		index++
	}
	return index
}

func (r *typesImp) Remove(predict func(TypeDesc) bool) {
	removeType(predict, &r.allBasics)
	removeType(predict, &r.allInterfaces)
	removeType(predict, &r.allNamed)
	removeType(predict, &r.allSignatures)
	removeType(predict, &r.allSolids)
	removeType(predict, &r.allStructs)
	removeType(predict, &r.allReferences)
	removeType(predict, &r.allUnions)
}

func removeType[T TypeDesc](predict func(TypeDesc) bool, s *[]T) {
	rs := *s
	zero := utils.Zero[T]()
	for i, td := range rs {
		if predict(td) {
			rs[i] = zero
		}
	}
	*s = utils.RemoveZeros(rs)
}

func (r *typesImp) RegisterBasic(t Basic) Basic {
	return registerType(t, &r.allBasics)
}

func (r *typesImp) RegisterInterface(t Interface) Interface {
	return registerType(t, &r.allInterfaces)
}

func (r *typesImp) RegisterNamed(t Named) Named {
	return registerType(t, &r.allNamed)
}

func (r *typesImp) RegisterSignature(t Signature) Signature {
	return registerType(t, &r.allSignatures)
}

func (r *typesImp) RegisterSolid(t Solid) Solid {
	return registerType(t, &r.allSolids)
}

func (r *typesImp) RegisterStruct(t Struct) Struct {
	return registerType(t, &r.allStructs)
}

func (r *typesImp) RegisterTypeDefRef(t Reference) Reference {
	return registerType(t, &r.allReferences)
}

func (r *typesImp) RegisterUnion(t Union) Union {
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
