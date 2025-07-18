package project

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/Snow-Gremlin/goToolbox/collections"
	"github.com/Snow-Gremlin/goToolbox/collections/enumerator"
	"github.com/Snow-Gremlin/goToolbox/comp"
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

func (p *projectImp) FindType(pkgPath, name string, nest constructs.NestType,
	implicitTypes, instanceTypes []constructs.TypeDesc,
	allowRef, panicOnNotFound bool) (constructs.TypeDesc, bool) {

	assert.ArgNotEmpty(`pkgPath`, pkgPath)

	if allowRef {
		ref, found := p.TempReferences().Enumerate().Where(func(ref constructs.TempReference) bool {
			return comp.Or(
				comp.DefaultPend(ref.PackagePath(), pkgPath),
				comp.DefaultPend(ref.Name(), name),
				constructs.SliceComparerPend(ref.ImplicitTypes(), implicitTypes),
				constructs.SliceComparerPend(ref.InstanceTypes(), instanceTypes),
				constructs.ComparerPend(ref.Nest(), nest),
			) == 0
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
			With(`nest`, nest).
			With(`implicit types`, implicitTypes).
			With(`instance types`, instanceTypes).
			With(`package path`, pkgPath).
			With(`existing paths`, `[`+names+`]`))
	}

	decl := pkg.FindTypeDecl(name, nest)
	if decl == nil {
		if !panicOnNotFound {
			return nil, false
		}
		panic(terror.New(`failed to find type for type reference`).
			With(`type name`, name).
			With(`nest`, nest).
			With(`implicit types`, implicitTypes).
			With(`instance types`, instanceTypes).
			With(`package path`, pkgPath))
	}

	if len(implicitTypes) > 0 || len(instanceTypes) > 0 {
		switch t := decl.(type) {
		case constructs.InterfaceDecl:
			if inst, ok := t.FindInstance(implicitTypes, instanceTypes); ok {
				return inst, true
			}
			return nil, false

		case constructs.Object:
			if inst, ok := t.FindInstance(implicitTypes, instanceTypes); ok {
				return inst, true
			}
			return nil, false

		case constructs.Method:
			panic(terror.New(`can not use method instance as type for type reference`).
				With(`type name`, name).
				With(`nest`, nest).
				With(`implicit types`, implicitTypes).
				With(`instance types`, instanceTypes).
				With(`method`, t).
				With(`package path`, pkgPath))

		case constructs.Value:
			panic(terror.New(`can not get an instance of a value for type reference`).
				With(`type name`, name).
				With(`nest`, nest).
				With(`implicit types`, implicitTypes).
				With(`instance types`, instanceTypes).
				With(`value`, t).
				With(`package path`, pkgPath))

		default:
			panic(terror.New(`unexpected type for type reference instance`).
				With(`type name`, name).
				With(`nest`, nest).
				With(`implicit types`, implicitTypes).
				With(`instance types`, instanceTypes).
				With(`declaration`, decl).
				With(`package path`, pkgPath))
		}
	}

	return decl, true
}

func (p *projectImp) FindDecl(pkgPath, name string, nest constructs.NestType,
	implicitTypes, instanceTypes []constructs.TypeDesc,
	allowRef, panicOnNotFound bool) (constructs.Construct, bool) {

	assert.ArgNotEmpty(`pkgPath`, pkgPath)

	if allowRef {
		ref, found := p.TempDeclRefs().Enumerate().Where(func(ref constructs.TempDeclRef) bool {
			return comp.Or(
				comp.DefaultPend(ref.PackagePath(), pkgPath),
				comp.DefaultPend(ref.Name(), name),
				constructs.SliceComparerPend(ref.ImplicitTypes(), implicitTypes),
				constructs.SliceComparerPend(ref.InstanceTypes(), instanceTypes),
				constructs.ComparerPend(ref.Nest(), nest),
			) == 0
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

	decl := pkg.FindDecl(name, nest)
	if decl == nil {
		if !panicOnNotFound {
			return nil, false
		}
		panic(terror.New(`failed to find declaration for declaration reference`).
			With(`decl name`, name).
			With(`package path`, pkgPath))
	}

	if len(implicitTypes) > 0 || len(instanceTypes) > 0 {
		switch t := decl.(type) {
		case constructs.InterfaceDecl:
			if inst, ok := t.FindInstance(implicitTypes, instanceTypes); ok {
				return inst, true
			}
			return nil, false

		case constructs.Object:
			if inst, ok := t.FindInstance(implicitTypes, instanceTypes); ok {
				return inst, true
			}
			return nil, false

		case constructs.Method:
			assert.ArgIsNil(`nest`, nest)
			assert.ArgIsNil(`implicit types`, implicitTypes)
			if inst, ok := t.FindInstance(instanceTypes); ok {
				return inst, true
			}
			return nil, false

		case constructs.Value:
			panic(terror.New(`can not get an instance of a value for declaration reference`).
				With(`type name`, name).
				With(`nest`, nest).
				With(`implicit types`, implicitTypes).
				With(`instance types`, instanceTypes).
				With(`value`, t).
				With(`package path`, pkgPath))

		default:
			panic(terror.New(`unexpected declaration for declaration reference instance`).
				With(`type name`, name).
				With(`nest`, nest).
				With(`implicit types`, implicitTypes).
				With(`instance types`, instanceTypes).
				With(`declaration`, decl).
				With(`package path`, pkgPath))
		}
	}

	return decl, true
}

func (p *projectImp) UpdateIndices(skipDead bool) {
	var index int
	var kind kind.Kind
	var prev constructs.Construct
	p.AllConstructs().Foreach(func(c constructs.Construct) {
		if cKind := c.Kind(); kind != cKind {
			kind = cKind
			index = 0
			prev = nil
		}
		alive := c.Alive() || !skipDead
		duplicate := true
		if constructs.ComparerPend(prev, c)() != 0 {
			if alive {
				index++
			}
			duplicate = false
		}
		if alive {
			c.SetIndex(index, duplicate)
		} else {
			c.SetIndex(0, duplicate)
		}
		prev = c
	})
}

func pl(k kind.Kind) string {
	s := string(k)
	if !strings.HasSuffix(s, `s`) {
		s += `s`
	}
	return s
}

func jsonList[T constructs.Construct](m *jsonify.Map, ctx *jsonify.Context, k kind.Kind, s collections.ReadonlySortedSet[T]) *jsonify.Map {
	list := s.Enumerate().Where(func(c T) bool { return !c.Duplicate() }).ToSlice()
	return m.AddNonZero(ctx, pl(k), list)
}

func (p *projectImp) ToJson(ctx *jsonify.Context) jsonify.Datum {
	m := jsonify.NewMap().
		Add(ctx, `language`, `go`).
		AddNonZero(ctx, `locs`, p.locations)

	jsonList(m, ctx, kind.Abstract, p.Abstracts())
	jsonList(m, ctx, kind.Argument, p.Arguments())
	jsonList(m, ctx, kind.Field, p.Fields())
	jsonList(m, ctx, kind.Package, p.Packages())
	jsonList(m, ctx, kind.Metrics, p.Metrics())
	jsonList(m, ctx, kind.Selection, p.Selections())
	jsonList(m, ctx, kind.TempDeclRef, p.TempDeclRefs())

	jsonList(m, ctx, kind.InterfaceDecl, p.InterfaceDecls())
	jsonList(m, ctx, kind.Method, p.Methods())
	jsonList(m, ctx, kind.Object, p.Objects())
	jsonList(m, ctx, kind.Value, p.Values())

	jsonList(m, ctx, kind.Basic, p.Basics())
	jsonList(m, ctx, kind.InterfaceDesc, p.InterfaceDescs())
	jsonList(m, ctx, kind.InterfaceInst, p.InterfaceInsts())
	jsonList(m, ctx, kind.MethodInst, p.MethodInsts())
	jsonList(m, ctx, kind.ObjectInst, p.ObjectInsts())
	jsonList(m, ctx, kind.TempReference, p.TempReferences())
	jsonList(m, ctx, kind.TempTypeParamRef, p.TempTypeParamRefs())
	jsonList(m, ctx, kind.Signature, p.Signatures())
	jsonList(m, ctx, kind.StructDesc, p.StructDescs())
	jsonList(m, ctx, kind.TypeParam, p.TypeParams())

	return m
}

func stringCon[T fmt.Stringer, S ~[]T](buf *strings.Builder, k kind.Kind, s S) {
	buf.WriteString(pl(k))
	if len(s) <= 0 {
		buf.WriteString(" { }\n")
		return
	}
	buf.WriteString(" {\n")
	for i, k := range s {
		buf.WriteString(fmt.Sprintf("  %2d. %q\n", i+1, k.String()))
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
