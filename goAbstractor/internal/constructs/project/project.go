package project

import (
	"fmt"
	"strconv"
	"strings"

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
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs/tempTypeParamRef"
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
	constructs.TempTypeParamRefFactory
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

		BasicFactory:            basic.New(),
		InterfaceDescFactory:    interfaceDesc.New(),
		InterfaceInstFactory:    interfaceInst.New(),
		MethodInstFactory:       methodInst.New(),
		ObjectInstFactory:       objectInst.New(),
		SignatureFactory:        signature.New(),
		StructDescFactory:       structDesc.New(),
		TempReferenceFactory:    tempReference.New(),
		TempTypeParamRefFactory: tempTypeParamRef.New(),
		TypeParamFactory:        typeParam.New(),

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
		enumerator.Cast[constructs.Construct](p.TempTypeParamRefs().Enumerate()),
		enumerator.Cast[constructs.Construct](p.Values().Enumerate()),
	)
}

func (p *projectImp) EntryPoint() constructs.Package {
	pkg, _ := p.Packages().Enumerate().Where(func(pkg constructs.Package) bool {
		return pkg.EntryPoint()
	}).First()
	return pkg
}

func (p *projectImp) FindType(pkgPath, name string, instTypes []constructs.TypeDesc, allowRef, panicOnNotFound bool) (constructs.TypeDesc, bool) {
	assert.ArgNotEmpty(`pkgPath`, pkgPath)

	if allowRef {
		itComp := constructs.SliceComparer[constructs.TypeDesc]()
		ref, found := p.TempReferences().Enumerate().Where(func(ref constructs.TempReference) bool {
			return ref.PackagePath() == pkgPath && ref.Name() == name && itComp(ref.InstanceTypes(), instTypes) == 0
		}).First()
		if found {
			return ref, true
		}
	}

	pkg := p.FindPackageByPath(pkgPath)
	if pkg == nil {
		if !panicOnNotFound {
			return nil, false
		}
		names := enumerator.Select(p.Packages().Enumerate(),
			func(pkg constructs.Package) string { return strconv.Quote(pkg.Path()) }).
			Join(`, `)
		panic(terror.New(`failed to find package for type reference`).
			With(`type name`, name).
			With(`instance types`, instTypes).
			With(`package path`, pkgPath).
			With(`existing paths`, `[`+names+`]`))
	}

	decl := pkg.FindTypeDecl(name)
	if decl == nil {
		if !panicOnNotFound {
			return nil, false
		}
		panic(terror.New(`failed to find type for type reference`).
			With(`type name`, name).
			With(`instance types`, instTypes).
			With(`package path`, pkgPath))
	}

	if len(instTypes) > 0 {
		switch t := decl.(type) {
		case constructs.InterfaceDecl:
			if inst, ok := t.FindInstance(instTypes); ok {
				return inst, true
			}
			return nil, false
		case constructs.Object:
			if inst, ok := t.FindInstance(instTypes); ok {
				return inst, true
			}
			return nil, false
		case constructs.Method:
			panic(terror.New(`can not use method instance as type for type reference`).
				With(`type name`, name).
				With(`instance types`, instTypes).
				With(`method`, t).
				With(`package path`, pkgPath))
		case constructs.Value:
			panic(terror.New(`can not get an instance of a value for type reference`).
				With(`type name`, name).
				With(`instance types`, instTypes).
				With(`value`, t).
				With(`package path`, pkgPath))
		}
		panic(terror.New(`unexpected type for type reference instance`).
			With(`type name`, name).
			With(`instance types`, instTypes).
			With(`declaration`, decl).
			With(`package path`, pkgPath))
	}

	return decl, true
}

