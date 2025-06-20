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
	// Substitute recursively replaces constructs contained in the given construct
	// to create new constructs with the replacements. The new constructs are
	// added to the project.
	Substitute(orig constructs.Construct)
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

func (s *substituterImp) Substitute(orig constructs.Construct) {
	if len(s.replacements) > 0 {
		s.subConstruct(orig)
	}
}

func (s *substituterImp) subConstruct(con constructs.Construct) (subbed constructs.Construct, changed bool) {
	if utils.IsNil(con) {
		return nil, false
	}

	// Check if there was a replacement, otherwise record the subbed result to
	// shortcut this substitution next time it occurs.
	if r, exists := s.replacements[con]; exists {
		return r, constructs.ComparerPend(con, r)() == 0
	}
	defer func() {
		s.replacements[con] = subbed
	}()

	// Check if the current substitution is already in progress, if it is
	// create a reference that will be resolved after the in progress
	// construct is finished being substituted.
	if s.inProgress[con] {
		if ref, exists := s.references[con]; exists {
			return ref, true
		}
		return s.createTempDeclRef(con), true
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
		return s.subAbstract(con.(constructs.Abstract))
	case kind.Argument:
		return s.subArgument(con.(constructs.Argument))
	case kind.Basic:
		return s.subBasic(con.(constructs.Basic))
	case kind.Field:
		return s.subField(con.(constructs.Field))
	case kind.InterfaceDecl:
		return s.subInterfaceDecl(con.(constructs.InterfaceDecl))
	case kind.InterfaceDesc:
		return s.subInterfaceDesc(con.(constructs.InterfaceDesc))
	case kind.InterfaceInst:
		return s.subInterfaceInst(con.(constructs.InterfaceInst))
	case kind.Method:
		return s.subMethod(con.(constructs.Method))
	case kind.MethodInst:
		return s.subMethodInst(con.(constructs.MethodInst))
	case kind.Metrics:
		return s.subMetrics(con.(constructs.Metrics))
	case kind.Object:
		return s.subObject(con.(constructs.Object))
	case kind.ObjectInst:
		return s.subObjectInst(con.(constructs.ObjectInst))
	case kind.Package:
		return s.subPackage(con.(constructs.Package))
	case kind.Selection:
		return s.subSelection(con.(constructs.Selection))
	case kind.Signature:
		return s.subSignature(con.(constructs.Signature))
	case kind.StructDesc:
		return s.subStructDesc(con.(constructs.StructDesc))
	case kind.TempDeclRef:
		return s.subTempDeclRef(con.(constructs.TempDeclRef))
	case kind.TempReference:
		return s.subTempReference(con.(constructs.TempReference))
	case kind.TempTypeParamRef:
		return s.subTempTypeParamRef(con.(constructs.TempTypeParamRef))
	case kind.TypeParam:
		return s.subTypeParam(con.(constructs.TypeParam))
	case kind.Value:
		return s.subValue(con.(constructs.Value))
	default:
		panic(terror.New(`unexpected construct kind in substituter`).
			With(`king`, con.Kind()).
			With(`con`, con))
	}
}

func subConstruct[T constructs.Construct](s *substituterImp, con T) (subbed T, changed bool) {
	conSubbed, changed := s.subConstruct(con)
	return conSubbed.(T), changed
}

func (s *substituterImp) createTempDeclRef(con constructs.Construct) constructs.TempDeclRef {
	switch con.Kind() {
	case kind.InterfaceDecl:
		return s.subAbstractForInterfaceDecl(con.(constructs.InterfaceDecl))
	case kind.Object:
		return s.subAbstractForObject(con.(constructs.Object))
	case kind.Method:
		return s.subAbstractForMethod(con.(constructs.Method))
	case kind.Value:
		return s.subAbstractForValue(con.(constructs.Value))
	default:
		panic(terror.New(`unexpected construct kind of referencing in substituter`).
			With(`king`, con.Kind()).
			With(`con`, con))
	}
}

func (s *substituterImp) subAbstractForInterfaceDecl(con constructs.InterfaceDecl) constructs.TempDeclRef {
	panic(terror.New(`unimplemented`)) // TODO: Implement
}

