package tests

import (
	"encoding/json"
	"os"
	"runtime"
	"slices"
	"strings"
	"testing"

	"github.com/Snow-Gremlin/goToolbox/differs/diff"
	"github.com/Snow-Gremlin/goToolbox/testers/check"
	"gopkg.in/yaml.v3"

	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/abstractor"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/jsonify"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/logger"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/reader"
)

const (
	pathToTestData = `../../testData/go/`
	expAbstraction = `/abstraction.yaml`
	writeOutFile   = `/out.json`
)

func Test_T0001(t *testing.T) { newTest(t, `test0001`).abstract().equals() }
func Test_T0002(t *testing.T) { newTest(t, `test0002`).abstract().equals() }
func Test_T0003(t *testing.T) { newTest(t, `test0003`).abstract().equals() }
func Test_T0004(t *testing.T) { newTest(t, `test0004`).abstract().equals() }
func Test_T0005(t *testing.T) { newTest(t, `test0005`).abstract(`cats.go`).equals() }
func Test_T0006(t *testing.T) {
	if runtime.GOOS != `windows` {
		t.Skip(`The OS changes the specific type indices, this test is for Windows.`)
	}
	newTest(t, `test0006`).abstract(`cats.go`).partial().save()
}
func Test_T0007(t *testing.T) { newTest(t, `test0007`).abstract().equals() }
func Test_T0008(t *testing.T) { newTest(t, `test0008`).abstract().equals() }

func newTest(t *testing.T, dir string) *testTool {
	return &testTool{
		t:   t,
		dir: dir,
	}
}

type testTool struct {
	t    *testing.T
	dir  string
	proj constructs.Project
}

func (tt *testTool) abstract(patterns ...string) *testTool {
	if len(patterns) <= 0 {
		patterns = []string{`main.go`}
	}

	verbose := testing.Verbose()
	basePath := pathToTestData + tt.dir
	ps, err := reader.Read(&reader.Config{
		Verbose:    verbose,
		Dir:        pathToTestData + tt.dir,
		Patterns:   patterns,
		BuildFlags: []string{`-tags=test`},
	})
	check.NoError(tt.t).
		Name(`Read project`).
		With(`Dir`, tt.dir).
		Require(err)

	tt.proj = abstractor.Abstract(abstractor.Config{
		Packages: ps,
		Log:      logger.New(verbose),
		BasePath: basePath,
	})
	return tt
}

func (tt *testTool) readExp(expData any) *testTool {
	expFile, err := os.ReadFile(pathToTestData + tt.dir + expAbstraction)
	check.NoError(tt.t).
		Name(`Read expected json`).
		With(`Dir`, tt.dir).
		Require(err)

	err = yaml.Unmarshal(expFile, expData)
	check.NoError(tt.t).
		Name(`Unmarshal expected json`).
		With(`Dir`, tt.dir).
		Require(err)
	return tt
}

func (tt *testTool) equals() *testTool {
	var expData any
	tt.readExp(&expData)

	exp, err := json.MarshalIndent(expData, ``, `  `)
	check.NoError(tt.t).
		Name(`Marshal expected json`).
		With(`Dir`, tt.dir).
		Require(err)

	gotten, err := jsonify.Marshal(jsonify.NewContext(), tt.proj)
	check.NoError(tt.t).
		Name(`Marshal project`).
		With(`Dir`, tt.dir).
		Require(err)

	if !slices.Equal(exp, gotten) {
		expLines := strings.Split(string(exp), "\n")
		gotLines := strings.Split(string(gotten), "\n")
		diffLines := diff.Default().PlusMinus(expLines, gotLines)
		tt.t.Error(strings.Join(diffLines, "\n"))
	}
	return tt
}

func (tt *testTool) partial() *testTool {
	var expParts []struct {
		Path []any `yaml:"path"`
		Data any   `yaml:"data"`
	}
	tt.readExp(&expParts)

	for _, part := range expParts {
		ctx := jsonify.NewContext()
		subData := tt.proj.ToJson(ctx).Seek(part.Path)

		exp, err := json.MarshalIndent(part.Data, ``, `  `)
		check.NoError(tt.t).
			Name(`Marshal expected json`).
			With(`Dir`, tt.dir).
			With(`Path`, part.Path).
			Require(err)

		gotten, err := json.MarshalIndent(subData, ``, `  `)
		check.NoError(tt.t).
			Name(`Marshal project`).
			With(`Dir`, tt.dir).
			With(`Path`, part.Path).
			Require(err)

		if !slices.Equal(exp, gotten) {
			expLines := strings.Split(string(exp), "\n")
			gotLines := strings.Split(string(gotten), "\n")
			diffLines := diff.Default().PlusMinus(expLines, gotLines)
			tt.t.Error(strings.Join(diffLines, "\n"))
		}
	}

	return tt
}

func (tt *testTool) save() *testTool {
	gotten, err := jsonify.Marshal(jsonify.NewContext(), tt.proj)
	check.NoError(tt.t).
		Name(`Marshal project`).
		With(`Dir`, tt.dir).
		Require(err)

	err = os.WriteFile(pathToTestData+tt.dir+writeOutFile, gotten, 0o644)
	check.NoError(tt.t).
		Name(`Save project`).
		With(`Dir`, tt.dir).
		Require(err)
	return tt
}
