package typeDesc

import (
	"fmt"

	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/jsonify"
)

type Interface struct {
	Methods map[string]*Signature

	Index      int
	Inherits   []*Interface
	Inheritors []*Interface
}

func NewInterface() *Interface {
	return &Interface{
		Methods: map[string]*Signature{},
	}
}

func (ti *Interface) _isTypeDesc() {}

func (ti *Interface) ToJson(ctx *jsonify.Context) jsonify.Datum {
	if ctx.Short {
		return jsonify.New(ctx, ti.Index)
	}

	ctx2 := ctx.Copy()
	ctx2.NoKind = false
	ctx2.Short = true

	return jsonify.NewMap().
		AddIf(ctx2, !ctx.NoKind, `kind`, `interface`).
		AddNonZero(ctx2, `inherits`, ti.Inherits).
		AddNonZero(ctx2, `methods`, ti.Methods)
}

func (ti *Interface) String() string {
	return jsonify.ToString(ti)
}

func (ti *Interface) HasFunc(name string, sig *Signature) bool {
	other, has := ti.Methods[name]
	// The signature types have been registers
	// so they can be compared by pointers.
	return has && sig == other
}

func (ti *Interface) AddFunc(name string, sig *Signature) bool {
	if other, has := ti.Methods[name]; has {
		if other != sig {
			panic(fmt.Errorf(`function %v already exists with a different signature`, name))
		}
		return false
	}
	ti.Methods[name] = sig
	return true
}

func (ti *Interface) IsSupertypeOf(other *Interface) bool {
	for name, sig := range other.Methods {
		if !ti.HasFunc(name, sig) {
			return false
		}
	}
	return true
}
