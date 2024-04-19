package cyclomatic

import (
	"fmt"
	"go/ast"
	"go/parser"
	"os"
	"strings"
	"testing"

	"github.com/Snow-Gremlin/goToolbox/differs/diff"
	"github.com/Snow-Gremlin/goToolbox/testers/check"
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
		`─┬─[enter_8]`,
		` └─┬─[if_17]`,
		`   ├─┬─[ifBody_26]`,
		`   │ └─┬─[endIf_37]`,
		`   │   └───[exit_50]`,
		`   └───<endIf_37>`)
}

func Test_SimpleIfElse(t *testing.T) {
	c := parseFunc(t,
		`x := 9`,
		`if x > 7 {`,
		`  x = 4`,
		`} else {`,
		`  x = 2`,
		`  print("cat")`,
		`}`,
		`println(x)`)
	checkCyclo(t, c,
		`─┬─[enter_8]`,
		` └─┬─[if_17]`,
		`   ├─┬─[ifBody_26]`,
		`   │ └─┬─[endIf_69]`,
		`   │   └───[exit_82]`,
		`   └─┬─[elseBody_43]`,
		`     └───<endIf_69>`)
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

const writeMermaidFile = true
const mermaidFile = `mermaid.md`

var mermaidText = "# Cyclomatic Test Output\n"

func checkCyclo(t *testing.T, c *Cyclomatic, expLines ...string) {
	gotten := strings.TrimSpace(c.String())
	if writeMermaidFile {
		mermaidText = fmt.Sprintf("%s\n## %s\n\n```mermaid\n%s```\n", mermaidText, t.Name(), c.Mermaid())
		check.NoError(t).Require(os.WriteFile(mermaidFile, []byte(mermaidText), 0o666))
	}
	exp := strings.Join(expLines, "\n")
	if gotten != exp {
		gotLines := strings.Split(gotten, "\n")
		d := diff.Default().PlusMinus(gotLines, expLines)
		fmt.Println(strings.Join(d, "\n"))
		t.Fail()
	}
}
