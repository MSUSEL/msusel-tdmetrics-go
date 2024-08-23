package abstractor

import (
	"github.com/Snow-Gremlin/goToolbox/collections"
	"github.com/Snow-Gremlin/goToolbox/collections/sortedSet"

	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs/interfaceDesc"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/logger"
)

func (ab *abstractor) resolveInheritance() {
	ab.log.Log(`resolve inheritance`)
	its := ab.proj.InterfaceDescs()
	roots := sortedSet.New(interfaceDesc.Comparer())
	log2 := ab.log.Group(`inheritance`)
	log3 := log2.Prefix(`  `)
	for i := range its.Count() {
		log2.Logf(`--(%d): %s`, i, its.Get(i))
		addInheritance(roots, its.Get(i), log3)
	}
	// throw away roots, they are no longer needed.
}

func addInheritance(siblings collections.SortedSet[constructs.InterfaceDesc], it constructs.InterfaceDesc, log *logger.Logger) {
	log2 := log.Prefix(` |  `)
	add := siblings.Count() <= 0
	for i := siblings.Count() - 1; i >= 0; i-- {
		a := siblings.Get(i)
		switch {
		case a.Implements(it):
			// Yi <: X
			log.Logf(` |--(%d) Yi<:X %s`, i, a)
			addInheritance(a.Inherits(), it, log2)
		case it.Implements(a):
			// Yi :> X
			log.Logf(` |--(%d) Yi:>X %s`, i, a)
			it.AddInherits(a)
			siblings.RemoveRange(i, 1)
			add = true
		default:
			// Possible overlap, check for super-types in subtree.
			log.Logf(` |--(%d) else  %s`, i, a)
			seekInherits(a.Inherits(), it, log2)
			add = true
		}
	}
	if add {
		log.Log(` '-- add`)
		siblings.Add(it)
	} else {
		log.Log(` '-- no-op`)
	}
}

func seekInherits(siblings collections.SortedSet[constructs.InterfaceDesc], it constructs.InterfaceDesc, log *logger.Logger) {
	for i := siblings.Count() - 1; i >= 0; i-- {
		a := siblings.Get(i)
		if it.Implements(a) {
			log.Log(` - `, a)
			it.AddInherits(a)
		} else {
			seekInherits(a.Inherits(), it, log)
		}
	}
}
