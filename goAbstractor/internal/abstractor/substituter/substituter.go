package substituter

import (
	"maps"

	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/assert"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs/kind"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/logger"
	"github.com/Snow-Gremlin/goToolbox/terrors/terror"
	"github.com/Snow-Gremlin/goToolbox/utils"
)

type Substituter interface {
	// Substitute recursively replaces constructs contained in the given
	// construct to create new constructs with the replacements.
	// The new constructs are added to the project.
	// Returns true if any substitution was performed.
	Substitute(orig constructs.Construct) bool
}

func New(log logger.Logger, proj constructs.Project, replacements map[constructs.Construct]constructs.Construct) Substituter {
	assert.ArgNotNil(`project`, proj)

	return &substituterImp{
		log:          log,
		proj:         proj,
		replacements: maps.Clone(replacements),
		inProgress:   map[constructs.Construct]bool{},
		references:   map[constructs.Construct]constructs.TempDeclRef{},
	}
}

type substituterImp struct {
	log          logger.Logger
	proj         constructs.Project
	replacements map[constructs.Construct]constructs.Construct
	inProgress   map[constructs.Construct]bool
	references   map[constructs.Construct]constructs.TempDeclRef
}

func (s *substituterImp) Substitute(orig constructs.Construct) bool {
	if len(s.replacements) <= 0 {
		return false
	}
	changed := false
	s.subConstruct(orig, &changed)
	return changed
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

func castList[TOut, TIn any, SIn ~[]TIn](list SIn) []TOut {
	casted := make([]TOut, len(list))
	for i, v := range list {
		casted[i] = any(v).(TOut)
	}
	return casted
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
	return s.proj.NewTempDeclRef(constructs.TempDeclRefArgs{
		PackagePath: con.Package().Path(),
		Name:        con.Name(),
	})
}

func (s *substituterImp) subAbstract(con constructs.Abstract, changed *bool) constructs.Abstract {
	panic(terror.New(`unimplemented`)) // TODO: Implement
}

func (s *substituterImp) subArgument(con constructs.Argument, changed *bool) constructs.Argument {
	panic(terror.New(`unimplemented`)) // TODO: Implement
}

func (s *substituterImp) subBasic(con constructs.Basic, changed *bool) constructs.Basic {
	panic(terror.New(`unimplemented`)) // TODO: Implement
}

func (s *substituterImp) subField(con constructs.Field, changed *bool) constructs.Field {
	panic(terror.New(`unimplemented`)) // TODO: Implement
}

func (s *substituterImp) subInterfaceDecl(con constructs.InterfaceDecl, changed *bool) constructs.InterfaceDecl {
	panic(terror.New(`unimplemented`)) // TODO: Implement
}

func (s *substituterImp) subInterfaceDesc(con constructs.InterfaceDesc, changed *bool) constructs.InterfaceDesc {
	panic(terror.New(`unimplemented`)) // TODO: Implement
}

func (s *substituterImp) subInterfaceInst(con constructs.InterfaceInst, changed *bool) constructs.InterfaceInst {
	panic(terror.New(`unimplemented`)) // TODO: Implement
}

func (s *substituterImp) subMethod(con constructs.Method, changed *bool) constructs.Method {
	panic(terror.New(`unimplemented`)) // TODO: Implement
}

func (s *substituterImp) subMethodInst(con constructs.MethodInst, changed *bool) constructs.MethodInst {
	panic(terror.New(`unimplemented`)) // TODO: Implement
}

func (s *substituterImp) subMetrics(con constructs.Metrics, changed *bool) constructs.Metrics {
	panic(terror.New(`unimplemented`)) // TODO: Implement
}

func (s *substituterImp) subObject(con constructs.Object, changed *bool) constructs.Object {
	panic(terror.New(`unimplemented`)) // TODO: Implement
}

func (s *substituterImp) subObjectInst(con constructs.ObjectInst, changed *bool) constructs.ObjectInst {
	panic(terror.New(`unimplemented`)) // TODO: Implement
}

func (s *substituterImp) subPackage(con constructs.Package, changed *bool) constructs.Package {
	panic(terror.New(`unimplemented`)) // TODO: Implement
}

func (s *substituterImp) subSelection(con constructs.Selection, changed *bool) constructs.Selection {
	panic(terror.New(`unimplemented`)) // TODO: Implement
}

func (s *substituterImp) subSignature(con constructs.Signature, changed *bool) constructs.Signature {
	panic(terror.New(`unimplemented`)) // TODO: Implement
}

func (s *substituterImp) subStructDesc(con constructs.StructDesc, changed *bool) constructs.StructDesc {
	panic(terror.New(`unimplemented`)) // TODO: Implement
}

func (s *substituterImp) subTempDeclRef(con constructs.TempDeclRef, changed *bool) constructs.TempDeclRef {
	panic(terror.New(`unimplemented`)) // TODO: Implement
}

func (s *substituterImp) subTempReference(con constructs.TempReference, changed *bool) constructs.TempReference {
	panic(terror.New(`unimplemented`)) // TODO: Implement
}

func (s *substituterImp) subTempTypeParamRef(con constructs.TempTypeParamRef, changed *bool) constructs.TempTypeParamRef {
	panic(terror.New(`unimplemented`)) // TODO: Implement
}

func (s *substituterImp) subTypeParam(con constructs.TypeParam, changed *bool) constructs.TypeParam {
	panic(terror.New(`unimplemented`)) // TODO: Implement
}

func (s *substituterImp) subValue(con constructs.Value, changed *bool) constructs.Value {
	panic(terror.New(`unimplemented`)) // TODO: Implement
}
