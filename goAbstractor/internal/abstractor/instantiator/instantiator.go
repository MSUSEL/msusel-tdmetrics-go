package instantiator

import (
	"go/types"
	"slices"

	"github.com/Snow-Gremlin/goToolbox/terrors/terror"
	"github.com/Snow-Gremlin/goToolbox/utils"

	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/abstractor/querier"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/assert"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs/kind"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/logger"
)

func InterfaceDecl(log *logger.Logger, querier *querier.Querier, proj constructs.Project, realType types.Type,
	decl constructs.InterfaceDecl, implicitTypes, instanceTypes []constructs.TypeDesc) constructs.TypeDesc {

	defer func() {
		if r := recover(); r != nil {
			panic(terror.New(`error in instantiator.InterfaceDecl`, terror.RecoveredPanic(r)).
				With(`realType`, realType).
				With(`decl`, decl).
				With(`implicitTypes`, implicitTypes).
				With(`instanceTypes`, instanceTypes))
		}
	}()

	assert.ArgNotNil(`project`, proj)
	assert.ArgNotNil(`interface decl`, decl)
	assert.ArgHasNoNils(`implicit types`, implicitTypes)
	assert.ArgHasNoNils(`instance types`, instanceTypes)

	itp := decl.ImplicitTypeParams()
	tp := decl.TypeParams()
	assert.ArgsHaveSameLength(`implicit types`, implicitTypes, itp)
	assert.ArgsHaveSameLength(`instance types`, instanceTypes, tp)

	log2 := log.Group(`instantiator`).Prefix(`|  `)
	i, existing, needsInstance := newInstantiator(log2, querier, proj, decl, itp, tp, implicitTypes, instanceTypes)
	if !needsInstance {
		return decl
	}
	if !utils.IsNil(existing) {
		return existing.(constructs.InterfaceInst)
	}

	assert.ArgNotNil(`real type`, realType)
	assert.ArgNotEmpty(`instance types`, instanceTypes)
	assert.ArgHasNoNils(`instance types`, instanceTypes)
	return i.createInstance(realType).(constructs.InterfaceInst)
}

func Object(log *logger.Logger, querier *querier.Querier, proj constructs.Project, realType types.Type,
	decl constructs.Object, implicitTypes, instanceTypes []constructs.TypeDesc) constructs.TypeDesc {

	defer func() {
		if r := recover(); r != nil {
			panic(terror.New(`error in instantiator.Object`, terror.RecoveredPanic(r)).
				With(`realType`, realType).
				With(`decl`, decl).
				With(`implicitTypes`, implicitTypes).
				With(`instanceTypes`, instanceTypes))
		}
	}()

	assert.ArgNotNil(`project`, proj)
	assert.ArgNotNil(`object`, decl)
	assert.ArgHasNoNils(`implicit types`, implicitTypes)
	assert.ArgHasNoNils(`instance types`, instanceTypes)

	itp := decl.ImplicitTypeParams()
	tp := decl.TypeParams()
	assert.ArgsHaveSameLength(`implicit types`, implicitTypes, itp)
	assert.ArgsHaveSameLength(`instance types`, instanceTypes, tp)

	log2 := log.Group(`instantiator`).Prefix(`|  `)
	i, existing, needsInstance := newInstantiator(log2, querier, proj, decl, itp, tp, implicitTypes, instanceTypes)
	if !needsInstance {
		return decl
	}
	if !utils.IsNil(existing) {
		return existing.(constructs.ObjectInst)
	}

	assert.ArgNotNil(`real type`, realType)
	assert.ArgNotEmpty(`instance types`, instanceTypes)
	assert.ArgHasNoNils(`instance types`, instanceTypes)
	return i.createInstance(realType).(constructs.ObjectInst)
}

