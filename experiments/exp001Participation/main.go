package main

import (
	"fmt"
	"go/ast"
	"go/importer"
	"go/parser"
	"go/token"
	"go/types"
	"io/ioutil"
	"log"
	"os"
	"path"
	"sort"
	"strings"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println(`Must provide only the path to the base folder of the go code.`)
		fmt.Println(`Example: go run main.go /Users/<user name>/go/src/github.com/MSUSEL/msusel-tdmetrics-go`)
		return
	}
	basePath := os.Args[1]

	// Find all the
	fileInfo, err := ioutil.ReadDir(basePath)
	if err != nil {
		panic(err)
	}
	filenames := []string{}
	for _, info := range fileInfo {
		filename := path.Clean(path.Join(basePath, info.Name()))

		if strings.HasSuffix(filename, ".go") && !strings.HasSuffix(filename, "_test.go") {
			filenames = append(filenames, filename)
		}
	}
	sort.Strings(filenames)

	fileset := token.NewFileSet()
	files := []*ast.File{}
	for _, filename := range filenames {
		f, err := parser.ParseFile(fileset, filename, nil, parser.ParseComments)
		if err != nil {
			panic(err)
		}
		files = append(files, f)
	}

	conf := types.Config{Importer: importer.Default()}
	pkg, err := conf.Check(basePath, fileset, files, nil)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Package  %q\n", pkg.Path())
	fmt.Printf("Name:    %s\n", pkg.Name())
	fmt.Printf("Imports: %s\n", pkg.Imports())
	fmt.Printf("Scope:   %s\n", pkg.Scope())
}
