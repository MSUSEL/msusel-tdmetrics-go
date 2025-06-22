package substituter

import (
	"fmt"
	"maps"

	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/assert"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs/kind"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/logger"
	"github.com/Snow-Gremlin/goToolbox/collections/sortedSet"
	"github.com/Snow-Gremlin/goToolbox/terrors/terror"
	"github.com/Snow-Gremlin/goToolbox/utils"
)

type Substituter interface {
	// Substitute recursively replaces constructs contained in the given
	// construct to create new constructs with the replacements.
	// The new constructs are added to the project.
	// Returns true if any substitution was performed.
	Substitute(orig constructs.Construct) (constructs.Construct, bool)
}

// New creates a new Substituter.
//
// The given replacements map is keyed by the construct to replace with
// the value of the construct to replace the key with. This is NOT recursive
// such that if a value is a key, it will NOT get that key's value.
//
// Some constructs may not able to be substituted in some constructs,
// such as packages, generics, receivers, and nests, because the substitution
// only heads down the constructs hierarchy.
//
// The given curPkg is the current package that is used when a package
// is needed but can't be determined from the construct itself.
func New(log *logger.Logger, proj constructs.Project, curPkg constructs.Package, replacements map[constructs.Construct]constructs.Construct) Substituter {
	assert.ArgNotNil(`project`, proj)
	assert.ArgNotNil(`current package`, curPkg)

	return &substituterImp{
		log:          log,
		proj:         proj,
		curPkg:       curPkg,
		replacements: maps.Clone(replacements),
		inProgress:   map[constructs.Construct]bool{},
		references:   map[constructs.Construct]constructs.TempDeclRef{},
	}
}

func Substitute[T constructs.Construct](s Substituter, orig T) (T, bool) {
	subbed, changed := s.Substitute(orig)
	return subbed.(T), changed
}

type substituterImp struct {
	log          *logger.Logger
	proj         constructs.Project
	curPkg       constructs.Package
	replacements map[constructs.Construct]constructs.Construct
	inProgress   map[constructs.Construct]bool
	references   map[constructs.Construct]constructs.TempDeclRef
}

func (s *substituterImp) Substitute(orig constructs.Construct) (constructs.Construct, bool) {
	if len(s.replacements) <= 0 {
		return orig, false
	}
	changed := false
	subbed := s.subConstruct(orig, &changed)
	return subbed, changed
}

func (s *substituterImp) subConstruct(con constructs.Construct, changed *bool) (subbed constructs.Construct) {
	if utils.IsNil(con) {
		return nil
	}

	// Check if there was a replacement, otherwise record the subbed result to
	// shortcut this substitution next time it occurs.
	if r, exists := s.replacements[con]; exists {
		*changed = *changed || constructs.ComparerPend(con, r)() == 0
		return r
	}
	defer func() {
		s.replacements[con] = subbed
	}()

	// Check if the current substitution is already in progress, if it is
	// create a reference that will be resolved after the in progress
	// construct is finished being substituted.
	if s.inProgress[con] {
		if ref, exists := s.references[con]; exists {
			*changed = true
			return ref
		}
		return s.createTempDeclRef(con, changed)
	}
	s.inProgress[con] = true
	defer func() {
		delete(s.inProgress, con)
		if ref, exists := s.references[con]; exists {
			ref.SetResolution(subbed.(constructs.TypeDesc))
			delete(s.references, con)
		}
	}()

	switch con.Kind() {
	case kind.Abstract:
		return s.subAbstract(con.(constructs.Abstract), changed)
	case kind.Argument:
		return s.subArgument(con.(constructs.Argument), changed)
	case kind.Basic:
		return s.subBasic(con.(constructs.Basic), changed)
	case kind.Field:
		return s.subField(con.(constructs.Field), changed)
	case kind.InterfaceDecl:
		return s.subInterfaceDecl(con.(constructs.InterfaceDecl), changed)
	case kind.InterfaceDesc:
		return s.subInterfaceDesc(con.(constructs.InterfaceDesc), changed)
	case kind.InterfaceInst:
		return s.subInterfaceInst(con.(constructs.InterfaceInst), changed)
	case kind.Method:
		return s.subMethod(con.(constructs.Method), changed)
	case kind.MethodInst:
		return s.subMethodInst(con.(constructs.MethodInst), changed)
	case kind.Metrics:
		return s.subMetrics(con.(constructs.Metrics), changed)
	case kind.Object:
		return s.subObject(con.(constructs.Object), changed)
	case kind.ObjectInst:
		return s.subObjectInst(con.(constructs.ObjectInst), changed)
	case kind.Package:
		return s.subPackage(con.(constructs.Package), changed)
	case kind.Selection:
		return s.subSelection(con.(constructs.Selection), changed)
	case kind.Signature:
		return s.subSignature(con.(constructs.Signature), changed)
	case kind.StructDesc:
		return s.subStructDesc(con.(constructs.StructDesc), changed)
	case kind.TempDeclRef:
		return s.subTempDeclRef(con.(constructs.TempDeclRef), changed)
	case kind.TempReference:
		return s.subTempReference(con.(constructs.TempReference), changed)
	case kind.TempTypeParamRef:
		return s.subTempTypeParamRef(con.(constructs.TempTypeParamRef), changed)
	case kind.TypeParam:
		return s.subTypeParam(con.(constructs.TypeParam), changed)
	case kind.Value:
		return s.subValue(con.(constructs.Value), changed)
	default:
		panic(terror.New(`unexpected construct kind in substituter`).
			With(`king`, con.Kind()).
			With(`con`, con))
	}
}

