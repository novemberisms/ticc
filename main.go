package main

import (
	"fmt"
	"os"
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
	for _, line := range lines(mainFile) {
		outline := fmt.Sprintf("[%s}\n", line)
		Args.outputFile.WriteString(outline)
	}
}
