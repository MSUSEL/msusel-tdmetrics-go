package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/Snow-Gremlin/goToolbox/argers/args"

	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/abstractor"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/reader"
)

type argObject struct {
	ShowHelp bool `args:"h, help"`
	Verbose  bool `args:"v, verbose"`
	Minimize bool `args:"m, minimize"`
	InPath   string
	OutPath  string `args:"optional"`
}

func main() {
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("Panic: %v\n", r)
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
		fmt.Println(os.Args[0], `<options> inputPath [outputPath]`)
		fmt.Println(`  --help|-h: Shows this help text.`)
		fmt.Println(`  --verbose|-v: Indicates the abstraction process should`,
			`output additional status information.`)
		fmt.Println(`  --minimize|-m: Indicates the JSON output should be`,
			`minimized instead of formatted.`)
		fmt.Println(`  inputPath: The path to the directory of the project to read.`,
			`The directory should have a go.mod file.`)
		fmt.Println(`  outputPath (optional): The file path to write the JSON to.`,
			`If not given, the JSON will be outputted to the console.`)
		os.Exit(0)
	}

	ps, err := reader.Read(&reader.Config{
		Verbose: ao.Verbose,
		Path:    ao.InPath,
	})
	if err != nil {
		fmt.Println(`Error reading project:`, err)
		os.Exit(1)
	}

	proj := abstractor.Abstract(ps)
	if err = writeJson(ao.OutPath, ao.Minimize, proj); err != nil {
		fmt.Println(`Error abstracting project:`, err)
		os.Exit(1)
	}

	os.Exit(0)
}

func jsonMarshal(minimize bool, data any) ([]byte, error) {
	if minimize {
		return json.Marshal(data)
	}
	return json.MarshalIndent(data, ``, `  `)
}

func writeJson(path string, minimize bool, data any) error {
	b, err := jsonMarshal(minimize, data)
	if err != nil {
		return err
	}

	if len(path) > 0 {
		return os.WriteFile(path, b, 0666)
	}

	_, err = fmt.Println(string(b))
	return err
}