func (p *projectImp) FindDecl(pkgPath, name string, instTypes []constructs.TypeDesc, allowRef, panicOnNotFound bool) (constructs.Construct, bool) {
	assert.ArgNotEmpty(`pkgPath`, pkgPath)

	if allowRef {
		itComp := constructs.SliceComparer[constructs.TypeDesc]()
		ref, found := p.TempDeclRefs().Enumerate().Where(func(ref constructs.TempDeclRef) bool {
			return ref.PackagePath() == pkgPath && ref.Name() == name && itComp(ref.InstanceTypes(), instTypes) == 0
		}).First()
		if found {
			return ref, true
		}
	}

	pkg := p.FindPackageByPath(pkgPath)
	if pkg == nil {
		if !panicOnNotFound {
			return nil, false
		}
		names := enumerator.Select(p.Packages().Enumerate(),
			func(pkg constructs.Package) string { return strconv.Quote(pkg.Path()) }).
			Join(`, `)
		panic(terror.New(`failed to find package for declaration reference`).
			With(`name`, name).
			With(`package path`, pkgPath).
			With(`existing paths`, `[`+names+`]`))
	}

	decl := pkg.FindDecl(name)
	if decl == nil {
		if !panicOnNotFound {
			return nil, false
		}
		panic(terror.New(`failed to find declaration for declaration reference`).
			With(`decl name`, name).
			With(`package path`, pkgPath))
	}

	if len(instTypes) > 0 {
		switch t := decl.(type) {
		case constructs.InterfaceDecl:
			if inst, ok := t.FindInstance(instTypes); ok {
				return inst, true
			}
			return nil, false
		case constructs.Object:
			if inst, ok := t.FindInstance(instTypes); ok {
				return inst, true
			}
			return nil, false
		case constructs.Method:
			if inst, ok := t.FindInstance(instTypes); ok {
				return inst, true
			}
			return nil, false
		case constructs.Value:
			panic(terror.New(`can not get an instance of a value for declaration reference`).
				With(`type name`, name).
				With(`instance types`, instTypes).
				With(`value`, t).
				With(`package path`, pkgPath))
		}
		panic(terror.New(`unexpected declaration for declaration reference instance`).
			With(`type name`, name).
			With(`instance types`, instTypes).
			With(`declaration`, decl).
			With(`package path`, pkgPath))
	}

	return decl, true
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

func pl(k kind.Kind) string {
	s := string(k)
	if !strings.HasSuffix(s, `s`) {
		s += `s`
	}
	return s
}

func (p *projectImp) ToJson(ctx *jsonify.Context) jsonify.Datum {
	m := jsonify.NewMap().
		Add(ctx, `language`, `go`).
		AddNonZero(ctx, `locs`, p.locations)

	m.AddNonZero(ctx, pl(kind.Abstract), p.Abstracts().ToSlice()).
		AddNonZero(ctx, pl(kind.Argument), p.Arguments().ToSlice()).
		AddNonZero(ctx, pl(kind.Field), p.Fields().ToSlice()).
		AddNonZero(ctx, pl(kind.Package), p.Packages().ToSlice()).
		AddNonZero(ctx, pl(kind.Metrics), p.Metrics().ToSlice()).
		AddNonZero(ctx, pl(kind.Selection), p.Selections().ToSlice()).
		AddNonZero(ctx, pl(kind.TempDeclRef), p.TempDeclRefs().ToSlice())

	m.AddNonZero(ctx, pl(kind.InterfaceDecl), p.InterfaceDecls().ToSlice()).
		AddNonZero(ctx, pl(kind.Method), p.Methods().ToSlice()).
		AddNonZero(ctx, pl(kind.Object), p.Objects().ToSlice()).
		AddNonZero(ctx, pl(kind.Value), p.Values().ToSlice())

	m.AddNonZero(ctx, pl(kind.Basic), p.Basics().ToSlice()).
		AddNonZero(ctx, pl(kind.InterfaceDesc), p.InterfaceDescs().ToSlice()).
		AddNonZero(ctx, pl(kind.InterfaceInst), p.InterfaceInsts().ToSlice()).
		AddNonZero(ctx, pl(kind.MethodInst), p.MethodInsts().ToSlice()).
		AddNonZero(ctx, pl(kind.ObjectInst), p.ObjectInsts().ToSlice()).
		AddNonZero(ctx, pl(kind.TempReference), p.TempReferences().ToSlice()).
		AddNonZero(ctx, pl(kind.TempTypeParamRef), p.TempTypeParamRefs().ToSlice()).
		AddNonZero(ctx, pl(kind.Signature), p.Signatures().ToSlice()).
		AddNonZero(ctx, pl(kind.StructDesc), p.StructDescs().ToSlice()).
		AddNonZero(ctx, pl(kind.TypeParam), p.TypeParams().ToSlice())

	return m
}

func stringCon[T any, S ~[]T](buf *strings.Builder, k kind.Kind, s S) {
	buf.WriteString(pl(k))
	if len(s) <= 0 {
		buf.WriteString(" { }\n")
		return
	}
	buf.WriteString(" {\n")
	for i, k := range s {
		buf.WriteString(fmt.Sprintf("  %2d. %q\n", i+1, fmt.Sprint(k)))
	}
	buf.WriteString("}\n")
}

func (p *projectImp) String() string {
	buf := &strings.Builder{}
	stringCon(buf, kind.Abstract, p.Abstracts().ToSlice())
	stringCon(buf, kind.Argument, p.Arguments().ToSlice())
	stringCon(buf, kind.Basic, p.Basics().ToSlice())
	stringCon(buf, kind.Field, p.Fields().ToSlice())
	stringCon(buf, kind.InterfaceDecl, p.InterfaceDecls().ToSlice())
	stringCon(buf, kind.InterfaceDesc, p.InterfaceDescs().ToSlice())
	stringCon(buf, kind.InterfaceInst, p.InterfaceInsts().ToSlice())
	stringCon(buf, kind.MethodInst, p.MethodInsts().ToSlice())
	stringCon(buf, kind.Method, p.Methods().ToSlice())
	stringCon(buf, kind.Metrics, p.Metrics().ToSlice())
	stringCon(buf, kind.ObjectInst, p.ObjectInsts().ToSlice())
	stringCon(buf, kind.Object, p.Objects().ToSlice())
	stringCon(buf, kind.Package, p.Packages().ToSlice())
	stringCon(buf, kind.Selection, p.Selections().ToSlice())
	stringCon(buf, kind.Signature, p.Signatures().ToSlice())
	stringCon(buf, kind.StructDesc, p.StructDescs().ToSlice())
	stringCon(buf, kind.TempDeclRef, p.TempDeclRefs().ToSlice())
	stringCon(buf, kind.TempReference, p.TempReferences().ToSlice())
	stringCon(buf, kind.TempTypeParamRef, p.TempTypeParamRefs().ToSlice())
	stringCon(buf, kind.TypeParam, p.TypeParams().ToSlice())
	stringCon(buf, kind.Value, p.Values().ToSlice())
	return buf.String()
}