func Method(log *logger.Logger, querier *querier.Querier, proj constructs.Project,
	decl constructs.Method, instanceTypes []constructs.TypeDesc) constructs.Construct {

	assert.ArgNotNil(`project`, proj)
	assert.ArgNotNil(`method`, decl)

	log2 := log.Group(`instantiator`).Prefix(`|  `)
	typeParams := decl.TypeParams()
	if decl.HasReceiver() {
		typeParams = decl.Receiver().TypeParams()
	}

	i, existing, needsInstance := newInstantiator(log2, querier, proj, decl, nil, typeParams, nil, instanceTypes)
	if !needsInstance {
		return decl
	}
	if !utils.IsNil(existing) {
		return existing.(constructs.MethodInst)
	}

	assert.ArgNotEmpty(`instance types`, instanceTypes)
	assert.ArgHasNoNils(`instance types`, instanceTypes)
	return i.createInstance(nil).(constructs.MethodInst)
}

type instantiator struct {
	log           *logger.Logger
	querier       *querier.Querier
	prior         *instantiator
	proj          constructs.Project
	decl          constructs.Declaration
	implicitTypes []constructs.TypeDesc
	instanceTypes []constructs.TypeDesc
	conversion    map[constructs.TypeParam]constructs.TypeDesc
}

func newInstantiator(log *logger.Logger, querier *querier.Querier, proj constructs.Project,
	decl constructs.Declaration, nestTypeParams, typeParams []constructs.TypeParam,
	implicitTypes, instanceTypes []constructs.TypeDesc) (*instantiator, constructs.Construct, bool) {

	if len(nestTypeParams) != len(implicitTypes) {
		panic(terror.New(`the amount of nested type params must match the implicit instance types`).
			With(`declaration`, decl).
			With(`nested type params count`, len(nestTypeParams)).
			With(`implicit types count`, len(implicitTypes)).
			With(`nested type params`, nestTypeParams).
			With(`implicit types`, implicitTypes).
			With(`type params`, typeParams).
			With(`instance types`, instanceTypes))
	}

	if len(typeParams) != len(instanceTypes) {
		panic(terror.New(`the amount of type params must match the instance types`).
			With(`declaration`, decl).
			With(`type params count`, len(typeParams)).
			With(`instance types count`, len(instanceTypes)).
			With(`nested type params`, nestTypeParams).
			With(`implicit types`, implicitTypes).
			With(`type params`, typeParams).
			With(`instance types`, instanceTypes))
	}

	diff := false
	cmp := constructs.Comparer[constructs.TypeDesc]()
	for i := len(implicitTypes) - 1; i >= 0; i-- {
		if cmp(nestTypeParams[i], implicitTypes[i]) != 0 {
			diff = true
			break
		}
	}
	if !diff {
		for i := len(instanceTypes) - 1; i >= 0; i-- {
			if cmp(typeParams[i], instanceTypes[i]) != 0 {
				diff = true
				break
			}
		}
	}
	if !diff {
		panic(terror.New(`attempted to make an instance with the general's type parameters`))
		// TODO: Add more info
	}

	log.Logf(`instantiating %v`, decl)
	log.Logf(`|- with %v`, instanceTypes)

	// Check declaration is a generic type, leave otherwise.
	if len(nestTypeParams) <= 0 && len(typeParams) <= 0 {
		log.Logf(`'- not generic`)
		return nil, nil, false
	}

	// Check if instance types match the declaration types.
	// For example if `func Foo[T any]() { ... Foo[T]() ... }` is given such that
	// the call to `Foo` doesn't need an instance since it matches the generic.
	tpSame := func(a constructs.TypeParam, b constructs.TypeDesc) bool { return a == b }
	same := slices.EqualFunc(nestTypeParams, instanceTypes, tpSame) &&
		slices.EqualFunc(typeParams, implicitTypes, tpSame)
	if same {
		log.Logf(`'- instantiation has same type arguments as type parameters`)
		return nil, nil, false
	}

	// Check that a matching instance doesn't already exist.
	var instance constructs.Construct
	found := false
	switch decl.Kind() {
	case kind.InterfaceDecl:
		instance, found = decl.(constructs.InterfaceDecl).FindInstance(implicitTypes, instanceTypes)
	case kind.Object:
		instance, found = decl.(constructs.Object).FindInstance(implicitTypes, instanceTypes)
	case kind.Method:
		instance, found = decl.(constructs.Method).FindInstance(instanceTypes)
	}
	if found {
		log.Logf(`'- instantiation found`)
		return nil, instance, true
	}

	// Create a new instantiator for the new instance.
	conversion := make(map[constructs.TypeParam]constructs.TypeDesc, len(nestTypeParams)+len(typeParams))
	for i, tp := range nestTypeParams {
		conversion[tp] = implicitTypes[i]
	}
	for i, tp := range typeParams {
		conversion[tp] = instanceTypes[i]
	}
	log.Logf(`|- instantiation needed`)
	return &instantiator{
		log:           log,
		querier:       querier,
		proj:          proj,
		decl:          decl,
		implicitTypes: implicitTypes,
		instanceTypes: instanceTypes,
		conversion:    conversion,
	}, nil, true
}

