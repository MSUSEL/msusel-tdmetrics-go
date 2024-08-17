package project

import (
	"go/token"

	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs/abstract"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs/argument"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs/basic"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs/field"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs/instance"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs/interfaceDecl"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs/interfaceDesc"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs/method"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs/object"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs/packageCon"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs/reference"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs/signature"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs/structDesc"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs/typeParam"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs/value"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/locs"
)

type factoriesImp struct {
	constructs.AbstractFactory
	constructs.ArgumentFactory
	constructs.FieldFactory
	constructs.PackageFactory

	constructs.InterfaceDeclFactory
	constructs.MethodFactory
	constructs.ObjectFactory
	constructs.ValueFactory

	constructs.BasicFactory
	constructs.InstanceFactory
	constructs.InterfaceDescFactory
	constructs.ReferenceFactory
	constructs.SignatureFactory
	constructs.StructDescFactory
	constructs.TypeParamFactory

	locations locs.Set
}

func newFactory(locs locs.Set) constructs.Factories {
	return &factoriesImp{
		AbstractFactory: abstract.New(),
		ArgumentFactory: argument.New(),
		FieldFactory:    field.New(),
		PackageFactory:  packageCon.New(),

		InterfaceDeclFactory: interfaceDecl.New(),
		MethodFactory:        method.New(),
		ObjectFactory:        object.New(),
		ValueFactory:         value.New(),

		BasicFactory:         basic.New(),
		InstanceFactory:      instance.New(),
		InterfaceDescFactory: interfaceDesc.New(),
		ReferenceFactory:     reference.New(),
		SignatureFactory:     signature.New(),
		StructDescFactory:    structDesc.New(),
		TypeParamFactory:     typeParam.New(),

		locations: locs,
	}
}

func (p *factoriesImp) NewLoc(pos token.Pos) locs.Loc {
	return p.locations.NewLoc(pos)
}
