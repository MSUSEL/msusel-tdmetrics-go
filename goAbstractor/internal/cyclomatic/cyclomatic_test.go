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

func Test_SimpleIfElseIf(t *testing.T) {
	c := parseFunc(t,
		`x := 9`,
		`if x > 7 {`,
		`  x = 4`,
		`} else if x > 4 {`,
		`  x = 3`,
		`}`,
		`println(x)`)
	checkCyclo(t, c,
		`─┬─[enter_8]`,
		` └─┬─[if_17]`,
		`   ├─┬─[ifBody_26]`,
		`   │ └─┬─[endIf_63]`,
		`   │   └───[exit_76]`,
		`   └─┬─[elseBody_43]`,
		`     └─┬─[if_43]`,
		`       ├─┬─[ifBody_52]`,
		`       │ └───<endIf_63>`,
		`       └───<endIf_63>`)
}

func Test_SimpleIfElseIfElse(t *testing.T) {
	c := parseFunc(t,
		`x := 9`,
		`if x > 7 {`,
		`  x = 4`,
		`} else if x > 4 {`,
		`  x = 3`,
		`} else {`,
		`  x = 2`,
		`}`,
		`println(x)`)
	checkCyclo(t, c,
		`─┬─[enter_8]`,
		` └─┬─[if_17]`,
		`   ├─┬─[ifBody_26]`,
		`   │ └─┬─[endIf_80]`,
		`   │   └───[exit_93]`,
		`   └─┬─[elseBody_43]`,
		`     └─┬─[if_43]`,
		`       ├─┬─[ifBody_52]`,
		`       │ └───<endIf_80>`,
		`       └─┬─[elseBody_69]`,
		`         └───<endIf_80>`)
}

func Test_DeferInBlock(t *testing.T) {
	c := parseFunc(t,
		`print("A ")`,
		`defer func() {`,
		`	print("B ")`,
		`}()`,
		`{`,
		`	print("C ")`,
		`	defer func() {`,
		`		print("D ")`,
		`	}()`,
		`	print("E ")`,
		`}`,
		`print("F ")`)
	// Output: A C E F D B
	checkCyclo(t, c,
		`─┬─[enter_8]`,
		` └─┬─[defer_70]`,
		`   └─┬─[defer_22]`,
		`     └───[exit_132]`)
}

func Test_DeferInFuncLiteral(t *testing.T) {
	c := parseFunc(t,
		`print("A ")`,
		`defer func() {`,
		`	print("B ")`,
		`}()`,
		`func() {`,
		`	print("C ")`,
		`	defer func() {`,
		`		print("D ")`,
		`	}()`,
		`	print("E ")`,
		`}()`,
		`print("F ")`)
	// Output: A C E D F B
	checkCyclo(t, c,
		``) // TODO: Need to include body in func literal
}

func Test_DeferWithComplexity(t *testing.T) {
	c := parseFunc(t,
		`print("A ")`,
		`defer func() {`,
		`   if r := recover(); r != nil {`,
		`		print("B ")`,
		`       return`,
		`   }`,
		`	print("C ")`,
		`}()`,
		`print("D ")`)
	checkCyclo(t, c,
		``) // TODO: Need to include if in defer.
}

func Test_ForRangeWithDefer(t *testing.T) {
	c := parseFunc(t,
		`print("A ")`,
		`for _ = range 4 {`,
		`	print("B ")`,
		`	defer func() {`,
		`		print("C ")`,
		`	}()`,
		`	print("D ")`,
		`}`,
		`print("E ")`)
	// Output: A B D B D B D B D E C C C C
	checkCyclo(t, c,
		``)
}

func Test_InfiniteGotoLoop(t *testing.T) {
	c := parseFunc(t,
		`print("A ")`,
		`BEANS:`,
		`print("B ")`,
		`goto BEANS`,
		`print("C ")`) // Unreachable
	// Output: A B B B...
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