func (s *substituterImp) subAbstractForObject(con constructs.Object) constructs.TempDeclRef {
	panic(terror.New(`unimplemented`)) // TODO: Implement
}

func (s *substituterImp) subAbstractForMethod(con constructs.Method) constructs.TempDeclRef {
	panic(terror.New(`unimplemented`)) // TODO: Implement
}

func (s *substituterImp) subAbstractForValue(con constructs.Value) constructs.TempDeclRef {
	panic(terror.New(`unimplemented`)) // TODO: Implement
}

func (s *substituterImp) subAbstract(con constructs.Abstract) (constructs.Abstract, bool) {
	panic(terror.New(`unimplemented`)) // TODO: Implement
}

func (s *substituterImp) subArgument(con constructs.Argument) (constructs.Argument, bool) {
	panic(terror.New(`unimplemented`)) // TODO: Implement
}

func (s *substituterImp) subBasic(con constructs.Basic) (constructs.Basic, bool) {
	panic(terror.New(`unimplemented`)) // TODO: Implement
}

func (s *substituterImp) subField(con constructs.Field) (constructs.Field, bool) {
	panic(terror.New(`unimplemented`)) // TODO: Implement
}

func (s *substituterImp) subInterfaceDecl(con constructs.InterfaceDecl) (constructs.InterfaceDecl, bool) {
	panic(terror.New(`unimplemented`)) // TODO: Implement
}

func (s *substituterImp) subInterfaceDesc(con constructs.InterfaceDesc) (constructs.InterfaceDesc, bool) {
	panic(terror.New(`unimplemented`)) // TODO: Implement
}

func (s *substituterImp) subInterfaceInst(con constructs.InterfaceInst) (constructs.InterfaceInst, bool) {
	panic(terror.New(`unimplemented`)) // TODO: Implement
}

func (s *substituterImp) subMethod(con constructs.Method) (constructs.Method, bool) {
	panic(terror.New(`unimplemented`)) // TODO: Implement
}

func (s *substituterImp) subMethodInst(con constructs.MethodInst) (constructs.MethodInst, bool) {
	panic(terror.New(`unimplemented`)) // TODO: Implement
}

func (s *substituterImp) subMetrics(con constructs.Metrics) (constructs.Metrics, bool) {
	panic(terror.New(`unimplemented`)) // TODO: Implement
}

func (s *substituterImp) subObject(con constructs.Object) (constructs.Object, bool) {
	panic(terror.New(`unimplemented`)) // TODO: Implement
}

func (s *substituterImp) subObjectInst(con constructs.ObjectInst) (constructs.ObjectInst, bool) {
	panic(terror.New(`unimplemented`)) // TODO: Implement
}

func (s *substituterImp) subPackage(con constructs.Package) (constructs.Package, bool) {
	panic(terror.New(`unimplemented`)) // TODO: Implement
}

func (s *substituterImp) subSelection(con constructs.Selection) (constructs.Selection, bool) {
	panic(terror.New(`unimplemented`)) // TODO: Implement
}

func (s *substituterImp) subSignature(con constructs.Signature) (constructs.Signature, bool) {
	panic(terror.New(`unimplemented`)) // TODO: Implement
}

func (s *substituterImp) subStructDesc(con constructs.StructDesc) (constructs.StructDesc, bool) {
	panic(terror.New(`unimplemented`)) // TODO: Implement
}

func (s *substituterImp) subTempDeclRef(con constructs.TempDeclRef) (constructs.TempDeclRef, bool) {
	panic(terror.New(`unimplemented`)) // TODO: Implement
}

func (s *substituterImp) subTempReference(con constructs.TempReference) (constructs.TempReference, bool) {
	panic(terror.New(`unimplemented`)) // TODO: Implement
}

func (s *substituterImp) subTempTypeParamRef(con constructs.TempTypeParamRef) (constructs.TempTypeParamRef, bool) {
	panic(terror.New(`unimplemented`)) // TODO: Implement
}

func (s *substituterImp) subTypeParam(con constructs.TypeParam) (constructs.TypeParam, bool) {
	panic(terror.New(`unimplemented`)) // TODO: Implement
}

func (s *substituterImp) subValue(con constructs.Value) (constructs.Value, bool) {
	panic(terror.New(`unimplemented`)) // TODO: Implement
}
