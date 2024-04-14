package abstractor

import (
	"errors"
	"fmt"

	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs/typeDesc"
)

func (ab *abstractor) resolveInheritance() {
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

	// TODO: Finish
}

func findInheritors(root, inter *typeDesc.Interface) bool {
	if !root.IsSupertypeOf(inter) {
		return false
	}

	homed := false
	for _, other := range inter.Inheritors {
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
			inter.Inheritors = append(inter.Inheritors, inter)
			inter.Inheritors[i] = nil
			changed = true
		}
	}
	if changed {
		inter.Inheritors = squeeze(inter.Inheritors)
	}

	root.Inheritors = append(root.Inheritors, inter)
	return true
}

func setInheritance(inter *typeDesc.Interface) {
	for _, i := range inter.Inheritors {
		i.Inherits = append(i.Inherits, inter)
	}
}
