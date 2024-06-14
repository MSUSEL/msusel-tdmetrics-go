package tests

import (
	"encoding/json"
	"os"
	"slices"
	"strings"
	"testing"

	"github.com/Snow-Gremlin/goToolbox/collections/set"
	"github.com/Snow-Gremlin/goToolbox/differs/diff"
	"github.com/Snow-Gremlin/goToolbox/testers/check"
	"gopkg.in/yaml.v3"

	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/abstractor"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/jsonify"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/reader"
)

func Test_T0001(t *testing.T) { newTest(t, `test0001`).abstract().equals() }
func Test_T0002(t *testing.T) { newTest(t, `test0002`).abstract().equals() }
func Test_T0003(t *testing.T) { newTest(t, `test0003`).abstract().equals() }
func Test_T0004(t *testing.T) { newTest(t, `test0004`).abstract().equals() }
func Test_T0005(t *testing.T) { newTest(t, `test0005`).abstract(`cats.go`).equals() }
func Test_T0006(t *testing.T) { newTest(t, `test0006`).abstract(`cats.go`).partial(`Cat`) }
func Test_T0007(t *testing.T) { newTest(t, `test0007`).abstract().equals() }
func Test_T0008(t *testing.T) { newTest(t, `test0008`).abstract().equals() }

func newTest(t *testing.T, dir string) *testTool {
	expFile, err := os.ReadFile(`./` + dir + `/expected.yaml`)
	check.NoError(t).Name(`Read expected json`).With(`Dir`, dir).Require(err)

	var expData any
	err = yaml.Unmarshal(expFile, &expData)
	check.NoError(t).Name(`Unmarshal expected json`).With(`Dir`, dir).Require(err)

	return &testTool{
		t:       t,
		dir:     dir,
		expData: expData,
	}
}

type testTool struct {
	t       *testing.T
	dir     string
	proj    constructs.Project
	expData any
}

func (tt *testTool) abstract(patterns ...string) *testTool {
	if len(patterns) <= 0 {
		patterns = []string{`main.go`}
	}
	const verbose = true
	ps, err := reader.Read(&reader.Config{
		Verbose:    verbose,
		Dir:        `./` + tt.dir,
		Patterns:   patterns,
		BuildFlags: []string{`-tags=test`},
	})
	check.NoError(tt.t).Name(`Read project`).With(`Dir`, tt.dir).Require(err)
	tt.proj = abstractor.Abstract(ps, verbose)
	return tt
}

func (tt *testTool) equals() *testTool {
	exp, err := json.MarshalIndent(tt.expData, ``, `  `)
	check.NoError(tt.t).Name(`Marshal expected json`).With(`Dir`, tt.dir).Require(err)

	gotten, err := jsonify.Marshal(jsonify.NewContext(), tt.proj)
	check.NoError(tt.t).Name(`Marshal project`).With(`Dir`, tt.dir).Require(err)

	if !slices.Equal(exp, gotten) {
		expLines := strings.Split(string(exp), "\n")
		gotLines := strings.Split(string(gotten), "\n")
		diffLines := diff.Default().PlusMinus(expLines, gotLines)
		tt.t.Error(strings.Join(diffLines, "\n"))
	}
	return tt
}

func (tt *testTool) partial(packages ...string) *testTool {
	keep := set.With(packages...)
	tt.proj.FilterPackage(func(p constructs.Package) bool {
		return !keep.Contains(p.Name())
	})
	return tt.equals()
}
