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
	t.Skip(`If is not implemented yet.`) // TODO: Finish
	c := parseFunc(t,
		`x := 9`,
		`if x > 7 {`,
		`  x = 4`,
		`}`,
		`println(x)`)
	checkCyclo(t, c,
		``)
}

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