func (i *instantiator) TypeDesc(td constructs.TypeDesc) constructs.TypeDesc {
	switch td.Kind() {
	case kind.Basic:
		return td
	case kind.InterfaceDesc:
		return i.InterfaceDesc(td.(constructs.InterfaceDesc))
	case kind.InterfaceInst:
		return i.InterfaceInst(td.(constructs.InterfaceInst))
	case kind.ObjectInst:
		return i.ObjectInst(td.(constructs.ObjectInst))
	case kind.Signature:
		return i.Signature(td.(constructs.Signature))
	case kind.StructDesc:
		return i.StructDesc(td.(constructs.StructDesc))
	case kind.TempReference:
		return i.TempReference(td.(constructs.TempReference))
	case kind.TypeParam:
		return i.TypeParam(td.(constructs.TypeParam))
	case kind.InterfaceDecl:
		return i.InterfaceDecl(td.(constructs.InterfaceDecl))
	case kind.Object:
		return i.Object(td.(constructs.Object))
	default:
		panic(terror.New(`unexpected type description kind`).
			With(`kind`, td.Kind()).
			With(`type desc`, td))
	}
}

func applyToSlice[T any, S ~[]T](s S, handle func(T) T) S {
	result := make(S, len(s))
	for i, e := range s {
		result[i] = handle(e)
	}
	return result
}

func (i *instantiator) Abstract(a constructs.Abstract) constructs.Abstract {
	a2 := i.proj.NewAbstract(constructs.AbstractArgs{
		Name:      a.Name(),
		Exported:  a.Exported(),
		Signature: i.Signature(a.Signature()),
	})
	i.log.Logf(`|- create Abstract: %v`, a2)
	return a2
}

func (i *instantiator) Argument(a constructs.Argument) constructs.Argument {
	a2 := i.proj.NewArgument(constructs.ArgumentArgs{
		Name: a.Name(),
		Type: i.TypeDesc(a.Type()),
	})
	i.log.Logf(`|- create Argument: %v`, a2)
	return a2
}

func (i *instantiator) Field(f constructs.Field) constructs.Field {
	f2 := i.proj.NewField(constructs.FieldArgs{
		Name:     f.Name(),
		Exported: f.Exported(),
		Type:     i.TypeDesc(f.Type()),
		Embedded: f.Embedded(),
	})
	i.log.Logf(`|- create Field: %v`, f2)
	return f2
}

func (i *instantiator) InterfaceInst(in constructs.InterfaceInst) constructs.TypeDesc {
	decl := in.Generic()
	in2 := i.typeDecl(decl, decl.ImplicitTypeParams(), decl.TypeParams(), in.ImplicitTypes(), in.InstanceTypes())
	i.log.Logf(`|- create InterfaceInst: %v`, in2)
	return in2
}

func (i *instantiator) ObjectInst(in constructs.ObjectInst) constructs.TypeDesc {
	decl := in.Generic()
	in2 := i.typeDecl(decl, decl.ImplicitTypeParams(), decl.TypeParams(), in.ImplicitTypes(), in.InstanceTypes())
	i.log.Logf(`|- create ObjectInst: %v`, in2)
	return in2
}

func (i *instantiator) InterfaceDecl(decl constructs.InterfaceDecl) constructs.TypeDesc {
	implicitTypes := constructs.Cast[constructs.TypeDesc](decl.ImplicitTypeParams())
	instanceTypes := constructs.Cast[constructs.TypeDesc](decl.TypeParams())
	decl2 := i.typeDecl(decl, decl.ImplicitTypeParams(), decl.TypeParams(), implicitTypes, instanceTypes)
	i.log.Logf(`|- create InterfaceDecl: %v`, decl2)
	return decl2
}

