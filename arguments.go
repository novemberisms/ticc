package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// Args holds all the optional arguments in the form of flags
var Args struct {
	language   Language
	directory  os.FileInfo
	positional []string
	outputFile string
	defines    map[string]string
	watchMode  bool
}

// Must be called as soon as the program starts to initialize Args
func getArguments() {
	// define pointers to the arguments which will be filled up when flag.Parse() is called
	langFlag := flag.String("l", string(auto), "Which language to use. Args are: lua | wren | moon | auto")
	dirFlag := flag.String("d", ".", "The directory containing the main file and the subfiles")
	outFlag := flag.String("o", "out", "The output file (sans extension)")
	watchFlag := flag.Bool("w", false, "Whether to enable Watch mode, which automatically recompiles if a file has changed in the directory")
	definesFlag := flag.String("D", "", "Used to pass in defines before compiling. Format is -D \"var1=value;var2=value;var3=value\"")

	// begin parsing the flags
	flag.Parse()

	// these setup functions have to be performed in this particular order
	// because they depend on certain fields of Args to be set when they are called
	_setDir(*dirFlag)
	_setLanguage(*langFlag)
	_setOutputFile(*outFlag)
	_setDefines(*definesFlag)

	Args.watchMode = *watchFlag

	// this gives all the non-flag command line args
	Args.positional = flag.Args()
}

func _setDir(dirname string) {
	// make sure it's actually a directory first
	stat, err := os.Stat(dirname)
	checkError(err)

	if !stat.IsDir() {
		checkError(errors.New("the argument to -d must be a directory"))
	}

	Args.directory = stat
}

func _setLanguage(rawInput string) {
	Args.language = Language(rawInput)

	if Args.language == auto {
		// automatically detect the language by finding a file called 'main' and checking its extension
		pathToMain, err := findMainFile(Args.directory.Name())
		checkError(err)
		ext := filepath.Ext(pathToMain)
		// trim out the dot in the beginning
		Args.language = Language(ext[1:])
	}

	if !isSupportedLanguage(Args.language) {
		checkError(fmt.Errorf("invalid language detected (%s) the supported languages are: lua | moon | wren", Args.language))
	}
}

func _setOutputFile(filename string) {
	// determine the name of the output file
	// if the file already has an extension, then use that
	ext := filepath.Ext(filename)
	if ext == "" {
		// auto fix
		filename += "." + string(Args.language)
	} else if ext != string(Args.language) {
		checkError(
			errors.New(
				`The output file must have the same extension as the detected language. 
				Alternatively, you may omit the extension and it will automatically be detected`,
			),
		)
	}
	// check to see if the output file already exists, and if so, delete it
	_deleteIfExists(filename)
	Args.outputFile = filename
}

func _setDefines(input string) {
	Args.defines = make(map[string]string)

	if len(input) == 0 {
		return
	}

	for _, define := range strings.Split(input, ";") {
		expressions := strings.Split(define, "=")

		lhs := strings.TrimSpace(expressions[0])

		var rhs string

		if len(expressions) == 1 {
			rhs = "true"
		} else {
			rhs = strings.TrimSpace(expressions[1])
		}

		Args.defines[lhs] = rhs
	}
}

func _deleteIfExists(filename string) {
	// does it exist?
	_, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return
	}
	// it does, sp delete it
	err = os.Remove(filename)
	checkError(err)
}
