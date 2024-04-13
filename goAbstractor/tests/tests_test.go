package tests

import (
	"encoding/json"
	"os"
	"slices"
	"strings"
	"testing"

	"github.com/Snow-Gremlin/goToolbox/differs/diff"
	"github.com/Snow-Gremlin/goToolbox/testers/check"

	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/abstractor"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/reader"
)

func runTest(t *testing.T, dir string) {
	const verbose = true
	ps, err := reader.Read(&reader.Config{
		Verbose:    verbose,
		Dir:        `./` + dir,
		Patterns:   []string{`main.go`},
		BuildFlags: []string{`-tags=test`},
	})
	check.NoError(t).Name(`Read project`).With(`Dir`, dir).Require(err)

	proj := abstractor.Abstract(ps, verbose)

	expFile, err := os.ReadFile(`./` + dir + `/expected.json`)
	check.NoError(t).Name(`Read expected json`).With(`Dir`, dir).Require(err)

	var expData any
	err = json.Unmarshal(expFile, &expData)
	check.NoError(t).Name(`Unmarshal expected json`).With(`Dir`, dir).Require(err)

	exp, err := json.MarshalIndent(expData, ``, `  `)
	check.NoError(t).Name(`Marshal expected json`).With(`Dir`, dir).Require(err)

	gotten, err := json.MarshalIndent(proj, ``, `  `)
	check.NoError(t).Name(`Marshal project`).With(`Dir`, dir).Require(err)

	if !slices.Equal(exp, gotten) {
		expLines := strings.Split(string(exp), "\n")
		gotLines := strings.Split(string(gotten), "\n")
		diffLines := diff.Default().PlusMinus(expLines, gotLines)
		t.Error(strings.Join(diffLines, "\n"))
	}
}

func Test_T0001(t *testing.T) { runTest(t, `test0001`) }

func Test_T0002(t *testing.T) { runTest(t, `test0002`) }
