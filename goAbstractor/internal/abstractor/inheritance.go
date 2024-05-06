package abstractor

import (
	"errors"
	"fmt"

	"github.com/Snow-Gremlin/goToolbox/utils"

	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs/typeDesc"
)

func (ab *abstractor) resolveInheritance() {
	ab.log(`resolve inheritance`)
	if len(ab.proj.AllInterfaces) <= 0 {
		panic(errors.New(`expected the object interface at minimum but found no interfaces`))
	}
	obj := ab.proj.AllInterfaces[0]
	if len(obj.Methods) > 0 {
		panic(fmt.Errorf(`expected the first interface to be the object interface but it had %d methods`, len(obj.Methods)))
	}
	for _, inter := range ab.proj.AllInterfaces[1:] {
		findInheritors(obj, inter)
	}
	for _, inter := range ab.proj.AllInterfaces {
		setInheritance(inter)
	}
	for _, p := range ab.proj.Packages {
		for _, td := range p.Types {
			if _, ok := td.Type.(*typeDesc.Interface); ok {
				// Skip over interfaces which may not have bodied methods.
				continue
			}
			findImplements(obj, td)
		}
	}
}

func findInheritors(root, inter *typeDesc.Interface) bool {
	if root == inter {
		return true
	}
	if !inter.IsSupertypeOf(root) {
		return false
	}

	homed := false
	for _, other := range root.Inheritors {
		if findInheritors(other, inter) {
			homed = true
		}
	}
	if homed {
		return true
	}

	changed := false
	for i, other := range root.Inheritors {
		if other.IsSupertypeOf(inter) {
			inter.Inheritors = append(inter.Inheritors, other)
			root.Inheritors[i] = nil
			changed = true
		}
	}
	if changed {
		root.Inheritors = utils.RemoveZeros(root.Inheritors)
	}

	root.Inheritors = append(root.Inheritors, inter)
	return true
}

func setInheritance(inter *typeDesc.Interface) {
	for _, i := range inter.Inheritors {
		i.Inherits = append(i.Inherits, inter)
	}
}

func findImplements(root *typeDesc.Interface, td *constructs.TypeDef) bool {
	if !td.IsSupertypeOf(root) {
		return false
	}

	homed := false
	for _, other := range root.Inheritors {
		if findImplements(other, td) {
			homed = true
		}
	}
	if homed {
		return true
	}

	td.Inherits = append(td.Inherits, root)
	return true
}
