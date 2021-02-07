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

	// Prepare the info for collecting data.
	info := &types.Info{
		Types: make(map[ast.Expr]types.TypeAndValue),
		Defs:  make(map[*ast.Ident]types.Object),
		Uses:  make(map[*ast.Ident]types.Object),
	}

	// Resolve types in the packages.
	imp := importer.ForCompiler(fileset, runtime.Compiler, nil)
	conf := types.Config{Importer: imp}
	pkg, err := conf.Check(basePath, fileset, files, info)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Package: %q\n", pkg.Path())
	fmt.Printf("Name:    %s\n", pkg.Name())

	fmt.Println("Defines:")
	for id, obj := range info.Defs {
		fmt.Printf("  %q defines %v\n", id.Name, obj)
	}
	fmt.Println("Uses:")
	for id, obj := range info.Uses {
		fmt.Printf("  %q uses %v\n", id.Name, obj)
	}
	fmt.Println()

}
