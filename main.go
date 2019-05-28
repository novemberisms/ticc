package main

import (
	"fmt"

	"github.com/novemberisms/ticc/compiler"
	"github.com/novemberisms/ticc/wrenparser"
)

func main() {
	// populate the Args global var with the proper command line args
	getArguments()
	defer Args.outputFile.Close()

	fmt.Printf("===================TICC=====================\n")
	fmt.Printf("language: %s\n", Args.language)
	fmt.Printf("dir: %s\n", Args.directory.Name())
	fmt.Printf("out: %s\n", Args.outputFile.Name())
	fmt.Printf("============================================\n")

	// get the name of the main file so we can pass it into the compiler
	mainFile, err := findMainFile(Args.directory.Name())
	checkError(err)

	// select a parser based on the supplied language
	var parser compiler.Parser

	switch Args.language {
	case wren:
		parser = wrenparser.WrenParser{}
	default:
		parser = wrenparser.WrenParser{}
	}

	// create the compiler struct
	comp := compiler.NewCompiler(
		parser,
		mainFile,
		Args.outputFile,
		Args.directory.Name(),
	)

	comp.Start()
}
