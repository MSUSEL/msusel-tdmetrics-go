package references

import (
	"fmt"

	"github.com/Snow-Gremlin/goToolbox/terrors/terror"
	"github.com/Snow-Gremlin/goToolbox/utils"

	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/abstractor/instantiator"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs/kind"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/logger"
)

func References(log *logger.Logger, proj constructs.Project, required bool) bool {
	defer func() {
		if rex := recover(); rex != nil {
			fmt.Println("\n" + proj.String())
			panic(terror.RecoveredPanic(rex))
		}
	}()

	r := &ref{
		log:      log,
		proj:     proj,
		required: required,
		changed:  false,
	}
	r.removeTempReferences()
	r.removeTempDeclRefs()
	r.tempReferences()
	r.tempDeclRefs()
	r.clearAllTemps()
	return r.changed
}

type ref struct {
	log      *logger.Logger
	proj     constructs.Project
	required bool
	changed  bool
}

func (r *ref) clearAllTemps() {
	if r.required {
		r.proj.ClearAllTempReferences()
		r.proj.ClearAllTempTypeParamRefs()
		r.proj.ClearAllTempDeclRefs()
	}
}

func (r *ref) removeTempReferences() {
	r.proj.AllConstructs().Foreach(func(c constructs.Construct) {
		if trc, has := c.(constructs.TempReferenceContainer); has {
			if trc.RemoveTempReferences(r.required) {
				r.changed = true
			}
		}
	})
}

func (r *ref) removeTempDeclRefs() {
	r.proj.AllConstructs().Foreach(func(c constructs.Construct) {
		if trc, has := c.(constructs.TempDeclRefContainer); has {
			if trc.RemoveTempDeclRefs(r.required) {
				r.changed = true
			}
		}
	})
}

func (r *ref) tempReferences() {
	refs := r.proj.TempReferences()
	for i := range refs.Count() {
		r.resolveTempRef(refs.Get(i))
	}
	r.removeTempReferences()
}

func (r *ref) tempDeclRefs() {
	refs := r.proj.TempDeclRefs()
	for i := range refs.Count() {
		r.resolveTempDeclRef(refs.Get(i))
	}
	r.removeTempDeclRefs()
}

func (r *ref) resolveTempRef(ref constructs.TempReference) {
	if ref.Resolved() {
		return
	}

	// Try to find instance of type or non-generic type.
	typ, ok := r.proj.FindType(ref.PackagePath(), ref.Name(), ref.Nest(), ref.ImplicitTypes(), ref.InstanceTypes(), false, false)
	if ok {
		ref.SetResolution(typ)
		r.changed = true
		return
	}

	// Try to find generic type and then create the instance if needed.
	typ, ok = r.proj.FindType(ref.PackagePath(), ref.Name(), ref.Nest(), nil, nil, false, r.required)
	if !ok {
		if !r.required {
			return
		}
		panic(terror.New(`failed to find temp referenced object`).
			With(`package path`, ref.PackagePath()).
			With(`name`, ref.Name()).
			With(`instance types`, ref.InstanceTypes()))
	}
	if len(ref.InstanceTypes()) <= 0 && len(ref.ImplicitTypes()) <= 0 {
		ref.SetResolution(typ)
		r.changed = true
		return
	}

	switch typ.Kind() {
	case kind.Object:
		res := instantiator.Object(r.log, r.proj, ref.GoType(), typ.(constructs.Object), ref.ImplicitTypes(), ref.InstanceTypes())
		if utils.IsNil(res) {
			panic(terror.New(`failed to resolve object type reference`).
				With(`package path`, ref.PackagePath()).
				With(`name`, ref.Name()).
				With(`instance types`, ref.InstanceTypes()))
		}
		ref.SetResolution(res)
		r.changed = true
		return

	case kind.InterfaceDecl:
		res := instantiator.InterfaceDecl(r.log, r.proj, ref.GoType(), typ.(constructs.InterfaceDecl), ref.ImplicitTypes(), ref.InstanceTypes())
		if utils.IsNil(res) {
			panic(terror.New(`failed to resolve interface type reference`).
				With(`package path`, ref.PackagePath()).
				With(`name`, ref.Name()).
				With(`instance types`, ref.InstanceTypes()))
		}
		ref.SetResolution(res)
		r.changed = true
		return

	default:
		panic(terror.New(`unexpected declaration type`).
			With(`kind`, typ.Kind()).
			With(`decl`, typ))
	}
}

func (r *ref) resolveTempDeclRef(ref constructs.TempDeclRef) {
	if ref.Resolved() {
		return
	}

	// Try to find instance of declaration or non-generic declaration.
	decl, ok := r.proj.FindDecl(ref.PackagePath(), ref.Name(), ref.Nest(), ref.ImplicitTypes(), ref.InstanceTypes(), false, false)
	if ok {
		ref.SetResolution(decl)
		r.changed = true
		return
	}

	// Try to find generic declaration and then create the instance if needed.
	decl, ok = r.proj.FindDecl(ref.PackagePath(), ref.Name(), ref.Nest(), nil, nil, false, r.required)
	if !ok {
		if !r.required {
			return
		}
		panic(terror.New(`failed to find temp declaration referenced`).
			With(`package path`, ref.PackagePath()).
			With(`name`, ref.Name()).
			With(`instance types`, ref.InstanceTypes()))
	}
	if len(ref.InstanceTypes()) <= 0 && len(ref.ImplicitTypes()) <= 0 {
		ref.SetResolution(decl)
		r.changed = true
		return
	}

	switch decl.Kind() {
	case kind.Object:
		res := instantiator.Object(r.log, r.proj, nil, decl.(constructs.Object), nil, ref.InstanceTypes())
		if utils.IsNil(res) {
			panic(terror.New(`failed to resolve object declaration reference`).
				With(`package path`, ref.PackagePath()).
				With(`name`, ref.Name()).
				With(`instance types`, ref.InstanceTypes()))
		}
		ref.SetResolution(res)
		r.changed = true
		return

	case kind.InterfaceDecl:
		res := instantiator.InterfaceDecl(r.log, r.proj, nil, decl.(constructs.InterfaceDecl), nil, ref.InstanceTypes())
		if utils.IsNil(res) {
			panic(terror.New(`failed to resolve interface declaration reference`).
				With(`package path`, ref.PackagePath()).
				With(`name`, ref.Name()).
				With(`instance types`, ref.InstanceTypes()))
		}
		ref.SetResolution(res)
		r.changed = true
		return

	case kind.Method:
		res := instantiator.Method(r.log, r.proj, decl.(constructs.Method), ref.InstanceTypes())
		if utils.IsNil(res) {
			panic(terror.New(`failed to resolve method declaration reference`).
				With(`package path`, ref.PackagePath()).
				With(`name`, ref.Name()).
				With(`instance types`, ref.InstanceTypes()))
		}
		ref.SetResolution(res)
		r.changed = true
		return

	default:
		panic(terror.New(`unexpected declaration type`).
			With(`kind`, decl.Kind()).
			With(`decl`, decl))
	}
}
