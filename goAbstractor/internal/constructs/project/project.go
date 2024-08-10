package constructs

import (
	"go/token"

	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs/components/argument"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs/components/instance.go"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs/typeDesc/basic"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs/typeDesc/field"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs/typeDesc/reference"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs/typeDesc/typeParam"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/jsonify"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/locs"
)

type Project interface {
	argument.ArgumentFactory
	instance.InstanceFactory
	basic.BasicFactory
	field.FieldFactory
	reference.ReferenceFactory
	typeParam.TypeParamFactory
	jsonify.Jsonable

	NewLoc(pos token.Pos) locs.Loc
}

type projectImp struct {
	argument.ArgumentFactory
	instance.InstanceFactory
	basic.BasicFactory
	field.FieldFactory
	reference.ReferenceFactory
	typeParam.TypeParamFactory
	locations locs.Set
}

func NewProject(locs locs.Set) Project {
	return &projectImp{
		ArgumentFactory:  argument.New(),
		InstanceFactory:  instance.New(),
		BasicFactory:     basic.New(),
		FieldFactory:     field.New(),
		ReferenceFactory: reference.New(),
		TypeParamFactory: typeParam.New(),
		locations:        locs,
	}
}

func (p *projectImp) NewLoc(pos token.Pos) locs.Loc {
	return p.locations.NewLoc(pos)
}

func (p *projectImp) ToJson(ctx *jsonify.Context) jsonify.Datum {
	ctx2 := ctx.HideKind()
	return jsonify.NewMap().
		Add(ctx2, `language`, `go`).
		AddNonZero(ctx2, argument.Kind, p.Arguments()).
		AddNonZero(ctx2, instance.Kind, p.Instances()).
		AddNonZero(ctx2, basic.Kind, p.Basics()).
		AddNonZero(ctx2, field.Kind, p.Fields()).
		AddNonZero(ctx2, reference.Kind, p.References()).
		AddNonZero(ctx2, `locs`, p.locations)
}
