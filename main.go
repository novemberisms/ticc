package main

import (
	"fmt"
	"os"

	"github.com/novemberisms/ticc/compiler"
	"github.com/novemberisms/ticc/wrenparser"
)

func main() {
	getArguments()
	// this is deferred from here so even if the program panics, the file will still be closed
	defer Args.outputFile.Close()

	fmt.Printf("language: %s\n", Args.language)
	fmt.Printf("dir: %s\n", Args.directory.Name())

	mainFileName, err := findMainFile(Args.directory.Name())
	checkError(err)
	mainFile, err := os.Open(mainFileName)
	checkError(err)
	defer mainFile.Close()

	var parser compiler.Parser

	switch Args.language {
	case wren:
		parser = wrenparser.WrenParser{}
	default:
		parser = wrenparser.WrenParser{}
	}

	comp := compiler.NewCompiler(
		parser,
		mainFile,
		Args.outputFile,
		Args.directory.Name(),
	)

	comp.Start()
}
