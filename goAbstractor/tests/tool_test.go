package tests

import (
	"encoding/json"
	"fmt"
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
	expPartials    = `/partial.yaml`
	writeOutFile   = `/out.json`
)

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

	verbose := true // testing.Verbose()
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

	var log *logger.Logger
	if verbose {
		log = logger.New()
		// Use the group filters to show specific algorithm logs while debugging.
		// log = log.Show(`analyze`)
		// log = log.Show(`converter`)
		// log = log.Show(`files`)
		// log = log.Show(`inheritance`)
		// log = log.Show(`instantiator`)
		// log = log.Show(`generateInterfaces`)
		// log = log.Show(`packages`)
		// log = log.Show(`usages`)
	}

	tt.proj = abstractor.Abstract(abstractor.Config{
		Packages: ps,
		Log:      log,
	})
	return tt
}

func (tt *testTool) readExp(expData any, file string) *testTool {
	expFile, err := os.ReadFile(pathToTestData + tt.dir + file)
	check.NoError(tt.t).
		Name(`Read expected json`).
		With(`Dir`, tt.dir).
		With(`File`, file).
		Require(err)

	err = yaml.Unmarshal(expFile, expData)
	check.NoError(tt.t).
		Name(`Unmarshal expected json`).
		With(`Dir`, tt.dir).
		With(`File`, file).
		Require(err)
	return tt
}

func (tt *testTool) full() *testTool {
	var expData any
	tt.readExp(&expData, expAbstraction)

	exp, err := json.MarshalIndent(expData, ``, `  `)
	check.NoError(tt.t).
		Name(`Marshal expected json`).
		With(`Dir`, tt.dir).
		Require(err)

	ctx := jsonify.NewContext()
	gotten, err := jsonify.Marshal(ctx, tt.proj)
	check.NoError(tt.t).
		Name(`Marshal project`).
		With(`Dir`, tt.dir).
		Require(err)

	if !slices.Equal(exp, gotten) {
		expLines := strings.Split(string(exp), "\n")
		gotLines := strings.Split(string(gotten), "\n")
		diffLines := diff.Default().PlusMinus(expLines, gotLines)
		tt.t.Error("\n" + strings.Join(diffLines, "\n"))
	}
	return tt
}

func (tt *testTool) dump() *testTool {
	fmt.Println()
	fmt.Println(tt.proj)
	return tt
}

type partialTest struct {
	Name      string `yaml:"name"`
	Path      []any  `yaml:"path"`
	Data      any    `yaml:"data"`
	OS        string `yaml:"os"`
	GoVersion string `yaml:"goVersion"`
}

func (tt *testTool) partial() *testTool {
	var partialTests []partialTest
	tt.readExp(&partialTests, expPartials)

	for _, pt := range partialTests {
		tt.runPartialTest(pt)
	}
	return tt
}

func (tt *testTool) runPartialTest(pt partialTest) {
	tt.t.Run(pt.Name, func(t *testing.T) {
		if len(pt.OS) > 0 && runtime.GOOS != pt.OS {
			t.Skip(`The OS changes the specific indices, this test is for ` + pt.OS + `.`)
		}
		if len(pt.GoVersion) > 0 && runtime.Version() != pt.GoVersion {
			t.Skip(`The Go version changes the specific indices, this test is for ` + pt.GoVersion + `.`)
		}

		ctx := jsonify.NewContext().IncludeDebugIndex(true)
		subData := tt.proj.ToJson(ctx).Seek(pt.Path)

		exp, err := json.MarshalIndent(pt.Data, ``, `  `)
		check.NoError(t).
			Name(`Marshal expected json`).
			With(`Dir`, tt.dir).
			With(`Path`, pt.Path).
			Require(err)

		gotten, err := json.MarshalIndent(subData, ``, `  `)
		check.NoError(t).
			Name(`Marshal project`).
			With(`Dir`, tt.dir).
			With(`Path`, pt.Path).
			Require(err)

		if !slices.Equal(exp, gotten) {
			expLines := strings.Split(string(exp), "\n")
			gotLines := strings.Split(string(gotten), "\n")
			diffLines := diff.Default().PlusMinus(expLines, gotLines)
			t.Error(strings.Join(diffLines, "\n"))
		}
	})
}

func (tt *testTool) save() *testTool {
	ctx := jsonify.NewContext().IncludeDebugIndex(true)
	gotten, err := jsonify.Marshal(ctx, tt.proj)
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
