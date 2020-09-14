package main

import (
	"errors"
	"fmt"
	"time"

	"github.com/radovskyb/watcher"

	"github.com/novemberisms/ticc/compiler"
	"github.com/novemberisms/ticc/moonlang"
	"github.com/novemberisms/ticc/wrenlang"
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
			case event := <-w.Event:
				// if ticc is being run inside the watched directory itself, and the output file is being
				// written to the same directory, then this will also pick up the output file being written
				// and cause a loop. So we ignore the output file here
				if event.Name() != Args.outputFile {
					doCompilation()
					fmt.Printf("--------------------------------------------\n")
				}
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
	case wren:
		langService = wrenlang.WrenLanguageService{}
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