func (i *instantiator) Object(decl constructs.Object) constructs.TypeDesc {
	implicitTypes := constructs.Cast[constructs.TypeDesc](decl.ImplicitTypeParams())
	instanceTypes := constructs.Cast[constructs.TypeDesc](decl.TypeParams())
	decl2 := i.typeDecl(decl, decl.ImplicitTypeParams(), decl.TypeParams(), implicitTypes, instanceTypes)
	i.log.Logf(`|- create Object: %v`, decl2)
	return decl2
}

func (i *instantiator) getInstanceTypeChange(tps []constructs.TypeDesc) ([]constructs.TypeDesc, bool) {
	anyReplaced := false
	its := make([]constructs.TypeDesc, len(tps))
	for j, td := range tps {
		its[j] = td
		td2 := i.TypeDesc(td)
		if td2.CompareTo(td) != 0 {
			its[j] = td2
			anyReplaced = true
		}
	}
	return its, anyReplaced
}

func (i *instantiator) inProgress(decl constructs.TypeDecl, implicitTypes, instanceTypes []constructs.TypeDesc) bool {
	for !utils.IsNil(i) {
		if i.decl == decl &&
			slices.Equal(i.implicitTypes, implicitTypes) &&
			slices.Equal(i.instanceTypes, instanceTypes) {
			return true
		}
		i = i.prior
	}
	return false
}

func (i *instantiator) typeDecl(decl constructs.TypeDecl,
	nestTypeParams, typeParams []constructs.TypeParam,
	implicitTypes, instanceTypes []constructs.TypeDesc) constructs.TypeDesc {

	implicitTypes, anyImplicitReplaced := i.getInstanceTypeChange(implicitTypes)
	instanceTypes, anyInstanceReplaced := i.getInstanceTypeChange(instanceTypes)
	if !(anyImplicitReplaced || anyInstanceReplaced) {
		return decl
	}

	// If the declaration is the same as a declaration being instantiated,
	// create a reference to avoid the cycle. Cycles are caused by the type
	// being instantiated being used as part of the type definition
	// directly, e.g. `type Foo[T any] interface { Get() Foo[T]  }` or
	// indirectly. e.g. `type Foo[T any] interface { Children() List[Foo[T]] }`
	if i.inProgress(decl, implicitTypes, instanceTypes) {
		typ, found := i.proj.FindType(decl.Package().Path(), decl.Name(), decl.Nest(), implicitTypes, instanceTypes, true, false)
		if found {
			return typ
		}

		realTyp := i.goType(decl, instanceTypes)
		return i.proj.NewTempReference(constructs.TempReferenceArgs{
			PackagePath:   decl.Package().Path(),
			RealType:      realTyp,
			Name:          decl.Name(),
			ImplicitTypes: implicitTypes,
			InstanceTypes: instanceTypes,
			Nest:          decl.Nest(),
			Package:       decl.Package().Source(),
		})
	}

	i2, existing, _ := newInstantiator(i.log, i.querier, i.proj, decl, nestTypeParams, typeParams, implicitTypes, instanceTypes)
	if !utils.IsNil(existing) {
		return existing.(constructs.TypeDesc)
	}
	i2.prior = i
	realTyp := i.goType(decl, instanceTypes)
	return i2.createInstance(realTyp).(constructs.TypeDesc)
}

func (i *instantiator) createInstance(realType types.Type) constructs.Construct {
	switch i.decl.Kind() {
	case kind.InterfaceDecl:
		d := i.decl.(constructs.InterfaceDecl)
		inst := i.proj.NewInterfaceInst(constructs.InterfaceInstArgs{
			RealType:      realType,
			Generic:       d,
			Resolved:      i.InterfaceDesc(d.Interface()),
			ImplicitTypes: i.implicitTypes,
			InstanceTypes: i.instanceTypes,
		})
		i.log.Logf(`'- instantiated interface: %v`, inst)
		return inst

	case kind.Object:
		d := i.decl.(constructs.Object)
		obj := i.proj.NewObjectInst(constructs.ObjectInstArgs{
			RealType:      realType,
			Generic:       d,
			ResolvedData:  i.StructDesc(d.Data()),
			ImplicitTypes: i.implicitTypes,
			InstanceTypes: i.instanceTypes,
		})
		i.log.Logf(`'- instantiated object: %v`, obj)
		return obj

	case kind.Method:
		d := i.decl.(constructs.Method)
		md := i.proj.NewMethodInst(constructs.MethodInstArgs{
			Generic:       d,
			Resolved:      i.Signature(d.Signature()),
			InstanceTypes: i.instanceTypes,
			Metrics:       nil, // This needs to be set later
		})
		i.log.Logf(`'- instantiated method: %v`, md)
		return md

	default:
		panic(terror.New(`unexpected declaration type`).
			With(`kind`, i.decl.Kind()).
			With(`decl`, i.decl))
	}
}

