package constructs

import (
	"fmt"

	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/jsonify"
	"github.com/Snow-Gremlin/goToolbox/utils"
)

type (
	Value interface {
		Construct
		_value()

		Name() string
		Type() TypeDesc
		Methods() []Method
		TypeParams() []Named
		AppendMethod(met ...Method)
		SetInterface(inter Interface)
	}

	valueImp struct {
		name       string
		typ        TypeDesc
		methods    []Method
		typeParams []Named
		inter      Interface
		index      int
	}
)

func newValue(name string, typ TypeDesc) Value {
	if len(name) <= 0 || utils.IsNil(typ) {
		panic(fmt.Errorf(`must have a name and type for a class definition: %q %v`,
			name, typ))
	}

	return &valueImp{
		name: name,
		typ:  typ,
	}
}

func (td *valueImp) _value() {}

func (td *valueImp) SetIndex(index int) {
	td.index = index
}

func (td *valueImp) ToJson(ctx *jsonify.Context) jsonify.Datum {
	if ctx.IsShort() {
		return jsonify.New(ctx, td.index)
	}

	ctx2 := ctx.HideKind().Short()
	return jsonify.NewMap().
		Add(ctx2, `name`, td.name).
		Add(ctx2, `type`, td.typ).
		AddNonZero(ctx2, `methods`, td.methods).
		AddNonZero(ctx2, `typeParams`, td.typeParams).
		AddNonZero(ctx2, `interface`, td.inter)
}

func (td *valueImp) Visit(v Visitor) {
	visitTest(v, td.typ)
	visitList(v, td.methods)
	visitList(v, td.typeParams)
	visitTest(v, td.inter)
}

func (td *valueImp) String() string {
	return jsonify.ToString(td)
}

func (td *valueImp) Name() string {
	return td.name
}

func (td *valueImp) Type() TypeDesc {
	return td.typ
}

func (td *valueImp) Methods() []Method {
	return td.methods
}

func (td *valueImp) TypeParams() []Named {
	return td.typeParams
}

func (td *valueImp) AppendMethod(met ...Method) {
	td.methods = append(td.methods, met...)
}

func (td *valueImp) SetInterface(inter Interface) {
	td.inter = inter
}
