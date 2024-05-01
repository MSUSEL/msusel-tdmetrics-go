package constructs

import (
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs/typeDesc"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/jsonify"
)

type TypeDef struct {
	Name       string
	Type       typeDesc.TypeDesc
	Inherits   []*typeDesc.Interface
	Methods    []*Method
	TypeParams []*typeDesc.Named
}

func NewTypeDef(name string, t typeDesc.TypeDesc) *TypeDef {
	return &TypeDef{
		Name: name,
		Type: t,
	}
}

func (td *TypeDef) ToJson(ctx *jsonify.Context) jsonify.Datum {
	return jsonify.NewMap().
		Add(ctx, `name`, td.Name).
		Add(ctx, `type`, td.Type).
		AddNonZero(ctx, `inherits`, td.Inherits).
		AddNonZero(ctx, `methods`, td.Methods).
		AddNonZero(ctx, `typeParams`, td.TypeParams)
}

func (td *TypeDef) HasFunc(name string, sig *typeDesc.Signature) bool {
	for _, other := range td.Methods {
		// The signature types have been registers
		// so they can be compared by pointers.
		if name == other.Name && sig == other.Signature {
			return true
		}
	}
	return false
}

func (td *TypeDef) IsSupertypeOf(inter *typeDesc.Interface) bool {
	for name, m := range inter.Methods {
		if !td.HasFunc(name, m) {
			return false
		}
	}
	return true
}

func (td *TypeDef) String() string {
	return jsonify.ToString(td)
}