func (i *instantiator) InterfaceDesc(it constructs.InterfaceDesc) constructs.InterfaceDesc {
	it2 := i.proj.NewInterfaceDesc(constructs.InterfaceDescArgs{
		Hint:      it.Hint(),
		Abstracts: applyToSlice(it.Abstracts(), i.Abstract),
		Exact:     applyToSlice(it.Exact(), i.TypeDesc),
		Approx:    applyToSlice(it.Approx(), i.TypeDesc),
		Package:   i.decl.Package().Source(),
	})
	i.log.Logf(`|- create InterfaceDesc: %v`, it2)
	return it2
}

func (i *instantiator) goType(desc constructs.TypeDesc, instanceTypes []constructs.TypeDesc) types.Type {
	origin := desc.GoType()
	if origin == nil {
		return nil
	}
	rtp, ok := origin.(interface {
		types.Type
		TypeParams() *types.TypeParamList
	})
	if !ok || rtp.TypeParams().Len() <= 0 {
		return origin
	}

	tArgs := make([]types.Type, len(instanceTypes))
	for i, ta := range instanceTypes {
		tArgs[i] = ta.GoType()
	}

	realTyp, err := types.Instantiate(i.querier.Context(), rtp, tArgs, true)
	if err != nil {
		panic(terror.New(`failed to instantiate declaration Go type`, err))
	}
	return realTyp
}

func (i *instantiator) TempReference(r constructs.TempReference) constructs.TypeDesc {
	if r.Resolved() {
		// The reference will probably not be resolved, but just in case
		// it has been, just return the instantiated resolved type and
		// skip the instantiated reference.
		return i.TypeDesc(r.ResolvedType())
	}

	implicitTypes := applyToSlice(r.ImplicitTypes(), i.TypeDesc)
	instanceTypes := applyToSlice(r.InstanceTypes(), i.TypeDesc)
	typ, found := i.proj.FindType(r.PackagePath(), r.Name(), r.Nest(), implicitTypes, instanceTypes, true, false)
	if found {
		return typ
	}

	realTyp := i.goType(r, instanceTypes)
	r2 := i.proj.NewTempReference(constructs.TempReferenceArgs{
		PackagePath:   r.PackagePath(),
		RealType:      realTyp,
		Name:          r.Name(),
		Nest:          r.Nest(),
		ImplicitTypes: implicitTypes,
		InstanceTypes: instanceTypes,
		Package:       i.decl.Package().Source(),
	})
	i.log.Logf(`|- create TempReference: %v`, r2)
	return r2
}

func (i *instantiator) Signature(s constructs.Signature) constructs.Signature {
	s2 := i.proj.NewSignature(constructs.SignatureArgs{
		Variadic: s.Variadic(),
		Params:   applyToSlice(s.Params(), i.Argument),
		Results:  applyToSlice(s.Results(), i.Argument),
		Package:  i.decl.Package().Source(),
	})
	i.log.Logf(`|- create Signature: %v`, s2)
	return s2
}

func (i *instantiator) StructDesc(s constructs.StructDesc) constructs.StructDesc {
	s2 := i.proj.NewStructDesc(constructs.StructDescArgs{
		Fields:  applyToSlice(s.Fields(), i.Field),
		Package: i.decl.Package().Source(),
	})
	i.log.Logf(`|- create StructDesc: %v`, s2)
	return s2
}

func (i *instantiator) TypeParam(tp constructs.TypeParam) constructs.TypeDesc {
	if t, has := i.conversion[tp]; has {
		return t
	}
	return tp
}
