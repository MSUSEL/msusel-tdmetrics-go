package constructs

import (
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/jsonify"
)

type TypeDef interface {
	Name() string
	Methods() []Method
	AppendMethod(met ...Method)
	SetInterface(inter Interface)
}

type typeDefImp struct {
	name       string
	typ        TypeDesc
	methods    []Method
	typeParams []Named
	inter      Interface
}

func NewTypeDef(name string, t TypeDesc) TypeDef {
	return &typeDefImp{
		name: name,
		typ:  t,
	}
}

func (td *typeDefImp) ToJson(ctx *jsonify.Context) jsonify.Datum {
	ctx2 := ctx.HideKind().Short()
	return jsonify.NewMap().
		Add(ctx2, `name`, td.name).
		Add(ctx2, `type`, td.typ).
		AddNonZero(ctx2, `methods`, td.methods).
		AddNonZero(ctx2, `typeParams`, td.typeParams).
		AddNonZero(ctx2, `interface`, td.inter)
}

func (td *typeDefImp) String() string {
	return jsonify.ToString(td)
}

func (td *typeDefImp) Name() string {
	return td.name
}

func (td *typeDefImp) Methods() []Method {
	return td.methods
}

func (td *typeDefImp) AppendMethod(met ...Method) {
	td.methods = append(td.methods, met...)
}

func (td *typeDefImp) SetInterface(inter Interface) {
	td.inter = inter
}