func subCon[T constructs.Construct](s *substituterImp, con T, changed *bool) T {
	return s.subConstruct(con, changed).(T)
}

func subConList[T constructs.Construct, S []T](s *substituterImp, list S, changed *bool) S {
	listSubbed := make(S, len(list))
	for i, con := range list {
		listSubbed[i] = subCon(s, con, changed)
	}
	return listSubbed
}

func (s *substituterImp) createTempDeclRef(con constructs.Construct, changed *bool) constructs.TempDeclRef {
	switch con.Kind() {
	case kind.InterfaceDecl:
		return s.createTempDeclRefForInterfaceDecl(con.(constructs.InterfaceDecl), changed)
	case kind.InterfaceInst:
		return s.createTempDeclRefForInterfaceInst(con.(constructs.InterfaceInst), changed)
	case kind.Object:
		return s.createTempDeclRefForObject(con.(constructs.Object), changed)
	case kind.ObjectInst:
		return s.createTempDeclRefForObjectInst(con.(constructs.ObjectInst), changed)
	case kind.Method:
		return s.createTempDeclRefForMethod(con.(constructs.Method), changed)
	case kind.MethodInst:
		return s.createTempDeclRefForMethodInst(con.(constructs.MethodInst), changed)
	case kind.Value:
		return s.createTempDeclRefForValue(con.(constructs.Value), changed)
	default:
		panic(terror.New(`unexpected construct kind of referencing in substituter`).
			With(`king`, con.Kind()).
			With(`con`, con))
	}
}

func (s *substituterImp) createTempDeclRefForInterfaceDecl(con constructs.InterfaceDecl, changed *bool) constructs.TempDeclRef {
	return s.proj.NewTempDeclRef(constructs.TempDeclRefArgs{
		PackagePath: con.Package().Path(),
		Name:        con.Name(),
		Nest:        subCon(s, con.Nest(), changed),
	})
}

func (s *substituterImp) createTempDeclRefForInterfaceInst(con constructs.InterfaceInst, changed *bool) constructs.TempDeclRef {
	return s.proj.NewTempDeclRef(constructs.TempDeclRefArgs{
		PackagePath:   con.Generic().Package().Path(),
		Name:          con.Generic().Name(),
		ImplicitTypes: subConList(s, con.ImplicitTypes(), changed),
		InstanceTypes: subConList(s, con.InstanceTypes(), changed),
		Nest:          subCon(s, con.Generic().Nest(), changed),
	})
}

func (s *substituterImp) createTempDeclRefForObject(con constructs.Object, changed *bool) constructs.TempDeclRef {
	return s.proj.NewTempDeclRef(constructs.TempDeclRefArgs{
		PackagePath: con.Package().Path(),
		Name:        con.Name(),
		Nest:        subCon(s, con.Nest(), changed),
	})
}

