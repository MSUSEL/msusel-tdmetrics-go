package converter

import (
	"fmt"
	"go/types"

	"github.com/Snow-Gremlin/goToolbox/terrors/terror"
	"github.com/Snow-Gremlin/goToolbox/utils"

	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/abstractor/baker"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/abstractor/instantiator"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/abstractor/querier"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs/hint"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs/innate"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs/kind"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/logger"
)

type Converter interface {
	Nest() constructs.NestType
	ImplicitTypes() []constructs.TypeDesc
	ConvertType(t types.Type, context string) constructs.TypeDesc
	ConvertSignature(t *types.Signature, context string) constructs.Signature
	ConvertInstanceTypes(t *types.TypeList, context string) []constructs.TypeDesc
}

func New(
	log *logger.Logger,
	querier *querier.Querier,
	baker baker.Baker,
	proj constructs.Project,
	curPkg constructs.Package,
	nest constructs.NestType,
	implicitTypes []constructs.TypeDesc,
	tpReplacer map[*types.TypeParam]*types.TypeParam,
	typeCache map[any]any,
) Converter {
	log2 := log.Group(`converter`).Prefix(`|  `)

	return &convImp{
		log:           log2,
		querier:       querier,
		baker:         baker,
		proj:          proj,
		curPkg:        curPkg,
		nest:          nest,
		implicitTypes: implicitTypes,
		tpReplacer:    tpReplacer,
		typeCache:     typeCache,
	}
}

type convImp struct {
	log           *logger.Logger
	querier       *querier.Querier
	baker         baker.Baker
	proj          constructs.Project
	curPkg        constructs.Package
	nest          constructs.NestType
	implicitTypes []constructs.TypeDesc
	tpReplacer    map[*types.TypeParam]*types.TypeParam
	context       string
	tpSeen        map[string]constructs.TempTypeParamRef
	typeCache     map[any]any
}

func (c *convImp) Nest() constructs.NestType            { return c.nest }
func (c *convImp) ImplicitTypes() []constructs.TypeDesc { return c.implicitTypes }

func (c *convImp) ConvertType(t types.Type, context string) constructs.TypeDesc {
	c.log.Logf("convert type: %v", t)
	c.tpSeen = map[string]constructs.TempTypeParamRef{}
	c.context = context
	t2 := cache(c, t, c.convertType)
	c.tpSeen = nil
	c.log.Logf("|  result: %v", t2)
	return t2
}

func (c *convImp) ConvertSignature(t *types.Signature, context string) constructs.Signature {
	c.log.Logf("convert signature: %v", t)
	c.tpSeen = map[string]constructs.TempTypeParamRef{}
	c.context = context
	t2 := cache(c, t, c.convertSignature)
	c.tpSeen = nil
	c.log.Logf("|  result: %v", t2)
	return t2
}

func (c *convImp) ConvertInstanceTypes(t *types.TypeList, context string) []constructs.TypeDesc {
	c.log.Logf("convert instance types: %v", t)
	c.tpSeen = map[string]constructs.TempTypeParamRef{}
	c.context = context
	t2 := cache(c, t, c.convertInstanceTypes)
	c.tpSeen = nil
	c.log.Logf("|  result: %v", t2)
	return t2
}

func cache[T, R any](c *convImp, t T, handle func(T) R) R {
	if cached, ok := c.typeCache[t]; ok {
		if v, ok := cached.(R); ok {
			return v
		}
		panic(fmt.Errorf("expected cached value for %v to be %T but got %T", t, cached, utils.Zero[R]()))
	}
	t2 := handle(t)
	c.typeCache[t] = t2
	return t2
}

func (c *convImp) convertType(t types.Type) constructs.TypeDesc {
	switch t2 := t.(type) {
	case *types.Alias:
		return cache(c, t2, c.convertAlias)
	case *types.Array:
		return cache(c, t2, c.convertArray)
	case *types.Basic:
		return cache(c, t2, c.convertBasic)
	case *types.Chan:
		return cache(c, t2, c.convertChan)
	case *types.Interface:
		return cache(c, t2, c.convertInterface)
	case *types.Map:
		return cache(c, t2, c.convertMap)
	case *types.Named:
		return cache(c, t2, c.convertNamed)
	case *types.Pointer:
		return cache(c, t2, c.convertPointer)
	case *types.Signature:
		return cache(c, t2, c.convertSignature)
	case *types.Slice:
		return cache(c, t2, c.convertSlice)
	case *types.Struct:
		return cache(c, t2, c.convertStruct)
	case *types.TypeParam:
		return cache(c, t2, c.convertTypeParam)
	case *types.Union:
		return cache(c, t2, c.convertUnion)
	default:
		panic(terror.New(`unhandled type`).
			WithType(`type`, t).
			With(`value`, t))
	}
}

func (c *convImp) convertAlias(t *types.Alias) constructs.TypeDesc {
	return cache(c, t.Rhs(), c.convertType)
}

