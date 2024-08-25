package project

import (
	"strconv"

	"github.com/Snow-Gremlin/goToolbox/collections"
	"github.com/Snow-Gremlin/goToolbox/collections/enumerator"
	"github.com/Snow-Gremlin/goToolbox/terrors/terror"

	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/assert"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs/abstract"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs/argument"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs/basic"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs/field"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs/instance"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs/interfaceDecl"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs/interfaceDesc"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs/method"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs/metrics"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs/object"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs/packageCon"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs/signature"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs/structDesc"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs/tempReference"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs/typeParam"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs/value"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/jsonify"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/locs"
)

type projectImp struct {
	constructs.AbstractFactory
	constructs.ArgumentFactory
	constructs.FieldFactory
	constructs.PackageFactory
	constructs.MetricsFactory

	constructs.InterfaceDeclFactory
	constructs.MethodFactory
	constructs.ObjectFactory
	constructs.ValueFactory

	constructs.BasicFactory
	constructs.InstanceFactory
	constructs.InterfaceDescFactory
	constructs.SignatureFactory
	constructs.StructDescFactory
	constructs.TempReferenceFactory
	constructs.TypeParamFactory

	locations locs.Set
}

func New(locs locs.Set) constructs.Project {
	return &projectImp{
		AbstractFactory: abstract.New(),
		ArgumentFactory: argument.New(),
		FieldFactory:    field.New(),
		PackageFactory:  packageCon.New(),
		MetricsFactory:  metrics.New(),

		InterfaceDeclFactory: interfaceDecl.New(),
		MethodFactory:        method.New(),
		ObjectFactory:        object.New(),
		ValueFactory:         value.New(),

		BasicFactory:         basic.New(),
		InstanceFactory:      instance.New(),
		InterfaceDescFactory: interfaceDesc.New(),
		SignatureFactory:     signature.New(),
		StructDescFactory:    structDesc.New(),
		TempReferenceFactory: tempReference.New(),
		TypeParamFactory:     typeParam.New(),

		locations: locs,
	}
}

func (p *projectImp) Locs() locs.Set { return p.locations }

func (p *projectImp) AllConstructs() collections.Enumerator[constructs.Construct] {
	return enumerator.Enumerate[constructs.Construct]().Concat(
		enumerator.Cast[constructs.Construct](p.Abstracts().Enumerate()),
		enumerator.Cast[constructs.Construct](p.Arguments().Enumerate()),
		enumerator.Cast[constructs.Construct](p.Basics().Enumerate()),
		enumerator.Cast[constructs.Construct](p.Fields().Enumerate()),
		enumerator.Cast[constructs.Construct](p.Instances().Enumerate()),
		enumerator.Cast[constructs.Construct](p.InterfaceDecls().Enumerate()),
		enumerator.Cast[constructs.Construct](p.InterfaceDescs().Enumerate()),
		enumerator.Cast[constructs.Construct](p.Methods().Enumerate()),
		enumerator.Cast[constructs.Construct](p.Metrics().Enumerate()),
		enumerator.Cast[constructs.Construct](p.Objects().Enumerate()),
		enumerator.Cast[constructs.Construct](p.Packages().Enumerate()),
		enumerator.Cast[constructs.Construct](p.Signatures().Enumerate()),
		enumerator.Cast[constructs.Construct](p.StructDescs().Enumerate()),
		enumerator.Cast[constructs.Construct](p.TypeParams().Enumerate()),
		enumerator.Cast[constructs.Construct](p.TempReferences().Enumerate()),
		enumerator.Cast[constructs.Construct](p.Values().Enumerate()),
	)
}

func (p *projectImp) FindType(pkgPath, typeName string, panicOnNotFound bool) (constructs.Package, constructs.TypeDecl, bool) {
	assert.ArgNotEmpty(`pkgPath`, pkgPath)

	pkg := p.FindPackageByPath(pkgPath)
	if pkg == nil {
		if !panicOnNotFound {
			return nil, nil, false
		}
		names := enumerator.Select(p.Packages().Enumerate(),
			func(pkg constructs.Package) string { return strconv.Quote(pkg.Path()) }).
			Join(`, `)
		panic(terror.New(`failed to find package for type reference`).
			With(`type name`, typeName).
			With(`package path`, pkgPath).
			With(`existing paths`, `[`+names+`]`))
	}

	decl := pkg.FindTypeDecl(typeName)
	if decl == nil {
		if !panicOnNotFound {
			return pkg, nil, false
		}
		panic(terror.New(`failed to find type declaration for type reference`).
			With(`type name`, typeName).
			With(`package path`, pkgPath))
	}

	return pkg, decl, true
}

func (p *projectImp) String() string {
	return jsonify.ToString(p)
}

func (p *projectImp) ToJson(ctx *jsonify.Context) jsonify.Datum {
	ctx2 := ctx.HideKind()
	m := jsonify.NewMap().
		Add(ctx2, `language`, `go`).
		AddNonZero(ctx2, `locs`, p.locations)

	m.AddNonZero(ctx2, `abstracts`, p.Abstracts().ToSlice()).
		AddNonZero(ctx2, `arguments`, p.Arguments().ToSlice()).
		AddNonZero(ctx2, `fields`, p.Fields().ToSlice()).
		AddNonZero(ctx2, `packages`, p.Packages().ToSlice()).
		AddNonZero(ctx2, `metrics`, p.Metrics().ToSlice())

	m.AddNonZero(ctx2, `interfaceDecls`, p.InterfaceDecls().ToSlice()).
		AddNonZero(ctx2, `methods`, p.Methods().ToSlice()).
		AddNonZero(ctx2, `objects`, p.Objects().ToSlice()).
		AddNonZero(ctx2, `values`, p.Values().ToSlice())

	m.AddNonZero(ctx2, `basics`, p.Basics().ToSlice()).
		AddNonZero(ctx2, `instances`, p.Instances().ToSlice()).
		AddNonZero(ctx2, `interfaceDescs`, p.InterfaceDescs().ToSlice()).
		AddNonZero(ctx2, `tempReferences`, p.TempReferences().ToSlice()).
		AddNonZero(ctx2, `signatures`, p.Signatures().ToSlice()).
		AddNonZero(ctx2, `structDescs`, p.StructDescs().ToSlice()).
		AddNonZero(ctx2, `typeParams`, p.TypeParams().ToSlice())

	return m
}
