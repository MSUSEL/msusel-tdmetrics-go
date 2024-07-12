package constructs

import (
	"fmt"

	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/jsonify"
	"github.com/Snow-Gremlin/goToolbox/utils"
)

type (
	Class interface {
		TypeDesc
		_class()

		Name() string
		Type() TypeDesc
		Methods() []Method
		TypeParams() []Named
		AppendMethod(met ...Method)
		SetInterface(inter Interface)
	}

	classImp struct {
		name       string
		typ        TypeDesc
		methods    []Method
		typeParams []Named
		inter      Interface
		index      int
	}
)

func newClass(name string, typ TypeDesc) Class {
	if len(name) <= 0 || utils.IsNil(typ) {
		panic(fmt.Errorf(`must have a name and type for a class definition: %q %v`,
			name, typ))
	}

	return &classImp{
		name: name,
		typ:  typ,
	}
}

func (td *classImp) _class() {}

func (td *classImp) SetIndex(index int) {
	td.index = index
}

func (td *classImp) ToJson(ctx *jsonify.Context) jsonify.Datum {
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

func (td *classImp) Visit(v Visitor) {
	visitTest(v, td.typ)
	visitList(v, td.methods)
	visitList(v, td.typeParams)
	visitTest(v, td.inter)
}

func (td *classImp) String() string {
	return jsonify.ToString(td)
}

func (td *classImp) Name() string {
	return td.name
}

func (td *classImp) Type() TypeDesc {
	return td.typ
}

func (td *classImp) Methods() []Method {
	return td.methods
}

func (td *classImp) TypeParams() []Named {
	return td.typeParams
}

func (td *classImp) AppendMethod(met ...Method) {
	td.methods = append(td.methods, met...)
}

func (td *classImp) SetInterface(inter Interface) {
	td.inter = inter
}
