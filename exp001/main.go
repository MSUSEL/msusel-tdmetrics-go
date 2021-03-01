package main

import (
	"fmt"

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

	basePath := `C:\Data\Code\Go\src\github.com\grant-nelson\goDiff\`

	project := reader.New().
		SetBasePath(basePath).
		AddDir(basePath).
		FilterFilenames(filter.Default()).
		Read()

	fmt.Printf("Package: %q\n", project.Package.Path())
	fmt.Printf("Name:    %s\n", project.Package.Name())

	fmt.Println("Defines:")
	for id, obj := range project.Defs {
		fmt.Printf("  %q defines %v\n", id.Name, obj)
	}
	fmt.Println("Uses:")
	for id, obj := range project.Uses {
		fmt.Printf("  %q uses %v\n", id.Name, obj)
	}
	fmt.Println()
}
