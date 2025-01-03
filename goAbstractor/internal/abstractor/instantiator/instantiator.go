package instantiator

import (
	"go/types"
	"slices"

	"github.com/Snow-Gremlin/goToolbox/terrors/terror"
	"github.com/Snow-Gremlin/goToolbox/utils"

	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/assert"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs/kind"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/logger"
)

func InterfaceDecl(log *logger.Logger, proj constructs.Project, realType types.Type, decl constructs.InterfaceDecl, instanceTypes ...constructs.TypeDesc) constructs.TypeDesc {
	assert.ArgNotNil(`project`, proj)
	assert.ArgNotNil(`interface decl`, decl)

	log2 := log.Group(`instantiator`).Prefix(`|  `)
	i, existing, needsInstance := newInstantiator(log2, proj, decl, decl.TypeParams(), instanceTypes)
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

func Object(log *logger.Logger, proj constructs.Project, realType types.Type, decl constructs.Object, instanceTypes ...constructs.TypeDesc) constructs.TypeDesc {
	assert.ArgNotNil(`project`, proj)
	assert.ArgNotNil(`object`, decl)

	log2 := log.Group(`instantiator`).Prefix(`|  `)
	i, existing, needsInstance := newInstantiator(log2, proj, decl, decl.TypeParams(), instanceTypes)
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

func Method(log *logger.Logger, proj constructs.Project, decl constructs.Method, instanceTypes ...constructs.TypeDesc) constructs.Construct {
	assert.ArgNotNil(`project`, proj)
	assert.ArgNotNil(`method`, decl)

	log2 := log.Group(`instantiator`).Prefix(`|  `)
	typeParams := decl.TypeParams()
	if decl.HasReceiver() {
		typeParams = decl.Receiver().TypeParams()
	}

	i, existing, needsInstance := newInstantiator(log2, proj, decl, typeParams, instanceTypes)
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
	prior         *instantiator
	proj          constructs.Project
	decl          constructs.Declaration
	instanceTypes []constructs.TypeDesc
	conversion    map[constructs.TypeParam]constructs.TypeDesc
}

func newInstantiator(log *logger.Logger, proj constructs.Project, decl constructs.Declaration, typeParams []constructs.TypeParam, instanceTypes []constructs.TypeDesc) (*instantiator, constructs.Construct, bool) {
	count := len(typeParams)
	if count != len(instanceTypes) {
		panic(terror.New(`the amount of type params must match the instance types`).
			With(`declaration`, decl).
			With(`type params count`, count).
			With(`instance types count`, len(instanceTypes)).
			With(`type params`, typeParams).
			With(`instance types`, instanceTypes))
	}

	log.Logf(`instantiating %v`, decl)
	log.Logf(`|- with %v`, instanceTypes)

	// Check declaration is a generic type, leave otherwise.
	if count <= 0 {
		log.Logf(`'- not generic`)
		return nil, nil, false
	}

	// Check if instance types match the declaration types.
	// For example if `func Foo[T any]() { ... Foo[T]() ... }` is given such that
	// the call to `Foo` doesn't need an instance since it matches the generic.
	same := true
	for i, tp := range typeParams {
		if tp != instanceTypes[i] {
			same = false
			break
		}
	}
	if same {
		log.Logf(`'- instantiation has same type arguments as type parameters`)
		return nil, nil, false
	}

	// Check that a matching instance doesn't already exist.
	var instance constructs.Construct
	found := false
	switch decl.Kind() {
	case kind.InterfaceDecl:
		instance, found = decl.(constructs.InterfaceDecl).FindInstance(instanceTypes)
	case kind.Object:
		instance, found = decl.(constructs.Object).FindInstance(instanceTypes)
	case kind.Method:
		instance, found = decl.(constructs.Method).FindInstance(instanceTypes)
	}
	if found {
		log.Logf(`'- instantiation found`)
		return nil, instance, true
	}

	// Create a new instantiator for the new instance.
	conversion := make(map[constructs.TypeParam]constructs.TypeDesc, count)
	for i, tp := range typeParams {
		conversion[tp] = instanceTypes[i]
	}
	log.Logf(`|- instantiation needed`)
	return &instantiator{
		log:           log,
		proj:          proj,
		decl:          decl,
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
	return i.proj.NewAbstract(constructs.AbstractArgs{
		Name:      a.Name(),
		Exported:  a.Exported(),
		Signature: i.Signature(a.Signature()),
	})
}

func (i *instantiator) Argument(a constructs.Argument) constructs.Argument {
	return i.proj.NewArgument(constructs.ArgumentArgs{
		Name: a.Name(),
		Type: i.TypeDesc(a.Type()),
	})
}

func (i *instantiator) Field(f constructs.Field) constructs.Field {
	return i.proj.NewField(constructs.FieldArgs{
		Name:     f.Name(),
		Exported: f.Exported(),
		Type:     i.TypeDesc(f.Type()),
		Embedded: f.Embedded(),
	})
}

func (i *instantiator) InterfaceInst(in constructs.InterfaceInst) constructs.TypeDesc {
	decl := in.Generic()
	return i.typeDecl(decl, decl.TypeParams(), in.InstanceTypes())
}

func (i *instantiator) ObjectInst(in constructs.ObjectInst) constructs.TypeDesc {
	decl := in.Generic()
	return i.typeDecl(decl, decl.TypeParams(), in.InstanceTypes())
}

func (i *instantiator) InterfaceDecl(decl constructs.InterfaceDecl) constructs.TypeDesc {
	tps := make([]constructs.TypeDesc, len(decl.TypeParams()))
	for i, tp := range decl.TypeParams() {
		tps[i] = tp
	}
	return i.typeDecl(decl, decl.TypeParams(), tps)
}

func (i *instantiator) Object(decl constructs.Object) constructs.TypeDesc {
	tps := make([]constructs.TypeDesc, len(decl.TypeParams()))
	for i, tp := range decl.TypeParams() {
		tps[i] = tp
	}
	return i.typeDecl(decl, decl.TypeParams(), tps)
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

func (i *instantiator) inProgress(decl constructs.TypeDecl, its []constructs.TypeDesc) bool {
	for !utils.IsNil(i) {
		if i.decl == decl && slices.Equal(i.instanceTypes, its) {
			return true
		}
		i = i.prior
	}
	return false
}

func (i *instantiator) typeDecl(decl constructs.TypeDecl, tps []constructs.TypeParam, its []constructs.TypeDesc) constructs.TypeDesc {
	its, anyReplaced := i.getInstanceTypeChange(its)
	if !anyReplaced {
		return decl
	}

	// If the declaration is the same as a declaration being instantiated,
	// create a reference to avoid the cycle. Cycles are caused by the type
	// being instantiated being used as part of the type definition
	// directly, e.g. `type Foo[T any] interface { Get() Foo[T]  }` or
	// indirectly. e.g. `type Foo[T any] interface { Children() List[Foo[T]] }`
	if i.inProgress(decl, its) {
		typ, found := i.proj.FindType(decl.Package().Path(), decl.Name(), its, true, false)
		if found {
			return typ
		}
		return i.proj.NewTempReference(constructs.TempReferenceArgs{
			PackagePath:   decl.Package().Path(),
			Name:          decl.Name(),
			InstanceTypes: its,
			Package:       decl.Package().Source(),
		})
	}

	i2, existing, _ := newInstantiator(i.log, i.proj, decl, tps, its)
	if !utils.IsNil(existing) {
		return existing.(constructs.TypeDesc)
	}
	i2.prior = i
	return i2.createInstance(nil).(constructs.TypeDesc)
}

func (i *instantiator) createInstance(realType types.Type) constructs.Construct {
	switch i.decl.Kind() {
	case kind.InterfaceDecl:
		d := i.decl.(constructs.InterfaceDecl)
		inst := i.proj.NewInterfaceInst(constructs.InterfaceInstArgs{
			RealType:      realType,
			Generic:       d,
			Resolved:      i.InterfaceDesc(d.Interface()),
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
	return i.proj.NewInterfaceDesc(constructs.InterfaceDescArgs{
		Hint:      it.Hint(),
		Abstracts: applyToSlice(it.Abstracts(), i.Abstract),
		Exact:     applyToSlice(it.Exact(), i.TypeDesc),
		Approx:    applyToSlice(it.Approx(), i.TypeDesc),
		Package:   i.decl.Package().Source(),
	})
}

func (i *instantiator) TempReference(r constructs.TempReference) constructs.TypeDesc {
	if r.Resolved() {
		// The reference will probably not be resolved, but just in case
		// it has been, just return the instantiated resolved type and
		// skip the instantiated reference.
		return i.TypeDesc(r.ResolvedType())
	}

	instTp := applyToSlice(r.InstanceTypes(), i.TypeDesc)
	typ, found := i.proj.FindType(r.PackagePath(), r.Name(), instTp, true, false)
	if found {
		return typ
	}
	return i.proj.NewTempReference(constructs.TempReferenceArgs{
		PackagePath:   r.PackagePath(),
		Name:          r.Name(),
		InstanceTypes: instTp,
		Package:       i.decl.Package().Source(),
	})
}

func (i *instantiator) Signature(s constructs.Signature) constructs.Signature {
	return i.proj.NewSignature(constructs.SignatureArgs{
		Variadic: s.Variadic(),
		Params:   applyToSlice(s.Params(), i.Argument),
		Results:  applyToSlice(s.Results(), i.Argument),
		Package:  i.decl.Package().Source(),
	})
}

func (i *instantiator) StructDesc(s constructs.StructDesc) constructs.StructDesc {
	return i.proj.NewStructDesc(constructs.StructDescArgs{
		Fields:  applyToSlice(s.Fields(), i.Field),
		Package: i.decl.Package().Source(),
	})
}

func (i *instantiator) TypeParam(tp constructs.TypeParam) constructs.TypeDesc {
	if t, has := i.conversion[tp]; has {
		return t
	}
	return tp
}
