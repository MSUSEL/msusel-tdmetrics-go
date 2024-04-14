package typeDesc

import "github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/jsonify"

type Interface struct {
	Index    int
	Inherits []*Interface
	Methods  []*Func

	Inheritors []*Interface
}

func (ti *Interface) _isTypeDesc() {}

func (ti *Interface) ToJson(ctx *jsonify.Context) jsonify.Datum {
	if ctx.GetBool(`onlyIndex`) {
		return jsonify.New(ctx, ti.Index)
	}

	ctx2 := ctx.Copy().
		Remove(`noKind`).
		Set(`onlyIndex`, true)

	if len(ti.Methods) <= 0 {
		return jsonify.New(ctx2, `object`)
	}

	showKind := !ctx.GetBool(`noKind`)
	return jsonify.NewMap().
		AddIf(ctx2, showKind, `kind`, `interface`).
		AddNonZero(ctx2, `inherits`, ti.Inherits).
		AddNonZero(ctx2, `methods`, ti.Methods)
}

func (ti *Interface) HasMethod(m *Func) bool {
	for _, other := range ti.Methods {
		// The signatures have been registers so they can be compared by pointers.
		if m.Name == other.Name && m.Signature == other.Signature {
			return true
		}
	}
	return false
}

func (ti *Interface) IsSupertypeOf(other *Interface) bool {
	for _, m := range other.Methods {
		if !ti.HasMethod(m) {
			return false
		}
	}
	return true
}

type Func struct {
	Name      string
	Signature *Signature
}

func (f *Func) ToJson(ctx *jsonify.Context) jsonify.Datum {
	return jsonify.NewMap().
		Add(ctx, `name`, f.Name).
		Add(ctx, `signature`, f.Signature)
}
