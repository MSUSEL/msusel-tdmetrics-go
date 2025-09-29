package project

import (
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

func (p *projectImp) Kind() kind.Kind { return kind.Project }

func (p *projectImp) Locs() locs.Set { return p.locations }

func (p *projectImp) Factories() collections.Enumerator[constructs.Factory] {
	return enumerator.Enumerate[constructs.Factory](
		p.AbstractFactory,
		p.ArgumentFactory,
		p.BasicFactory,
		p.FieldFactory,
		p.InterfaceDeclFactory,
		p.InterfaceDescFactory,
		p.InterfaceInstFactory,
		p.MethodFactory,
		p.MethodInstFactory,
		p.MetricsFactory,
		p.ObjectFactory,
		p.ObjectInstFactory,
		p.PackageFactory,
		p.SelectionFactory,
		p.SignatureFactory,
		p.StructDescFactory,
		p.TypeParamFactory,
		p.TempDeclRefFactory,
		p.TempReferenceFactory,
		p.TempTypeParamRefFactory,
		p.ValueFactory,
	)
}

func (p *projectImp) Enumerate() collections.Enumerator[constructs.Construct] {
	return enumerator.Expand(p.Factories(), func(f constructs.Factory) collections.Iterable[constructs.Construct] {
		return f.Enumerate().Iterate
	})
}

func (p *projectImp) Refresh() {
	p.Factories().Foreach(constructs.Factory.Refresh)
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

func (p *projectImp) RemoveDuplicates() {
	for {
		m := map[constructs.Construct]constructs.Construct{}
		var prev constructs.Construct
		for c := range p.Enumerate().Seq() {
			if c.Duplicate() {
				// Skip over any constructs already labelled as duplicate.
				continue
			}

			if constructs.ComparerPend(prev, c)() == 0 {
				c.SetDuplicate(true)
				m[c] = prev
				// Don't set `prev` so that we use the same construct for replacement.
				continue
			}

			prev = c
		}
		if len(m) <= 0 {
			// No new duplicates were found so exit.
			break
		}

		// TODO: Use an any to find Constructs.

		// TODO: Need the factories to resort and remove duplicates.

	}
}

func (p *projectImp) UpdateIndices(skipDead bool) {
	for f := range p.Factories().Seq() {
		index := 0
		for c := range f.Enumerate().Seq() {
			if !c.Duplicate() && (c.Alive() || !skipDead) {
				index++
				c.SetIndex(index)
			} else {
				c.SetIndex(0)
			}
		}
	}
}

func (p *projectImp) ToJson(ctx *jsonify.Context) jsonify.Datum {
	m := jsonify.NewMap().
		Add(ctx, `language`, `go`).
		AddNonZero(ctx, `locs`, p.locations)
	for f := range p.Factories().Seq() {
		list := f.Enumerate().WhereNot(constructs.Construct.Duplicate).ToSlice()
		m.AddNonZero(ctx, f.Kind().Plural(), list)
	}
	return m
}

func (p *projectImp) String() string {
	buf := &strings.Builder{}
	for f := range p.Factories().Seq() {
		buf.WriteString(f.String() + "\n")
	}
	return buf.String()
}