func (c *convImp) convertArray(t *types.Array) constructs.TypeDesc {
	elem := cache(c, t.Elem(), c.convertType)
	generic := c.baker.BakeList()
	return instantiator.InterfaceDecl(c.log, c.proj, t.Underlying(), generic, nil, []constructs.TypeDesc{elem})
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
	elem := cache(c, t.Elem(), c.convertType)
	generic := c.baker.BakeChan()
	return instantiator.InterfaceDecl(c.log, c.proj, t.Underlying(), generic, nil, []constructs.TypeDesc{elem})
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
		sig := cache(c, f.Type().(*types.Signature), c.convertSignature)
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
	key := cache(c, t.Key(), c.convertType)
	value := cache(c, t.Elem(), c.convertType)
	generic := c.baker.BakeMap()
	return instantiator.InterfaceDecl(c.log, c.proj, t.Underlying(), generic, nil, []constructs.TypeDesc{key, value})
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
	instanceTypes := cache(c, t.TypeArgs(), c.convertInstanceTypes)

	// Check if the reference can already be found.
	typ, found := c.proj.FindType(pkgPath, name, c.nest, c.implicitTypes, instanceTypes, true, false)
	if !found {
		// Otherwise, create a temporary reference that will be filled later.
		return c.proj.NewTempReference(constructs.TempReferenceArgs{
			RealType:      t,
			PackagePath:   pkgPath,
			Name:          name,
			Nest:          c.nest,
			ImplicitTypes: c.implicitTypes,
			InstanceTypes: instanceTypes,
			Package:       c.curPkg.Source(),
		})
	}

	switch typ.Kind() {
	case kind.TempReference:
		return typ
	case kind.InterfaceDecl:
		return instantiator.InterfaceDecl(c.log, c.proj, t.Underlying(), typ.(constructs.InterfaceDecl), c.implicitTypes, instanceTypes)
	case kind.Object:
		return instantiator.Object(c.log, c.proj, t.Underlying(), typ.(constructs.Object), c.implicitTypes, instanceTypes)
	default:
		panic(terror.New(`unexpected declaration type`).
			With(`kind`, typ.Kind()).
			With(`decl`, typ))
	}
}

func (c *convImp) convertPointer(t *types.Pointer) constructs.TypeDesc {
	elem := cache(c, t.Elem(), c.convertType)
	generic := c.baker.BakePointer()
	return instantiator.InterfaceDecl(c.log, c.proj, t.Underlying(), generic, nil, []constructs.TypeDesc{elem})
}

func (c *convImp) convertSignature(t *types.Signature) constructs.Signature {
	// Don't output receiver or receiver type here.
	// Don't convert type parameters here.
	t2 := c.proj.NewSignature(constructs.SignatureArgs{
		RealType: t,
		Variadic: t.Variadic(),
		Params:   cache(c, t.Params(), c.convertArguments),
		Results:  cache(c, t.Results(), c.convertArguments),
		Package:  c.curPkg.Source(),
	})
	return t2
}

func (c *convImp) convertSlice(t *types.Slice) constructs.TypeDesc {
	elem := cache(c, t.Elem(), c.convertType)
	generic := c.baker.BakeList()
	return instantiator.InterfaceDecl(c.log, c.proj, t.Underlying(), generic, nil, []constructs.TypeDesc{elem})
}

func (c *convImp) convertStruct(t *types.Struct) constructs.StructDesc {
	fields := make([]constructs.Field, 0, t.NumFields())
	for i := range t.NumFields() {
		f := t.Field(i)
		if !constructs.BlankName(f.Name()) {
			field := c.proj.NewField(constructs.FieldArgs{
				Name:     f.Name(),
				Exported: f.Exported(),
				Type:     cache(c, f.Type(), c.convertType),
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
			Type: cache(c, t2.Type(), c.convertType),
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
		it := cache(c, term.Type(), c.convertType)
		if term.Tilde() {
			approx = append(approx, it)
		} else {
			exact = append(exact, it)
		}
	}
	return exact, approx
}

func (c *convImp) convertTypeParam(t *types.TypeParam) constructs.TypeDesc {
	if tr, ok := c.tpReplacer[t]; ok {
		t = tr
	}

	name := t.Obj().Name()
	if ref, ok := c.tpSeen[name]; ok {
		if ref == nil {
			ref = c.proj.NewTempTypeParamRef(constructs.TempTypeParamRefArgs{
				RealType: t,
				Context:  c.context,
				Name:     name,
			})
			c.tpSeen[name] = ref
		}

		// TODO: Make sure that the references are being replaced?
		return ref
	}

	c.tpSeen[name] = nil
	t2 := t.Obj().Type().Underlying()
	tpDesc := c.proj.NewTypeParam(constructs.TypeParamArgs{
		Name: name,
		Type: cache(c, t2, c.convertType),
	})

	if ref, ok := c.tpSeen[name]; ok && ref != nil {
		ref.SetResolution(tpDesc)
	}

	return tpDesc
}

func (c *convImp) convertInstanceTypes(t *types.TypeList) []constructs.TypeDesc {
	list := make([]constructs.TypeDesc, t.Len())
	for i := range t.Len() {
		list[i] = cache(c, t.At(i), c.convertType)
	}
	return list
}
