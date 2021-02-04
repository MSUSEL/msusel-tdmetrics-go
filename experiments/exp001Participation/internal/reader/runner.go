package reader

import (
	"fmt"
	"go/ast"
	"go/importer"
	"go/parser"
	"go/token"
	"go/types"
	"log"
	"runtime"
	"sort"
)

type runner struct {
}

func newRunner() *runner {
	return &runner{}
}

func (r *runner) run(basePath string, sources map[string]interface{}) {
	// Sort the file names
	filenames := []string{}
	for filename := range sources {
		filenames = append(filenames, filename)
	}
	sort.Strings(filenames)

	// Read and parse all the sources
	fileset := token.NewFileSet()
	files := []*ast.File{}
	for _, filename := range filenames {
		source := sources[filename]
		f, err := parser.ParseFile(fileset, filename, source, parser.ParseComments)
		if err != nil {
			panic(err)
		}
		files = append(files, f)
	}

	// Resolve types in the packages.
	imp := importer.ForCompiler(fileset, runtime.Compiler, nil)
	conf := types.Config{Importer: imp}
	pkg, err := conf.Check(basePath, fileset, files, nil)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Package  %q\n", pkg.Path())
	fmt.Printf("Name:    %s\n", pkg.Name())
	fmt.Printf("Imports: %s\n", pkg.Imports())
	fmt.Printf("Scope:   %s\n", pkg.Scope())
}
