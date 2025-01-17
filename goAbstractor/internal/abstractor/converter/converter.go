package converter

import (
	"go/types"

	"github.com/Snow-Gremlin/goToolbox/terrors/terror"
	"github.com/Snow-Gremlin/goToolbox/utils"

	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/abstractor/baker"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/abstractor/instantiator"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs/hint"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs/innate"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs/kind"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/logger"
)

type Converter interface {
	ConvertType(t types.Type) constructs.TypeDesc
	ConvertSignature(t *types.Signature) constructs.Signature
	ConvertInstanceTypes(t *types.TypeList) []constructs.TypeDesc
}

func New(
	log *logger.Logger,
	baker baker.Baker,
	proj constructs.Project,
	curPkg constructs.Package,
	tpReplacer map[*types.TypeParam]*types.TypeParam,
) Converter {
	log2 := log.Group(`converter`).Prefix(`|  `)
	return &convImp{
		log:        log2,
		baker:      baker,
		proj:       proj,
		curPkg:     curPkg,
		tpReplacer: tpReplacer,
	}
}

type convImp struct {
	log        *logger.Logger
	baker      baker.Baker
	proj       constructs.Project
	curPkg     constructs.Package
	tpReplacer map[*types.TypeParam]*types.TypeParam
}

func (c *convImp) ConvertType(t types.Type) constructs.TypeDesc {
	c.log.Logf("convert type: %v", t)
	t2 := c.convertType(t)
	c.log.Logf("|  result: %v", t2)
	return t2
}

func (c *convImp) ConvertSignature(t *types.Signature) constructs.Signature {
	c.log.Logf("convert signature: %v", t)
	t2 := c.convertSignature(t)
	c.log.Logf("|  result: %v", t2)
	return t2
}

func (c *convImp) ConvertInstanceTypes(t *types.TypeList) []constructs.TypeDesc {
	c.log.Logf("convert instance types: %v", t)
	t2 := c.convertInstanceTypes(t)
	c.log.Logf("|  result: %v", t2)
	return t2
}

func (c *convImp) convertType(t types.Type) constructs.TypeDesc {
	switch t2 := t.(type) {
	case *types.Alias:
		return c.convertAlias(t2)
	case *types.Array:
		return c.convertArray(t2)
	case *types.Basic:
		return c.convertBasic(t2)
	case *types.Chan:
		return c.convertChan(t2)
	case *types.Interface:
		return c.convertInterface(t2)
	case *types.Map:
		return c.convertMap(t2)
	case *types.Named:
		return c.convertNamed(t2)
	case *types.Pointer:
		return c.convertPointer(t2)
	case *types.Signature:
		return c.convertSignature(t2)
	case *types.Slice:
		return c.convertSlice(t2)
	case *types.Struct:
		return c.convertStruct(t2)
	case *types.TypeParam:
		return c.convertTypeParam(t2)
	case *types.Union:
		return c.convertUnion(t2)
	default:
		panic(terror.New(`unhandled type`).
			WithType(`type`, t).
			With(`value`, t))
	}
}

func (c *convImp) convertAlias(t *types.Alias) constructs.TypeDesc {
	return c.convertType(t.Rhs())
}

func (c *convImp) convertArray(t *types.Array) constructs.TypeDesc {
	elem := c.convertType(t.Elem())
	generic := c.baker.BakeList()
	return instantiator.InterfaceDecl(c.log, c.proj, t.Underlying(), generic, elem)
}

func (c *convImp) convertBasic(t *types.Basic) constructs.TypeDesc {
	switch t.Kind() {
	case types.Complex64:
		return c.baker.BakeComplex64()
	case types.Complex128:
		return c.baker.BakeComplex128()
	default:
		return c.proj.NewBasic(constructs.BasicArgs{
			RealType: t,
		})
	}
}

func (c *convImp) convertChan(t *types.Chan) constructs.TypeDesc {
	elem := c.convertType(t.Elem())
	generic := c.baker.BakeChan()
	return instantiator.InterfaceDecl(c.log, c.proj, t.Underlying(), generic, elem)
}

func (c *convImp) convertInterface(t *types.Interface) constructs.InterfaceDesc {
	t.Complete()

	if t.Empty() { // any
		return c.baker.BakeAny()
	}

	if t.IsComparable() && t.NumMethods() <= 0 && t.NumEmbeddeds() <= 0 {
		return c.baker.BakeComparable().Interface() // comparable
	}

	h := hint.None
	abstracts := []constructs.Abstract{}
	if t.IsComparable() && t.NumMethods() > 0 {
		h = hint.Comparable
		comp := c.baker.BakeComparableAbstract()
		abstracts = append(abstracts, comp)
		c.log.Logf(`add comparable to %v`, t)
	}

	pinned := false
	for i := range t.NumMethods() {
		f := t.Method(i)
		sig := c.convertSignature(f.Type().(*types.Signature))
		abstract := c.proj.NewAbstract(constructs.AbstractArgs{
			Name:      f.Name(),
			Exported:  f.Exported(),
			Signature: sig,
		})
		abstracts = append(abstracts, abstract)
		pinned = pinned || !f.Exported()
	}

	var exact, approx []constructs.TypeDesc
	for i := range t.NumEmbeddeds() {
		et := t.EmbeddedType(i)
		if union, ok := et.(*types.Union); ok {
			exact2, approx2 := c.readUnionTerms(union)
			exact = append(exact, exact2...)
			approx = append(approx, approx2...)
		}
	}

	var pinnedPkg constructs.Package
	if pinned {
		pinnedPkg = c.curPkg
	}

	return c.proj.NewInterfaceDesc(constructs.InterfaceDescArgs{
		Hint:      h,
		RealType:  t,
		PinnedPkg: pinnedPkg,
		Exact:     exact,
		Approx:    approx,
		Abstracts: abstracts,
		Package:   c.curPkg.Source(),
	})
}

