package cyclomatic

import (
	"fmt"
	"go/ast"
	"go/parser"
	"strings"
	"testing"

	"github.com/Snow-Gremlin/goToolbox/differs/diff"
)

func Test_Simple(t *testing.T) {
	c := parseFunc(t,
		`x := max(1, 3, 5)`)
	checkCyclo(t, c,
		`─┬─[enter_8]`,
		` └───[exit_29]`)
}

func Test_SimpleWithReturn(t *testing.T) {
	c := parseFunc(t,
		`x := max(1, 3, 5)`,
		`return x`)
	checkCyclo(t, c,
		`─┬─[enter_8]`,
		` └───[exit_38]`)
}

func Test_SimpleWithDefer(t *testing.T) {
	c := parseFunc(t,
		`x := open()`,
		`defer x.close()`,
		`x.doStuff()`)
	checkCyclo(t, c,
		`─┬─[enter_8]`,
		` └─┬─[defer_22]`,
		`   └───[exit_51]`)
}

func Test_SimpleIf(t *testing.T) {
	c := parseFunc(t,
		`x := 9`,
		`if x > 7 {`,
		`  x = 4`,
		`}`,
		`println(x)`)
	checkCyclo(t, c,
		``)
}

/*
	print("A ")
	defer func() {
		print("B ")
	}()
	{
		print("C ")
		defer func() {
			print("D ")
		}()
		print("E ")
	}
	print("F ")

	// Output: A C E F D B
*/

/*
	print("A ")
	defer func() {
		print("B ")
	}()
	func() {
		print("C ")
		defer func() {
			print("D ")
		}()
		print("E ")
	}()
	print("F ")

	// Output: A C E D F B
*/

/*
	print("A ")
	for _ = range 4 {
		print("B ")
		defer func() {
			print("C ")
		}()
		print("D ")
	}
	print("E ")

	// Output: A B D B D B D B D E C C C C
*/

func parseFunc(t *testing.T, lines ...string) *Cyclomatic {
	code := fmt.Sprintf("func() {\n%s\n}\n", strings.Join(lines, "\n"))
	expr, err := parser.ParseExpr(code)
	if err != nil {
		t.Error(err)
	}
	block := expr.(*ast.FuncLit).Body
	return New(block)
}

func checkCyclo(t *testing.T, c *Cyclomatic, expLines ...string) {
	gotten := strings.TrimSpace(c.String())
	exp := strings.Join(expLines, "\n")
	if gotten != exp {
		d := diff.Default().PlusMinus(strings.Split(gotten, "\n"), expLines)
		fmt.Println(strings.Join(d, "\n"))
		t.Fail()
	}
}
