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
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs/interfaceDecl"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs/interfaceDesc"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs/interfaceInst"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs/kind"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs/method"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs/methodInst"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs/metrics"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs/object"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs/objectInst"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs/packageCon"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs/selection"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs/signature"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs/structDesc"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs/tempDeclRef"
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
	constructs.SelectionFactory

	constructs.InterfaceDeclFactory
	constructs.MethodFactory
	constructs.ObjectFactory
	constructs.ValueFactory
	constructs.TempDeclRefFactory

	constructs.BasicFactory
	constructs.InterfaceDescFactory
	constructs.InterfaceInstFactory
	constructs.MethodInstFactory
	constructs.ObjectInstFactory
	constructs.SignatureFactory
	constructs.StructDescFactory
	constructs.TempReferenceFactory
	constructs.TypeParamFactory

	locations locs.Set
}

func New(locs locs.Set) constructs.Project {
	return &projectImp{
		AbstractFactory:  abstract.New(),
		ArgumentFactory:  argument.New(),
		FieldFactory:     field.New(),
		PackageFactory:   packageCon.New(),
		MetricsFactory:   metrics.New(),
		SelectionFactory: selection.New(),

		InterfaceDeclFactory: interfaceDecl.New(),
		MethodFactory:        method.New(),
		ObjectFactory:        object.New(),
		ValueFactory:         value.New(),
		TempDeclRefFactory:   tempDeclRef.New(),

		BasicFactory:         basic.New(),
		InterfaceDescFactory: interfaceDesc.New(),
		InterfaceInstFactory: interfaceInst.New(),
		MethodInstFactory:    methodInst.New(),
		ObjectInstFactory:    objectInst.New(),
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
		enumerator.Cast[constructs.Construct](p.InterfaceDecls().Enumerate()),
		enumerator.Cast[constructs.Construct](p.InterfaceDescs().Enumerate()),
		enumerator.Cast[constructs.Construct](p.InterfaceInsts().Enumerate()),
		enumerator.Cast[constructs.Construct](p.Methods().Enumerate()),
		enumerator.Cast[constructs.Construct](p.MethodInsts().Enumerate()),
		enumerator.Cast[constructs.Construct](p.Metrics().Enumerate()),
		enumerator.Cast[constructs.Construct](p.Objects().Enumerate()),
		enumerator.Cast[constructs.Construct](p.ObjectInsts().Enumerate()),
		enumerator.Cast[constructs.Construct](p.Packages().Enumerate()),
		enumerator.Cast[constructs.Construct](p.Selections().Enumerate()),
		enumerator.Cast[constructs.Construct](p.Signatures().Enumerate()),
		enumerator.Cast[constructs.Construct](p.StructDescs().Enumerate()),
		enumerator.Cast[constructs.Construct](p.TypeParams().Enumerate()),
		enumerator.Cast[constructs.Construct](p.TempDeclRefs().Enumerate()),
		enumerator.Cast[constructs.Construct](p.TempReferences().Enumerate()),
		enumerator.Cast[constructs.Construct](p.Values().Enumerate()),
	)
}

func (p *projectImp) EntryPoint() constructs.Package {
	pkg, _ := p.Packages().Enumerate().Where(func(pkg constructs.Package) bool {
		return pkg.EntryPoint()
	}).First()
	return pkg
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

func (p *projectImp) FindDecl(pkgPath, name string, panicOnNotFound bool) (constructs.Package, constructs.Declaration, bool) {
	assert.ArgNotEmpty(`pkgPath`, pkgPath)

	pkg := p.FindPackageByPath(pkgPath)
	if pkg == nil {
		if !panicOnNotFound {
			return nil, nil, false
		}
		names := enumerator.Select(p.Packages().Enumerate(),
			func(pkg constructs.Package) string { return strconv.Quote(pkg.Path()) }).
			Join(`, `)
		panic(terror.New(`failed to find package for declaration reference`).
			With(`name`, name).
			With(`package path`, pkgPath).
			With(`existing paths`, `[`+names+`]`))
	}

	method := pkg.FindDecl(name)
	if method == nil {
		if !panicOnNotFound {
			return pkg, nil, false
		}
		panic(terror.New(`failed to find declaration for declaration reference`).
			With(`decl name`, name).
			With(`package path`, pkgPath))
	}

	return pkg, method, true
}

func (p *projectImp) UpdateIndices() {
	var index int
	var kind kind.Kind
	p.AllConstructs().Foreach(func(c constructs.Construct) {
		if cKind := c.Kind(); kind != cKind {
			kind = cKind
			index = 0
		}
		index++
		c.SetIndex(index)
	})
}

func (p *projectImp) String() string {
	return jsonify.ToString(p)
}

func (p *projectImp) ToJson(ctx *jsonify.Context) jsonify.Datum {
	m := jsonify.NewMap().
		Add(ctx, `language`, `go`).
		AddNonZero(ctx, `locs`, p.locations)

	m.AddNonZero(ctx, `abstracts`, p.Abstracts().ToSlice()).
		AddNonZero(ctx, `arguments`, p.Arguments().ToSlice()).
		AddNonZero(ctx, `fields`, p.Fields().ToSlice()).
		AddNonZero(ctx, `packages`, p.Packages().ToSlice()).
		AddNonZero(ctx, `metrics`, p.Metrics().ToSlice()).
		AddNonZero(ctx, `selections`, p.Selections().ToSlice())

	m.AddNonZero(ctx, `interfaceDecls`, p.InterfaceDecls().ToSlice()).
		AddNonZero(ctx, `methods`, p.Methods().ToSlice()).
		AddNonZero(ctx, `objects`, p.Objects().ToSlice()).
		AddNonZero(ctx, `values`, p.Values().ToSlice()).
		AddNonZero(ctx, `tempDeclRef`, p.TempDeclRefs().ToSlice())

	m.AddNonZero(ctx, `basics`, p.Basics().ToSlice()).
		AddNonZero(ctx, `interfaceDescs`, p.InterfaceDescs().ToSlice()).
		AddNonZero(ctx, `interfaceInst`, p.InterfaceInsts().ToSlice()).
		AddNonZero(ctx, `methodInst`, p.MethodInsts().ToSlice()).
		AddNonZero(ctx, `objectInst`, p.ObjectInsts().ToSlice()).
		AddNonZero(ctx, `tempReferences`, p.TempReferences().ToSlice()).
		AddNonZero(ctx, `signatures`, p.Signatures().ToSlice()).
		AddNonZero(ctx, `structDescs`, p.StructDescs().ToSlice()).
		AddNonZero(ctx, `typeParams`, p.TypeParams().ToSlice())

	return m
}
