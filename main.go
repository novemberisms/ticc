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

	comp := compiler.NewCompiler(
		wrenparser.WrenParser{},
		mainFile,
		Args.outputFile,
		Args.directory.Name(),
	)

	fmt.Print(comp)

	for _, line := range lines(mainFile) {
		outline := fmt.Sprintf("[%s}\n", line)
		Args.outputFile.WriteString(outline)
	}
}
