package references

import (
	"github.com/Snow-Gremlin/goToolbox/terrors/terror"
	"github.com/Snow-Gremlin/goToolbox/utils"

	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/abstractor/instantiator"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs/kind"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/logger"
)

func References(log *logger.Logger, proj constructs.Project) {
	tempReferences(log, proj)
	tempDeclRefs(log, proj)
}

func tempReferences(log *logger.Logger, proj constructs.Project) {
	refs := proj.TempReferences()
	for i := range refs.Count() {
		resolveTempRef(log, proj, refs.Get(i))
	}

	proj.AllConstructs().Foreach(func(c constructs.Construct) {
		if trc, has := c.(constructs.TempReferenceContainer); has {
			trc.RemoveTempReferences()
		}
	})
	proj.ClearAllTempReferences()
	proj.ClearAllTempTypeParamRefs()
}

func resolveTempRef(log *logger.Logger, proj constructs.Project, ref constructs.TempReference) {
	if ref.Resolved() {
		return
	}

	// Try to find instance of type or non-generic type.
	typ, ok := proj.FindType(ref.PackagePath(), ref.Name(), ref.InstanceTypes(), false, false)
	if ok {
		ref.SetResolution(typ)
		return
	}

	// Try to find generic type and then create the instance if needed.
	typ, ok = proj.FindType(ref.PackagePath(), ref.Name(), []constructs.TypeDesc{}, false, true)
	if !ok {
		panic(terror.New(`failed to find temp referenced object`).
			With(`package path`, ref.PackagePath()).
			With(`name`, ref.Name()).
			With(`instance types`, ref.InstanceTypes()))
	}
	if len(ref.InstanceTypes()) <= 0 {
		ref.SetResolution(typ)
		return
	}

	switch typ.Kind() {
	case kind.Object:
		res := instantiator.Object(log, proj, ref.GoType(), typ.(constructs.Object), nil, ref.InstanceTypes())
		if utils.IsNil(res) {
			panic(terror.New(`failed to resolve object type reference`).
				With(`package path`, ref.PackagePath()).
				With(`name`, ref.Name()).
				With(`instance types`, ref.InstanceTypes()))
		}
		ref.SetResolution(res)

	case kind.InterfaceDecl:
		res := instantiator.InterfaceDecl(log, proj, ref.GoType(), typ.(constructs.InterfaceDecl), nil, ref.InstanceTypes())
		if utils.IsNil(res) {
			panic(terror.New(`failed to resolve interface type reference`).
				With(`package path`, ref.PackagePath()).
				With(`name`, ref.Name()).
				With(`instance types`, ref.InstanceTypes()))
		}
		ref.SetResolution(res)

	default:
		panic(terror.New(`unexpected declaration type`).
			With(`kind`, typ.Kind()).
			With(`decl`, typ))
	}
}

func tempDeclRefs(log *logger.Logger, proj constructs.Project) {
	refs := proj.TempDeclRefs()
	for i := range refs.Count() {
		resolveTempDeclRef(log, proj, refs.Get(i))
	}

	proj.AllConstructs().Foreach(func(c constructs.Construct) {
		if trc, has := c.(constructs.TempDeclRefContainer); has {
			trc.RemoveTempDeclRefs()
		}
	})
	proj.ClearAllTempDeclRefs()
}

func resolveTempDeclRef(log *logger.Logger, proj constructs.Project, ref constructs.TempDeclRef) {
	if ref.Resolved() {
		return
	}

	// Try to find instance of declaration or non-generic declaration.
	decl, ok := proj.FindDecl(ref.PackagePath(), ref.Name(), ref.InstanceTypes(), false, false)
	if ok {
		ref.SetResolution(decl)
		return
	}

	// Try to find generic declaration and then create the instance if needed.
	decl, ok = proj.FindDecl(ref.PackagePath(), ref.Name(), []constructs.TypeDesc{}, false, true)
	if !ok {
		panic(terror.New(`failed to find temp declaration referenced`).
			With(`package path`, ref.PackagePath()).
			With(`name`, ref.Name()).
			With(`instance types`, ref.InstanceTypes()))
	}
	if len(ref.InstanceTypes()) <= 0 {
		ref.SetResolution(decl)
		return
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

	case kind.InterfaceDecl:
		res := instantiator.InterfaceDecl(log, proj, nil, decl.(constructs.InterfaceDecl), nil, ref.InstanceTypes())
		if utils.IsNil(res) {
			panic(terror.New(`failed to resolve interface declaration reference`).
				With(`package path`, ref.PackagePath()).
				With(`name`, ref.Name()).
				With(`instance types`, ref.InstanceTypes()))
		}
		ref.SetResolution(res)

	case kind.Method:
		res := instantiator.Method(log, proj, decl.(constructs.Method), ref.InstanceTypes())
		if utils.IsNil(res) {
			panic(terror.New(`failed to resolve method declaration reference`).
				With(`package path`, ref.PackagePath()).
				With(`name`, ref.Name()).
				With(`instance types`, ref.InstanceTypes()))
		}
		ref.SetResolution(res)

	default:
		panic(terror.New(`unexpected declaration type`).
			With(`kind`, decl.Kind()).
			With(`decl`, decl))
	}
}
