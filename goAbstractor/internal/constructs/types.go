package constructs

import (
	"go/types"

	"golang.org/x/tools/go/packages"

	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/jsonify"
)

type Types interface {
	NewBasic(typ *types.Basic) Basic
	NewBasicFromName(pkg *packages.Package, typeName string) Basic
	NewInterface(args InterfaceArgs) Interface
	NewNamed(name string, typ TypeDesc) Named
	NewReference(realType *types.Named, pkgPath, name string) Reference
	NewSignature(args SignatureArgs) Signature
	NewSolid(typ types.Type, target TypeDesc, tp ...TypeDesc) Solid
	NewStruct(args StructArgs) Struct
	NewUnion(args UnionArgs) Union

	AllInterfaces() []Interface
	AllReferences() []Reference
	UpdateIndices(index int) int
	Remove(predict func(TypeDesc) bool)
}

func newTypes() Types {
	return &typesImp{
		allBasics:     newTypeSet[Basic](),
		allInterfaces: newTypeSet[Interface](),
		allNamed:      newTypeSet[Named](),
		allReferences: newTypeSet[Reference](),
		allSignatures: newTypeSet[Signature](),
		allSolids:     newTypeSet[Solid](),
		allStructs:    newTypeSet[Struct](),
		allUnions:     newTypeSet[Union](),
	}
}

type typesImp struct {
	allBasics     *typeSet[Basic]
	allInterfaces *typeSet[Interface]
	allNamed      *typeSet[Named]
	allReferences *typeSet[Reference]
	allSignatures *typeSet[Signature]
	allSolids     *typeSet[Solid]
	allStructs    *typeSet[Struct]
	allUnions     *typeSet[Union]
}

func (r *typesImp) ToJson(ctx *jsonify.Context) jsonify.Datum {
	ctx2 := ctx.HideKind()
	return jsonify.NewMap().
		AddNonZero(ctx2, `basics`, r.allBasics.values).
		AddNonZero(ctx2, `interfaces`, r.allInterfaces.values).
		AddNonZero(ctx2, `named`, r.allNamed.values).
		// Don't output r.allReferences
		AddNonZero(ctx2, `signatures`, r.allSignatures.values).
		AddNonZero(ctx2, `solids`, r.allSolids.values).
		AddNonZero(ctx2, `structs`, r.allStructs.values).
		AddNonZero(ctx2, `unions`, r.allUnions.values)
}

func (r *typesImp) NewBasic(typ *types.Basic) Basic {
	return r.allBasics.Insert(newBasic(typ))
}

func (r *typesImp) NewBasicFromName(pkg *packages.Package, typeName string) Basic {
	return r.allBasics.Insert((newBasicFromName(pkg, typeName)))
}

func (r *typesImp) NewInterface(args InterfaceArgs) Interface {
	return r.allInterfaces.Insert(newInterface(args))
}

func (r *typesImp) NewNamed(name string, typ TypeDesc) Named {
	return r.allNamed.Insert(newNamed(name, typ))
}

func (r *typesImp) NewReference(realType *types.Named, pkgPath, name string) Reference {
	return r.allReferences.Insert(newReference(realType, pkgPath, name))
}

func (r *typesImp) NewSignature(args SignatureArgs) Signature {
	return r.allSignatures.Insert(newSignature(args))
}

func (r *typesImp) NewSolid(realType types.Type, target TypeDesc, tp ...TypeDesc) Solid {
	return r.allSolids.Insert(newSolid(realType, target, tp...))
}

func (r *typesImp) NewStruct(args StructArgs) Struct {
	return r.allStructs.Insert(newStruct(args))
}

func (r *typesImp) NewUnion(args UnionArgs) Union {
	return r.allUnions.Insert(newUnion(args))
}

func (r *typesImp) AllInterfaces() []Interface {
	return r.allInterfaces.values
}

func (r *typesImp) AllReferences() []Reference {
	return r.allReferences.values
}

func (r *typesImp) UpdateIndices(index int) int {
	// Type indices compound so that each has a unique offset.
	// The typeDefs in each package are also uniquely offset.
	index = r.allBasics.SetIndices(index)
	index = r.allInterfaces.SetIndices(index)
	index = r.allNamed.SetIndices(index)
	// Don't index r.allReferences
	index = r.allSignatures.SetIndices(index)
	index = r.allSolids.SetIndices(index)
	index = r.allStructs.SetIndices(index)
	index = r.allUnions.SetIndices(index)
	return index
}

func (r *typesImp) Remove(predict func(TypeDesc) bool) {
	r.allBasics.Remove(predict)
	r.allInterfaces.Remove(predict)
	r.allNamed.Remove(predict)
	r.allSignatures.Remove(predict)
	r.allSolids.Remove(predict)
	r.allStructs.Remove(predict)
	r.allReferences.Remove(predict)
	r.allUnions.Remove(predict)
}
