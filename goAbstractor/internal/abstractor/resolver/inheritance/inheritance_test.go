package inheritance

import (
	"sort"
	"strings"
	"testing"

	"github.com/Snow-Gremlin/goToolbox/collections"
	"github.com/Snow-Gremlin/goToolbox/collections/sortedSet"
	"github.com/Snow-Gremlin/goToolbox/comp"
	"github.com/Snow-Gremlin/goToolbox/differs/diff"
	"github.com/Snow-Gremlin/goToolbox/testers/check"

	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/logger"
)

func Test_Diamond_Right(t *testing.T) {
	in := newTest(t)
	in.Check()

	in.Process(`A`, `B`, `C`)
	in.Check(
		`╶──{A, B, C}`)

	in.Process(`A`, `B`)
	in.Check(
		`╶──{A, B, C}`,
		`   └──{A, B}`)

	in.Process(`B`, `C`)
	in.Check(
		`╶──{A, B, C}`,
		`   ├──{A, B}`,
		`   └──{B, C}`)

	in.Process(`B`)
	in.Check(
		`╶──{A, B, C}`,
		`   ├──{A, B}`,
		`   │  └──{B}`,
		`   └──{B, C}`,
		`      └──{B}`)
}

func Test_Diamond_Left(t *testing.T) {
	in := newTest(t)
	in.Process(`B`)
	in.Check(
		`╶──{B}`)

	in.Process(`B`, `C`)
	in.Check(
		`╶──{B, C}`,
		`   └──{B}`)

	in.Process(`A`, `B`)
	in.Check(
		`┌──{A, B}`,
		`│  └──{B}`,
		`└──{B, C}`,
		`   └──{B}`)

	in.Process(`A`, `B`, `C`)
	in.Check(
		`╶──{A, B, C}`,
		`   ├──{A, B}`,
		`   │  └──{B}`,
		`   └──{B, C}`,
		`      └──{B}`)
}

func Test_Diamond_Middle(t *testing.T) {
	in := newTest(t)
	in.Process(`B`)
	in.Check(
		`╶──{B}`)

	in.Process(`A`, `B`, `C`)
	in.Check(
		`╶──{A, B, C}`,
		`   └──{B}`)

	in.Process(`B`, `C`)
	in.Check(
		`╶──{A, B, C}`,
		`   └──{B, C}`,
		`      └──{B}`)

	in.Process(`A`, `B`)
	in.Check(
		`╶──{A, B, C}`,
		`   ├──{A, B}`,
		`   │  └──{B}`,
		`   └──{B, C}`,
		`      └──{B}`)
}

func Test_PickingUpLeaves(t *testing.T) {
	in := newTest(t)
	in.Process(`A`)
	in.Process(`B`)
	in.Process(`C`)
	in.Check(
		`┌──{A}`,
		`├──{B}`,
		`└──{C}`)

	in.Process(`A`, `B`, `C`)
	in.Check(
		`╶──{A, B, C}`,
		`   ├──{A}`,
		`   ├──{B}`,
		`   └──{C}`)

	in.Process(`A`, `B`)
	in.Check(
		`╶──{A, B, C}`,
		`   ├──{A, B}`,
		`   │  ├──{A}`,
		`   │  └──{B}`,
		`   └──{C}`)

	in.Process(`B`, `C`)
	in.Check(
		`╶──{A, B, C}`,
		`   ├──{A, B}`,
		`   │  ├──{A}`,
		`   │  └──{B}`,
		`   └──{B, C}`,
		`      ├──{B}`,
		`      └──{C}`)

	in.Process(`A`, `C`)
	in.Check(
		`╶──{A, B, C}`,
		`   ├──{A, B}`,
		`   │  ├──{A}`,
		`   │  └──{B}`,
		`   ├──{A, C}`,
		`   │  ├──{A}`,
		`   │  └──{C}`,
		`   └──{B, C}`,
		`      ├──{B}`,
		`      └──{C}`)
}

