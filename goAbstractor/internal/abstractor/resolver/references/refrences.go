package references

import (
	"github.com/Snow-Gremlin/goToolbox/terrors/terror"
	"github.com/Snow-Gremlin/goToolbox/utils"

	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/abstractor/instantiator"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs/kind"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/logger"
)

func References(log *logger.Logger, proj constructs.Project, required bool) bool {
	changed := removeTempReferences(proj, required)
	changed = removeTempDeclRefs(proj, required) || changed
	changed = tempReferences(log, proj) || changed
	changed = tempDeclRefs(log, proj) || changed
	if required {
		clearAllTemps(proj)
	}
	return changed
}

func clearAllTemps(proj constructs.Project) {
	proj.ClearAllTempReferences()
	proj.ClearAllTempTypeParamRefs()
	proj.ClearAllTempDeclRefs()
}

func removeTempReferences(proj constructs.Project, required bool) bool {
	changed := false
	proj.AllConstructs().Foreach(func(c constructs.Construct) {
		if trc, has := c.(constructs.TempReferenceContainer); has {
			changed = trc.RemoveTempReferences(required) || changed
		}
	})
	return changed
}

func removeTempDeclRefs(proj constructs.Project, required bool) bool {
	changed := false
	proj.AllConstructs().Foreach(func(c constructs.Construct) {
		if trc, has := c.(constructs.TempDeclRefContainer); has {
			changed = trc.RemoveTempDeclRefs(required) || changed
		}
	})
	return changed
}

func tempReferences(log *logger.Logger, proj constructs.Project) bool {
	changed := false
	refs := proj.TempReferences()
	for i := range refs.Count() {
		changed = resolveTempRef(log, proj, refs.Get(i)) || changed
	}
	return removeTempReferences(proj, true) || changed
}

func resolveTempRef(log *logger.Logger, proj constructs.Project, ref constructs.TempReference) bool {
	if ref.Resolved() {
		return false
	}

	// Try to find instance of type or non-generic type.
	typ, ok := proj.FindType(ref.PackagePath(), ref.Name(), ref.Nest(), ref.ImplicitTypes(), ref.InstanceTypes(), false, false)
	if ok {
		ref.SetResolution(typ)
		return true
	}

	// Try to find generic type and then create the instance if needed.
	typ, ok = proj.FindType(ref.PackagePath(), ref.Name(), ref.Nest(), nil, nil, false, true)
	if !ok {
		panic(terror.New(`failed to find temp referenced object`).
			With(`package path`, ref.PackagePath()).
			With(`name`, ref.Name()).
			With(`instance types`, ref.InstanceTypes()))
	}
	if len(ref.InstanceTypes()) <= 0 && len(ref.ImplicitTypes()) <= 0 {
		ref.SetResolution(typ)
		return true
	}

	switch typ.Kind() {
	case kind.Object:
		res := instantiator.Object(log, proj, ref.GoType(), typ.(constructs.Object), ref.ImplicitTypes(), ref.InstanceTypes())
		if utils.IsNil(res) {
			panic(terror.New(`failed to resolve object type reference`).
				With(`package path`, ref.PackagePath()).
				With(`name`, ref.Name()).
				With(`instance types`, ref.InstanceTypes()))
		}
		ref.SetResolution(res)
		return true

	case kind.InterfaceDecl:
		res := instantiator.InterfaceDecl(log, proj, ref.GoType(), typ.(constructs.InterfaceDecl), ref.ImplicitTypes(), ref.InstanceTypes())
		if utils.IsNil(res) {
			panic(terror.New(`failed to resolve interface type reference`).
				With(`package path`, ref.PackagePath()).
				With(`name`, ref.Name()).
				With(`instance types`, ref.InstanceTypes()))
		}
		ref.SetResolution(res)
		return true

	default:
		panic(terror.New(`unexpected declaration type`).
			With(`kind`, typ.Kind()).
			With(`decl`, typ))
	}
}

func tempDeclRefs(log *logger.Logger, proj constructs.Project) bool {
	changed := false
	refs := proj.TempDeclRefs()
	for i := range refs.Count() {
		changed = resolveTempDeclRef(log, proj, refs.Get(i)) || changed
	}
	return removeTempDeclRefs(proj, true) || changed
}

func resolveTempDeclRef(log *logger.Logger, proj constructs.Project, ref constructs.TempDeclRef) bool {
	if ref.Resolved() {
		return false
	}

	// Try to find instance of declaration or non-generic declaration.
	decl, ok := proj.FindDecl(ref.PackagePath(), ref.Name(), ref.Nest(), ref.ImplicitTypes(), ref.InstanceTypes(), false, false)
	if ok {
		ref.SetResolution(decl)
		return true
	}

	// Try to find generic declaration and then create the instance if needed.
	decl, ok = proj.FindDecl(ref.PackagePath(), ref.Name(), ref.Nest(), nil, nil, false, true)
	if !ok {
		panic(terror.New(`failed to find temp declaration referenced`).
			With(`package path`, ref.PackagePath()).
			With(`name`, ref.Name()).
			With(`instance types`, ref.InstanceTypes()))
	}
	if len(ref.InstanceTypes()) <= 0 && len(ref.ImplicitTypes()) <= 0 {
		ref.SetResolution(decl)
		return true
	}

	switch decl.Kind() {
	case kind.Object:
		res := instantiator.Object(log, proj, nil, decl.(constructs.Object), nil, ref.InstanceTypes())
		if utils.IsNil(res) {
			panic(terror.New(`failed to resolve object declaration reference`).
				With(`package path`, ref.PackagePath()).
				With(`name`, ref.Name()).
				With(`instance types`, ref.InstanceTypes()))
		}
		ref.SetResolution(res)
		return true

	case kind.InterfaceDecl:
		res := instantiator.InterfaceDecl(log, proj, nil, decl.(constructs.InterfaceDecl), nil, ref.InstanceTypes())
		if utils.IsNil(res) {
			panic(terror.New(`failed to resolve interface declaration reference`).
				With(`package path`, ref.PackagePath()).
				With(`name`, ref.Name()).
				With(`instance types`, ref.InstanceTypes()))
		}
		ref.SetResolution(res)
		return true

	case kind.Method:
		res := instantiator.Method(log, proj, decl.(constructs.Method), ref.InstanceTypes())
		if utils.IsNil(res) {
			panic(terror.New(`failed to resolve method declaration reference`).
				With(`package path`, ref.PackagePath()).
				With(`name`, ref.Name()).
				With(`instance types`, ref.InstanceTypes()))
		}
		ref.SetResolution(res)
		return true

	default:
		panic(terror.New(`unexpected declaration type`).
			With(`kind`, decl.Kind()).
			With(`decl`, decl))
	}
}