func (c *convImp) convertMap(t *types.Map) constructs.TypeDesc {
	key := c.convertType(t.Key())
	value := c.convertType(t.Elem())
	generic := c.baker.BakeMap()
	return instantiator.InterfaceDecl(c.log, c.proj, t.Underlying(), generic, key, value)
}

func (c *convImp) convertNamed(t *types.Named) constructs.TypeDesc {
	pkgPath := ``
	if !utils.IsNil(t.Obj().Pkg()) {
		pkgPath = t.Obj().Pkg().Path()
	}
	name := t.Obj().Name()

	// Check for builtin types that need to be baked.
	if len(pkgPath) <= 0 {
		if typ := c.baker.TypeByName(name); !utils.IsNil(typ) {
			return typ
		}
		pkgPath = innate.Builtin
	}

	// Get any type parameters.
	instanceTp := c.convertInstanceTypes(t.TypeArgs())

	// Check if the reference can already be found.
	typ, found := c.proj.FindType(pkgPath, name, instanceTp, true, false)
	if !found {
		// Otherwise, create a temporary reference that will be filled later.
		return c.proj.NewTempReference(constructs.TempReferenceArgs{
			RealType:      t,
			PackagePath:   pkgPath,
			Name:          name,
			InstanceTypes: instanceTp,
			Package:       c.curPkg.Source(),
		})
	}

	switch typ.Kind() {
	case kind.TempReference:
		return typ
	case kind.InterfaceDecl:
		return instantiator.InterfaceDecl(c.log, c.proj, t.Underlying(), typ.(constructs.InterfaceDecl), instanceTp...)
	case kind.Object:
		return instantiator.Object(c.log, c.proj, t.Underlying(), typ.(constructs.Object), instanceTp...)
	default:
		panic(terror.New(`unexpected declaration type`).
			With(`kind`, typ.Kind()).
			With(`decl`, typ))
	}
}

func (c *convImp) convertPointer(t *types.Pointer) constructs.TypeDesc {
	elem := c.convertType(t.Elem())
	generic := c.baker.BakePointer()
	return instantiator.InterfaceDecl(c.log, c.proj, t.Underlying(), generic, elem)
}

func (c *convImp) convertSignature(t *types.Signature) constructs.Signature {
	// Don't output receiver or receiver type here.
	// Don't convert type parameters here.
	t2 := c.proj.NewSignature(constructs.SignatureArgs{
		RealType: t,
		Variadic: t.Variadic(),
		Params:   c.convertArguments(t.Params()),
		Results:  c.convertArguments(t.Results()),
		Package:  c.curPkg.Source(),
	})
	return t2
}

func (c *convImp) convertSlice(t *types.Slice) constructs.TypeDesc {
	elem := c.convertType(t.Elem())
	generic := c.baker.BakeList()
	return instantiator.InterfaceDecl(c.log, c.proj, t.Underlying(), generic, elem)
}

func (c *convImp) convertStruct(t *types.Struct) constructs.StructDesc {
	fields := make([]constructs.Field, 0, t.NumFields())
	for i := range t.NumFields() {
		f := t.Field(i)
		if !constructs.BlankName(f.Name()) {
			field := c.proj.NewField(constructs.FieldArgs{
				Name:     f.Name(),
				Exported: f.Exported(),
				Type:     c.convertType(f.Type()),
				Embedded: f.Embedded(),
			})
			fields = append(fields, field)
		}
	}

	return c.proj.NewStructDesc(constructs.StructDescArgs{
		RealType: t,
		Fields:   fields,
		Package:  c.curPkg.Source(),
	})
}

func (c *convImp) convertArguments(t *types.Tuple) []constructs.Argument {
	count := t.Len()
	list := make([]constructs.Argument, count)
	for i := range count {
		t2 := t.At(i)
		list[i] = c.proj.NewArgument(constructs.ArgumentArgs{
			Name: t2.Name(),
			Type: c.convertType(t2.Type()),
		})
	}
	return list
}

func (c *convImp) convertUnion(t *types.Union) constructs.InterfaceDesc {
	exact, approx := c.readUnionTerms(t)
	return c.proj.NewInterfaceDesc(constructs.InterfaceDescArgs{
		Exact:   exact,
		Approx:  approx,
		Package: c.curPkg.Source(),
	})
}

func (c *convImp) readUnionTerms(t *types.Union) (exact, approx []constructs.TypeDesc) {
	for i := range t.Len() {
		term := t.Term(i)
		it := c.convertType(term.Type())
		if term.Tilde() {
			approx = append(approx, it)
		} else {
			exact = append(exact, it)
		}
	}
	return exact, approx
}

func (c *convImp) convertTypeParam(t *types.TypeParam) constructs.TypeParam {
	if tr, ok := c.tpReplacer[t]; ok {
		t = tr
	}

	t2 := t.Obj().Type().Underlying()
	return c.proj.NewTypeParam(constructs.TypeParamArgs{
		Name: t.Obj().Name(),
		Type: c.convertType(t2),
	})
}

func (c *convImp) convertInstanceTypes(t *types.TypeList) []constructs.TypeDesc {
	list := make([]constructs.TypeDesc, t.Len())
	for i := range t.Len() {
		list[i] = c.convertType(t.At(i))
	}
	return list
}