func (s *substituterImp) createTempDeclRefForObjectInst(con constructs.ObjectInst, changed *bool) constructs.TempDeclRef {
	return s.proj.NewTempDeclRef(constructs.TempDeclRefArgs{
		PackagePath:   con.Generic().Package().Path(),
		Name:          con.Generic().Name(),
		ImplicitTypes: subConList(s, con.ImplicitTypes(), changed),
		InstanceTypes: subConList(s, con.InstanceTypes(), changed),
		Nest:          subCon(s, con.Generic().Nest(), changed),
	})
}

func (s *substituterImp) createTempDeclRefForMethod(con constructs.Method, changed *bool) constructs.TempDeclRef {
	_ = changed
	return s.proj.NewTempDeclRef(constructs.TempDeclRefArgs{
		PackagePath: con.Package().Path(),
		Name:        con.Name(),
	})
}

func (s *substituterImp) createTempDeclRefForMethodInst(con constructs.MethodInst, changed *bool) constructs.TempDeclRef {
	return s.proj.NewTempDeclRef(constructs.TempDeclRefArgs{
		PackagePath:   con.Generic().Package().Path(),
		Name:          con.Generic().Name(),
		InstanceTypes: subConList(s, con.InstanceTypes(), changed),
	})
}

func (s *substituterImp) createTempDeclRefForValue(con constructs.Value, changed *bool) constructs.TempDeclRef {
	_ = changed
	return s.proj.NewTempDeclRef(constructs.TempDeclRefArgs{
		PackagePath: con.Package().Path(),
		Name:        con.Name(),
	})
}

func finishSubCon[TCon constructs.Construct, TArg any](subChanged bool, orig TCon, factory func(TArg) TCon, args TArg, changed *bool) TCon {
	if subChanged {
		fmt.Printf(">> Original: %+v\n", orig) // TODO: REMOVE
		fmt.Printf("   Replaced: %+v\n", args) // TODO: REMOVE
		*changed = true
		return factory(args)
	}
	return orig
}

func (s *substituterImp) subAbstract(con constructs.Abstract, changed *bool) constructs.Abstract {
	subChanged := false
	args := constructs.AbstractArgs{
		Name:      con.Name(),
		Exported:  con.Exported(),
		Signature: subCon(s, con.Signature(), &subChanged),
	}
	return finishSubCon(subChanged, con, s.proj.NewAbstract, args, changed)
}

func (s *substituterImp) subArgument(con constructs.Argument, changed *bool) constructs.Argument {
	subChanged := false
	args := constructs.ArgumentArgs{
		Name: con.Name(),
		Type: subCon(s, con.Type(), &subChanged),
	}
	return finishSubCon(subChanged, con, s.proj.NewArgument, args, changed)
}

func (s *substituterImp) subBasic(con constructs.Basic, changed *bool) constructs.Basic {
	_ = changed
	return con
}

func (s *substituterImp) subField(con constructs.Field, changed *bool) constructs.Field {
	subChanged := false
	args := constructs.FieldArgs{
		Name: con.Name(),
		Type: subCon(s, con.Type(), &subChanged),
	}
	return finishSubCon(subChanged, con, s.proj.NewField, args, changed)
}

func (s *substituterImp) subInterfaceDecl(con constructs.InterfaceDecl, changed *bool) constructs.InterfaceDecl {
	subChanged := false
	args := constructs.InterfaceDeclArgs{
		Package:    con.Package(),
		Name:       con.Name(),
		Exported:   con.Exported(),
		Location:   con.Location(),
		Nest:       con.Nest(),
		TypeParams: subConList(s, con.TypeParams(), &subChanged),
		Interface:  subCon(s, con.Interface(), &subChanged),
	}
	return finishSubCon(subChanged, con, s.proj.NewInterfaceDecl, args, changed)
}

