package main

import (
	"errors"
	"fmt"
	"time"

	"github.com/radovskyb/watcher"

	"github.com/novemberisms/ticc/compiler"
	"github.com/novemberisms/ticc/moonlang"
)

func main() {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println(r)
		}
	}()

	// populate the Args global var with the proper command line args
	getArguments()

	fmt.Printf("===================TICC=====================\n")
	fmt.Printf("language: %s\n", Args.language)
	fmt.Printf("dir: %s\n", Args.directory.Name())
	fmt.Printf("out: %s\n", Args.outputFile)
	fmt.Printf("============================================\n")

	doCompilation()

	if !Args.watchMode {
		return
	}

	w := watcher.New()
	defer w.Close()

	w.FilterOps(watcher.Rename, watcher.Remove, watcher.Write, watcher.Move, watcher.Create)
	w.SetMaxEvents(1)

	if err := w.AddRecursive(Args.directory.Name()); err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("--------------------------------------------\n")
	fmt.Printf("Starting to watch directory '%s'...\n", Args.directory.Name())
	fmt.Printf("--------------------------------------------\n")

	go func() {
		for {
			select {
			case <-w.Event:
				doCompilation()
				fmt.Printf("--------------------------------------------\n")
			case err := <-w.Error:
				fmt.Println(err)
			case <-w.Closed:
				return
			}
		}
	}()

	if err := w.Start(time.Millisecond * 500); err != nil {
		fmt.Println(err)
	}
}

func doCompilation() {

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

	fmt.Println("Compiling...")

	err = comp.Start()

	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println("OK")
	}

}
