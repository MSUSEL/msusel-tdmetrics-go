package constructs

import (
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs/typeDesc"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/jsonify"
)

type TypeDef struct {
	Name     string
	Type     typeDesc.TypeDesc
	Inherits []*typeDesc.Interface
	Methods  []*Method
}

func (td *TypeDef) ToJson(ctx *jsonify.Context) jsonify.Datum {
	return jsonify.NewMap().
		Add(ctx, `name`, td.Name).
		Add(ctx, `type`, td.Type).
		AddNonZero(ctx, `inherits`, td.Inherits).
		AddNonZero(ctx, `methods`, td.Methods)
}

func (td *TypeDef) HasFunc(m *typeDesc.Func) bool {
	for _, other := range td.Methods {
		// The signatures have been registers so they can be compared by pointers.
		if m.Name == other.Name && m.Signature == other.Signature {
			return true
		}
	}
	return false
}

func (td *TypeDef) IsSupertypeOf(inter *typeDesc.Interface) bool {
	for _, m := range inter.Methods {
		if !td.HasFunc(m) {
			return false
		}
	}
	return true
}

func (td *TypeDef) String() string {
	return jsonify.ToString(td)
}