func (s *substituterImp) subInterfaceDesc(con constructs.InterfaceDesc, changed *bool) constructs.InterfaceDesc {
	subChanged := false
	args := constructs.InterfaceDescArgs{
		Hint:      con.Hint(),
		PinnedPkg: con.PinnedPackage(),
		Abstracts: subConList(s, con.Abstracts(), &subChanged),
		Exact:     subConList(s, con.Exact(), &subChanged),
		Approx:    subConList(s, con.Approx(), &subChanged),
		Package:   con.PinnedPackage().Source(),
	}
	return finishSubCon(subChanged, con, s.proj.NewInterfaceDesc, args, changed)
}

func (s *substituterImp) subInterfaceInst(con constructs.InterfaceInst, changed *bool) constructs.InterfaceInst {
	subChanged := false
	args := constructs.InterfaceInstArgs{
		Generic:       con.Generic(),
		Resolved:      subCon(s, con.Resolved(), &subChanged),
		ImplicitTypes: subConList(s, con.ImplicitTypes(), &subChanged),
		InstanceTypes: subConList(s, con.InstanceTypes(), &subChanged),
	}
	return finishSubCon(subChanged, con, s.proj.NewInterfaceInst, args, changed)
}

func (s *substituterImp) subMethod(con constructs.Method, changed *bool) constructs.Method {
	subChanged := false
	args := constructs.MethodArgs{
		Package:     con.Package(),
		Name:        con.Name(),
		Exported:    con.Exported(),
		Location:    con.Location(),
		TypeParams:  subConList(s, con.TypeParams(), &subChanged),
		Signature:   subCon(s, con.Signature(), &subChanged),
		Metrics:     subCon(s, con.Metrics(), &subChanged),
		RecvName:    con.ReceiverName(),
		Receiver:    con.Receiver(),
		PointerRecv: con.PointerRecv(),
	}
	return finishSubCon(subChanged, con, s.proj.NewMethod, args, changed)
}

func (s *substituterImp) subMethodInst(con constructs.MethodInst, changed *bool) constructs.MethodInst {
	subChanged := false
	args := constructs.MethodInstArgs{
		Generic:       con.Generic(),
		Resolved:      subCon(s, con.Resolved(), &subChanged),
		InstanceTypes: subConList(s, con.InstanceTypes(), &subChanged),
		Metrics:       subCon(s, con.Metrics(), &subChanged),
	}
	return finishSubCon(subChanged, con, s.proj.NewMethodInst, args, changed)
}

func (s *substituterImp) subMetrics(con constructs.Metrics, changed *bool) constructs.Metrics {
	subChanged := false
	cmp := constructs.Comparer[constructs.Construct]()
	args := constructs.MetricsArgs{
		Location:   con.Location(),
		Complexity: con.Complexity(),
		LineCount:  con.LineCount(),
		CodeCount:  con.CodeCount(),
		Indents:    con.Indents(),
		Getter:     con.Getter(),
		Setter:     con.Setter(),
		SideEffect: con.SideEffect(),
		Reads:      sortedSet.With(subConList(s, con.Reads().ToSlice(), &subChanged), cmp),
		Writes:     sortedSet.With(subConList(s, con.Writes().ToSlice(), &subChanged), cmp),
		Invokes:    sortedSet.With(subConList(s, con.Invokes().ToSlice(), &subChanged), cmp),
	}
	return finishSubCon(subChanged, con, s.proj.NewMetrics, args, changed)
}

func (s *substituterImp) subObject(con constructs.Object, changed *bool) constructs.Object {
	subChanged := false
	args := constructs.ObjectArgs{
		Package:    con.Package(),
		Name:       con.Name(),
		Exported:   con.Exported(),
		Location:   con.Location(),
		Nest:       con.Nest(),
		TypeParams: subConList(s, con.TypeParams(), &subChanged),
		Data:       subCon(s, con.Data(), &subChanged),
	}
	// TODO: Add methods, interface, instances, etc?
	return finishSubCon(subChanged, con, s.proj.NewObject, args, changed)
}

func (s *substituterImp) subObjectInst(con constructs.ObjectInst, changed *bool) constructs.ObjectInst {
	subChanged := false
	args := constructs.ObjectInstArgs{
		Generic:       con.Generic(),
		ResolvedData:  subCon(s, con.ResolvedData(), &subChanged),
		ImplicitTypes: subConList(s, con.ImplicitTypes(), &subChanged),
		InstanceTypes: subConList(s, con.InstanceTypes(), &subChanged),
	}
	return finishSubCon(subChanged, con, s.proj.NewObjectInst, args, changed)
}

