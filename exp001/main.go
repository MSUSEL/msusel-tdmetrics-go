package main

import (
	"fmt"
	"os"
	"path"

	"github.com/MSUSEL/msusel-tdmetrics-go/exp001/internal/filter"
	"github.com/MSUSEL/msusel-tdmetrics-go/exp001/internal/reader"
)

func main() {
	// if len(os.Args) != 2 {
	// 	fmt.Println(`Must provide only the path to the base folder of the go code.`)
	// 	fmt.Println(`Example: go run main.go /Users/<user name>/go/src/github.com/MSUSEL/msusel-tdmetrics-go`)
	// 	return
	// }
	// basePath := os.Args[1]

	basePath := path.Join(os.Getenv("GOPATH"), `src/github.com/Grant-Nelson/goDiff/`)

	project := reader.New().
		SetBasePath(basePath).
		AddDir(basePath).
		FilterFilenames(filter.Default()).
		Read()

	fmt.Printf("Package: %q\n", project.Package.Path())
	fmt.Printf("Name:    %s\n", project.Package.Name())

	// fmt.Println("===========================")
	// for t, ids := range project.UsedSignatures() {
	// 	fmt.Println(t, "=>", ids)
	// }
	// fmt.Println("===========================")
	// for id, f := range project.DefinedFuncs() {
	// 	fmt.Println(id, "=> ", f)
	// }
	fmt.Println("===========================")
	project.TypeDefs()
	fmt.Println()
}
