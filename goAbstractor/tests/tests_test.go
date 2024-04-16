package tests

import (
	"encoding/json"
	"os"
	"slices"
	"strings"
	"testing"

	"github.com/Snow-Gremlin/goToolbox/differs/diff"
	"github.com/Snow-Gremlin/goToolbox/testers/check"
	"gopkg.in/yaml.v3"

	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/abstractor"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/jsonify"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/reader"
)

func Test_T0001(t *testing.T) { runTest(t, `test0001`) }
func Test_T0002(t *testing.T) { runTest(t, `test0002`) }
func Test_T0003(t *testing.T) { runTest(t, `test0003`) }
func Test_T0004(t *testing.T) { runTest(t, `test0004`) }

func runTest(t *testing.T, dir string, patterns ...string) {
	if len(patterns) <= 0 {
		patterns = []string{`main.go`}
	}
	const verbose = true
	ps, err := reader.Read(&reader.Config{
		Verbose:    verbose,
		Dir:        `./` + dir,
		Patterns:   patterns,
		BuildFlags: []string{`-tags=test`},
	})
	check.NoError(t).Name(`Read project`).With(`Dir`, dir).Require(err)

	proj := abstractor.Abstract(ps, verbose)

	expFile, err := os.ReadFile(`./` + dir + `/expected.yaml`)
	check.NoError(t).Name(`Read expected json`).With(`Dir`, dir).Require(err)

	var expData any
	err = yaml.Unmarshal(expFile, &expData)
	check.NoError(t).Name(`Unmarshal expected json`).With(`Dir`, dir).Require(err)

	exp, err := json.MarshalIndent(expData, ``, `  `)
	check.NoError(t).Name(`Marshal expected json`).With(`Dir`, dir).Require(err)

	gotten, err := jsonify.Marshal(jsonify.NewContext(), proj)
	check.NoError(t).Name(`Marshal project`).With(`Dir`, dir).Require(err)

	if !slices.Equal(exp, gotten) {
		expLines := strings.Split(string(exp), "\n")
		gotLines := strings.Split(string(gotten), "\n")
		diffLines := diff.Default().PlusMinus(expLines, gotLines)
		t.Error(strings.Join(diffLines, "\n"))
	}
}