func Test_AddingLeaves(t *testing.T) {
	in := newTest(t)
	in.Process(`A`, `B`)
	in.Process(`B`, `C`)
	in.Process(`A`, `C`)
	in.Check(
		`┌──{A, B}`,
		`├──{A, C}`,
		`└──{B, C}`)

	in.Process(`A`, `B`, `C`)
	in.Check(
		`╶──{A, B, C}`,
		`   ├──{A, B}`,
		`   ├──{A, C}`,
		`   └──{B, C}`)

	in.Process(`A`)
	in.Check(
		`╶──{A, B, C}`,
		`   ├──{A, B}`,
		`   │  └──{A}`,
		`   ├──{A, C}`,
		`   │  └──{A}`,
		`   └──{B, C}`)

	in.Process(`B`)
	in.Check(
		`╶──{A, B, C}`,
		`   ├──{A, B}`,
		`   │  ├──{A}`,
		`   │  └──{B}`,
		`   ├──{A, C}`,
		`   │  └──{A}`,
		`   └──{B, C}`,
		`      └──{B}`)

	in.Process(`C`)
	in.Check(
		`╶──{A, B, C}`,
		`   ├──{A, B}`,
		`   │  ├──{A}`,
		`   │  └──{B}`,
		`   ├──{A, C}`,
		`   │  ├──{A}`,
		`   │  └──{C}`,
		`   └──{B, C}`,
		`      ├──{B}`,
		`      └──{C}`)
}

func Test_InjectBetween(t *testing.T) {
	in := newTest(t)
	in.Process(`A`)
	in.Process(`A`, `B`, `C`)
	in.Process(`A`, `B`, `D`)
	in.Check(
		`┌──{A, B, C}`,
		`│  └──{A}`,
		`└──{A, B, D}`,
		`   └──{A}`)

	in.Process(`A`, `B`)
	in.Check(
		`┌──{A, B, C}`,
		`│  └──{A, B}`,
		`│     └──{A}`,
		`└──{A, B, D}`,
		`   └──{A, B}`,
		`      └──{A}`)
}

func Test_TNodeTest(t *testing.T) {
	n1 := newTNode(`A`, `B`, `C`)
	n2 := newTNode(`A`, `B`)
	n3 := newTNode(`B`, `C`)
	n4 := newTNode(`A`)
	n5 := newTNode(`B`)
	n6 := newTNode(`C`)
	n7 := newTNode(`D`)

	check.True(t).Withf(`test`, `%v implements %v`, n1, n1).Assert(n1.Implements(n1))
	check.True(t).Withf(`test`, `%v implements %v`, n1, n2).Assert(n1.Implements(n2))
	check.True(t).Withf(`test`, `%v implements %v`, n1, n3).Assert(n1.Implements(n3))
	check.True(t).Withf(`test`, `%v implements %v`, n1, n4).Assert(n1.Implements(n4))
	check.True(t).Withf(`test`, `%v implements %v`, n1, n5).Assert(n1.Implements(n5))
	check.True(t).Withf(`test`, `%v implements %v`, n1, n6).Assert(n1.Implements(n6))
	check.False(t).Withf(`test`, `%v implements %v`, n1, n7).Assert(n1.Implements(n7))

	check.False(t).Withf(`test`, `%v implements %v`, n2, n1).Assert(n2.Implements(n1))
	check.True(t).Withf(`test`, `%v implements %v`, n2, n2).Assert(n2.Implements(n2))
	check.False(t).Withf(`test`, `%v implements %v`, n2, n3).Assert(n2.Implements(n3))
	check.True(t).Withf(`test`, `%v implements %v`, n2, n4).Assert(n2.Implements(n4))
	check.True(t).Withf(`test`, `%v implements %v`, n2, n5).Assert(n2.Implements(n5))
	check.False(t).Withf(`test`, `%v implements %v`, n2, n6).Assert(n2.Implements(n6))
	check.False(t).Withf(`test`, `%v implements %v`, n2, n7).Assert(n2.Implements(n7))

	check.False(t).Withf(`test`, `%v implements %v`, n3, n1).Assert(n3.Implements(n1))
	check.False(t).Withf(`test`, `%v implements %v`, n3, n2).Assert(n3.Implements(n2))
	check.True(t).Withf(`test`, `%v implements %v`, n3, n3).Assert(n3.Implements(n3))
	check.False(t).Withf(`test`, `%v implements %v`, n3, n4).Assert(n3.Implements(n4))
	check.True(t).Withf(`test`, `%v implements %v`, n3, n5).Assert(n3.Implements(n5))
	check.True(t).Withf(`test`, `%v implements %v`, n3, n6).Assert(n3.Implements(n6))
	check.False(t).Withf(`test`, `%v implements %v`, n3, n7).Assert(n3.Implements(n7))

	check.True(t).Withf(`test`, `%v implements %v`, n4, n4).Assert(n4.Implements(n4))
	check.False(t).Withf(`test`, `%v implements %v`, n4, n5).Assert(n4.Implements(n5))
	check.False(t).Withf(`test`, `%v implements %v`, n5, n4).Assert(n5.Implements(n4))
}

