package constructs

import (
	"fmt"

	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/jsonify"
	"github.com/Snow-Gremlin/goToolbox/utils"
)

type TypeDef interface {
	Visitable
	_typeDef()

	SetIndex(index int)
	Name() string
	Type() TypeDesc
	Methods() []Method
	TypeParams() []Named
	AppendMethod(met ...Method)
	SetInterface(inter Interface)
}

type typeDefImp struct {
	name       string
	typ        TypeDesc
	methods    []Method
	typeParams []Named
	inter      Interface
	index      int
}

func NewTypeDef(name string, typ TypeDesc) TypeDef {
	if len(name) <= 0 || utils.IsNil(typ) {
		panic(fmt.Errorf(`must have a name and type for a type definition: %q %v`,
			name, typ))
	}

	return &typeDefImp{
		name: name,
		typ:  typ,
	}
}

func (td *typeDefImp) _typeDef() {}

func (td *typeDefImp) SetIndex(index int) {
	td.index = index
}

func (td *typeDefImp) ToJson(ctx *jsonify.Context) jsonify.Datum {
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

func (td *typeDefImp) Visit(v Visitor) {
	visitTest(v, td.typ)
	visitList(v, td.methods)
	visitList(v, td.typeParams)
	visitTest(v, td.inter)
}

func (td *typeDefImp) String() string {
	return jsonify.ToString(td)
}

func (td *typeDefImp) Name() string {
	return td.name
}

func (td *typeDefImp) Type() TypeDesc {
	return td.typ
}

func (td *typeDefImp) Methods() []Method {
	return td.methods
}

func (td *typeDefImp) TypeParams() []Named {
	return td.typeParams
}

func (td *typeDefImp) AppendMethod(met ...Method) {
	td.methods = append(td.methods, met...)
}

func (td *typeDefImp) SetInterface(inter Interface) {
	td.inter = inter
}
