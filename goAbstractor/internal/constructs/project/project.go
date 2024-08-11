package constructs

import (
	"go/token"

	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs/components/argument"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs/components/field"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs/components/instance"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs/declarations/interfaceDecl"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs/declarations/method"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs/declarations/object"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs/declarations/value"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs/typeDescs/basic"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs/typeDescs/interfaceDesc"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs/typeDescs/reference"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs/typeDescs/signature"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs/typeDescs/structDesc"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs/typeDescs/typeParam"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/jsonify"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/locs"
)

type Project interface {
	argument.ArgumentFactory
	field.FieldFactory
	instance.InstanceFactory

	interfaceDecl.InterfaceDeclFactory
	method.MethodFactory
	object.ObjectFactory
	value.ValueFactory

	basic.BasicFactory
	interfaceDesc.InterfaceDescFactory
	reference.ReferenceFactory
	signature.SignatureFactory
	structDesc.StructDescFactory
	typeParam.TypeParamFactory

	jsonify.Jsonable

	NewLoc(pos token.Pos) locs.Loc
}

type projectImp struct {
	argument.ArgumentFactory
	field.FieldFactory
	instance.InstanceFactory

	interfaceDecl.InterfaceDeclFactory
	method.MethodFactory
	object.ObjectFactory
	value.ValueFactory

	basic.BasicFactory
	interfaceDesc.InterfaceDescFactory
	reference.ReferenceFactory
	signature.SignatureFactory
	structDesc.StructDescFactory
	typeParam.TypeParamFactory

	locations locs.Set
}

func NewProject(locs locs.Set) Project {
	return &projectImp{
		ArgumentFactory: argument.New(),
		FieldFactory:    field.New(),
		InstanceFactory: instance.New(),

		InterfaceDeclFactory: interfaceDecl.New(),
		MethodFactory:        method.New(),
		ObjectFactory:        object.New(),
		ValueFactory:         value.New(),

		BasicFactory:         basic.NewFactory(),
		InterfaceDescFactory: interfaceDesc.NewFactory(),
		ReferenceFactory:     reference.NewFactory(),
		SignatureFactory:     signature.NewFactory(),
		StructDescFactory:    structDesc.NewFactory(),
		TypeParamFactory:     typeParam.NewFactory(),

		locations: locs,
	}
}

func (p *projectImp) NewLoc(pos token.Pos) locs.Loc {
	return p.locations.NewLoc(pos)
}

func (p *projectImp) ToJson(ctx *jsonify.Context) jsonify.Datum {
	ctx2 := ctx.HideKind()
	m := jsonify.NewMap().
		Add(ctx2, `language`, `go`)

	m.AddNonZero(ctx2, argument.Kind, p.Arguments()).
		AddNonZero(ctx2, field.Kind, p.Fields()).
		AddNonZero(ctx2, instance.Kind, p.Instances())

	m.AddNonZero(ctx2, basic.Kind, p.Basics()).
		AddNonZero(ctx2, reference.Kind, p.References()).
		AddNonZero(ctx2, signature.Kind, p.Signatures()).
		AddNonZero(ctx2, typeParam.Kind, p.TypeParams())

	m.AddNonZero(ctx2, `locs`, p.locations)

	return m
}
