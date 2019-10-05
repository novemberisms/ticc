package main

import (
	"errors"
	"fmt"

	"github.com/novemberisms/ticc/compiler"
	"github.com/novemberisms/ticc/moonlang"
)

func main() {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println(r)
			return
		}
	}()

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

	// select a langserver based on the supplied language
	var langService compiler.LangService

	switch Args.language {
	case moon:
		langService = moonlang.MoonscriptLanguageService{}
	default:
		checkError(errors.New("language not yet implemented"))
	}

	// create the compiler struct
	comp := compiler.NewCompiler(
		langService,
		mainFile,
		Args.outputFile,
		Args.directory.Name(),
	)

	comp.Start()
}
