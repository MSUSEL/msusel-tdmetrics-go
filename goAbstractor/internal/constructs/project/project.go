package constructs

import (
	"go/token"

	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs/components/argument"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs/components/instance.go"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs/typeDesc/basic"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs/typeDesc/field"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs/typeDesc/interfaceDesc"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs/typeDesc/reference"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs/typeDesc/signature"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs/typeDesc/typeParam"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/jsonify"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/locs"
)

type Project interface {
	argument.ArgumentFactory
	instance.InstanceFactory

	basic.BasicFactory
	field.FieldFactory
	interfaceDesc.InterfaceDescFactory
	reference.ReferenceFactory
	signature.SignatureFactory
	typeParam.TypeParamFactory

	jsonify.Jsonable

	NewLoc(pos token.Pos) locs.Loc
}

type projectImp struct {
	argument.ArgumentFactory
	instance.InstanceFactory

	basic.BasicFactory
	field.FieldFactory
	interfaceDesc.InterfaceDescFactory
	reference.ReferenceFactory
	signature.SignatureFactory
	typeParam.TypeParamFactory

	locations locs.Set
}

func NewProject(locs locs.Set) Project {
	return &projectImp{
		ArgumentFactory: argument.New(),
		InstanceFactory: instance.New(),

		BasicFactory:         basic.New(),
		FieldFactory:         field.New(),
		InterfaceDescFactory: interfaceDesc.New(),
		ReferenceFactory:     reference.New(),
		SignatureFactory:     signature.New(),
		TypeParamFactory:     typeParam.New(),

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
		AddNonZero(ctx2, instance.Kind, p.Instances())

	m.AddNonZero(ctx2, basic.Kind, p.Basics()).
		AddNonZero(ctx2, field.Kind, p.Fields()).
		AddNonZero(ctx2, reference.Kind, p.References()).
		AddNonZero(ctx2, signature.Kind, p.Signatures()).
		AddNonZero(ctx2, typeParam.Kind, p.TypeParams())

	m.AddNonZero(ctx2, `locs`, p.locations)

	return m
}
