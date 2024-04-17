package cyclomatic

import (
	"fmt"
	"go/ast"
	"go/parser"
	"strings"
	"testing"
)

func Test_Simple(t *testing.T) {
	c := parseFunc(t,
		`x := max(1, 3, 5)`)
	checkCyclo(t, c,
		`─┬─[start_8]`,
		` └───[exit_29]`)
}

func Test_SimpleWithReturn(t *testing.T) {
	c := parseFunc(t,
		`x := max(1, 3, 5)`,
		`return x`)
	checkCyclo(t, c,
		`─┬─[start_8]`,
		` └───[exit_38]`)
}

func Test_SimpleWithDefer(t *testing.T) {
	c := parseFunc(t,
		`x := open()`,
		`defer x.close()`,
		`x.doStuff()`)
	checkCyclo(t, c,
		`─┬─[start_8]`,
		` └─┬─[defer_22]`,
		`   └───[exit_51]`)
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
		fmt.Println(`Gotten:`)
		fmt.Println(gotten)
		fmt.Println(`Expected:`)
		fmt.Println(exp)
		t.Fail()
	}
}