func (s *substituterImp) subPackage(con constructs.Package, changed *bool) constructs.Package {
	// The package real-types are required for a lot of the abstractor and
	// a new package real-type can not be created, so we can't substitute it.
	panic(terror.New(`may not substitute a package`))
}

func (s *substituterImp) subSelection(con constructs.Selection, changed *bool) constructs.Selection {
	subChanged := false
	args := constructs.SelectionArgs{
		Name:   con.Name(),
		Origin: subCon(s, con.Origin(), &subChanged),
	}
	return finishSubCon(subChanged, con, s.proj.NewSelection, args, changed)
}

func (s *substituterImp) subSignature(con constructs.Signature, changed *bool) constructs.Signature {
	subChanged := false
	args := constructs.SignatureArgs{
		Variadic: con.Variadic(),
		Params:   subConList(s, con.Params(), &subChanged),
		Results:  subConList(s, con.Results(), &subChanged),
		Package:  s.curPkg.Source(),
	}
	return finishSubCon(subChanged, con, s.proj.NewSignature, args, changed)
}

func (s *substituterImp) subStructDesc(con constructs.StructDesc, changed *bool) constructs.StructDesc {
	subChanged := false
	args := constructs.StructDescArgs{
		Fields:  subConList(s, con.Fields(), &subChanged),
		Package: s.curPkg.Source(),
	}
	return finishSubCon(subChanged, con, s.proj.NewStructDesc, args, changed)
}

func (s *substituterImp) subTempDeclRef(con constructs.TempDeclRef, changed *bool) constructs.Construct {
	if con.Resolved() {
		return subCon(s, con.ResolvedType(), changed)
	}
	subChanged := false
	args := constructs.TempDeclRefArgs{
		PackagePath:   con.PackagePath(),
		Name:          con.Name(),
		ImplicitTypes: subConList(s, con.ImplicitTypes(), &subChanged),
		InstanceTypes: subConList(s, con.InstanceTypes(), &subChanged),
		Nest:          con.Nest(),
	}
	return finishSubCon(subChanged, con, s.proj.NewTempDeclRef, args, changed)
}

func (s *substituterImp) subTempReference(con constructs.TempReference, changed *bool) constructs.Construct {
	if con.Resolved() {
		return subCon(s, con.ResolvedType(), changed)
	}
	subChanged := false
	args := constructs.TempReferenceArgs{
		PackagePath:   con.PackagePath(),
		Name:          con.Name(),
		ImplicitTypes: subConList(s, con.ImplicitTypes(), &subChanged),
		InstanceTypes: subConList(s, con.InstanceTypes(), &subChanged),
		Nest:          con.Nest(),
		Package:       s.curPkg.Source(),
	}
	return finishSubCon(subChanged, con, s.proj.NewTempReference, args, changed)
}

func (s *substituterImp) subTempTypeParamRef(con constructs.TempTypeParamRef, changed *bool) constructs.Construct {
	if con.Resolved() {
		return subCon(s, con.ResolvedType(), changed)
	}
	return con
}

func (s *substituterImp) subTypeParam(con constructs.TypeParam, changed *bool) constructs.TypeParam {
	subChanged := false
	args := constructs.TypeParamArgs{
		Name: con.Name(),
		Type: subCon(s, con.Type(), &subChanged),
	}
	return finishSubCon(subChanged, con, s.proj.NewTypeParam, args, changed)
}

func (s *substituterImp) subValue(con constructs.Value, changed *bool) constructs.Value {
	subChanged := false
	args := constructs.ValueArgs{
		Package:  con.Package(),
		Name:     con.Name(),
		Exported: con.Exported(),
		Location: con.Location(),
		Type:     subCon(s, con.Type(), &subChanged),
		Const:    con.Const(),
		Metrics:  subCon(s, con.Metrics(), &subChanged),
	}
	return finishSubCon(subChanged, con, s.proj.NewValue, args, changed)
}
