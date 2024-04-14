package main

import (
	"fmt"
	"os"
	"runtime/debug"

	"github.com/Snow-Gremlin/goToolbox/argers/args"

	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/abstractor"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/jsonify"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/reader"
)

type argObject struct {
	ShowHelp bool   `args:"flag, h, help"`
	Verbose  bool   `args:"flag, v, verbose"`
	Minimize bool   `args:"flag, m, minimize"`
	InPath   string `args:"i, in"`
	OutPath  string `args:"o, out"`
}

func main() {
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("Panic: %v\n", r)
			debug.PrintStack()
			os.Exit(1)
		}
	}()

	ao := &argObject{}
	in := args.New().Struct(ao)
	if err := in.Process(os.Args[1:]); err != nil {
		fmt.Println(err.Error())
		fmt.Println(`Use "-h" argument to show help.`)
		os.Exit(1)
	}

	if ao.ShowHelp {
		fmt.Println(`Abstractor will read a Go project and output a JSON file`,
			`designed to be used in a design recovery and participation analysis`,
			`of the Go project.`)
		fmt.Println(os.Args[0], `<options> -i <inputPath> [ -o <outputPath> ]`)
		fmt.Println(`  --help|-h: Shows this help text.`)
		fmt.Println(`  --verbose|-v: Indicates the abstraction process should`,
			`output additional status information.`)
		fmt.Println(`  --minimize|-m: Indicates the JSON output should be`,
			`minimized instead of formatted.`)
		fmt.Println(`  --in|-i: The input path to the directory of the project`,
			`to read. The directory should have a go.mod file.`)
		fmt.Println(`  --out|-o: The output file path to write the JSON to.`,
			`If not given, the JSON will be outputted to the console.`)
		os.Exit(0)
	}

	ps, err := reader.Read(&reader.Config{
		Verbose: ao.Verbose,
		Dir:     ao.InPath,
	})
	if err != nil {
		fmt.Println(`Error reading project:`, err)
		os.Exit(1)
	}

	proj := abstractor.Abstract(ps, ao.Verbose)
	if err = writeJson(ao.OutPath, ao.Minimize, proj); err != nil {
		fmt.Println(`Error abstracting project:`, err)
		os.Exit(1)
	}

	os.Exit(0)
}

func writeJson(path string, minimize bool, p *constructs.Project) error {
	ctx := jsonify.NewContext()
	ctx.Set(`minimize`, minimize)
	b := jsonify.Marshal(ctx, p)

	if len(path) > 0 {
		return os.WriteFile(path, b, 0666)
	}

	_, err := fmt.Println(string(b))
	return err
}
