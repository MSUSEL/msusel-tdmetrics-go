package typeDesc

import (
	"sort"

	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/jsonify"
)

type Interface struct {
	Index    int
	Inherits []*Interface
	Methods  []*Func

	Inheritors []*Interface
}

func (ti *Interface) _isTypeDesc() {}

func (ti *Interface) ToJson(ctx *jsonify.Context) jsonify.Datum {
	if ctx.OnlyIndex {
		return jsonify.New(ctx, ti.Index)
	}

	ctx2 := ctx.Copy()
	ctx2.NoKind = false
	ctx2.OnlyIndex = true

	return jsonify.NewMap().
		AddIf(ctx2, !ctx.NoKind, `kind`, `interface`).
		AddNonZero(ctx2, `inherits`, ti.Inherits).
		AddNonZero(ctx2, `methods`, ti.Methods)
}

func (ti *Interface) HasFunc(m *Func) bool {
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
		if !ti.HasFunc(m) {
			return false
		}
	}
	return true
}

func (ti *Interface) String() string {
	return jsonify.ToString(ti)
}

func (ti *Interface) SortMethods() {
	sort.SliceIsSorted(ti.Methods, func(i, j int) bool {
		return ti.Methods[i].Name < ti.Methods[j].Name
	})
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

func (f *Func) String() string {
	return jsonify.ToString(f)
}