type tInheritance struct {
	t  *testing.T
	in *inheritanceImp[*tNode]
}

func newTest(t *testing.T) tInheritance {
	in := New(Compare(), logger.New())
	return tInheritance{
		t:  t,
		in: in.(*inheritanceImp[*tNode]),
	}
}

func (ti tInheritance) Process(parts ...string) {
	ti.in.Process(newTNode(parts...))
	ti.in.log.Log()
}

func (ti tInheritance) String() string {
	buf := &strings.Builder{}
	touched := map[*tNode]bool{}
	write := func(text string) {
		if _, err := buf.WriteString(text); err != nil {
			panic(err)
		}
	}
	writeNodes(ti.in.roots, write, touched, true, ``)
	return buf.String()
}

func writeNodes(siblings collections.SortedSet[*tNode], write func(text string),
	touched map[*tNode]bool, first bool, indent string,
) {
	count := siblings.Count()
	for i := range count {
		node := siblings.Get(i)

		nodePrefix, subIndent := `├──`, `│  `
		if first && i <= 0 {
			nodePrefix = `┌──`
			if i+1 >= count {
				nodePrefix, subIndent = `╶──`, `   `
			}
			first = false
		} else if i+1 >= count {
			nodePrefix, subIndent = `└──`, `   `
		}

		write(indent + nodePrefix + node.String())

		if touched[node] {
			write("[Loop]\n")
		} else {
			write("\n")
			touched[node] = true
			writeNodes(node.Inherits(), write, touched, false, indent+subIndent)
			touched[node] = false
		}
	}
}

func (ti tInheritance) Check(expLines ...string) {
	exp := strings.Join(expLines, "\n")
	result := strings.TrimSpace(ti.String())
	if exp != result {
		resultLines := strings.Split(result, "\n")
		d := diff.Default().PlusMinus(resultLines, expLines)
		ti.t.Error("\n" + strings.Join(d, "\n"))
	}
}

type tNode struct {
	parts   []string
	parents collections.SortedSet[*tNode]
}

func newTNode(parts ...string) *tNode {
	sort.Strings(parts)
	return &tNode{
		parts:   parts,
		parents: sortedSet.New(Compare()),
	}
}

func Compare() comp.Comparer[*tNode] {
	partComp := comp.Slice[[]string](comp.Ordered[string]())
	return func(a, b *tNode) int {
		return partComp(a.parts, b.parts)
	}
}

func (tn *tNode) Implements(other *tNode) bool {
	j, jCount := 0, len(other.parts)
	for i, iCount := 0, len(tn.parts); i < iCount; i++ {
		if j >= jCount {
			return true
		}
		cmp := strings.Compare(tn.parts[i], other.parts[j])
		switch {
		case cmp > 0:
			return false
		case cmp == 0:
			j++
		}
	}
	return j >= jCount
}

func (tn *tNode) AddInherits(parent *tNode) *tNode {
	v, _ := tn.parents.TryAdd(parent)
	return v
}

func (tn *tNode) Inherits() collections.SortedSet[*tNode] {
	return tn.parents
}

func (tn *tNode) String() string {
	return `{` + strings.Join(tn.parts, `, `) + `}`
}
