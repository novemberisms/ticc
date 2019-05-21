package compiler

import "os"

// A Parser is one of the language implementations that will read a line of code and determine certain
// qualities about it
type Parser interface {
	IsLineImport(line string) bool
}

// Compiler is the central control struct that reads input files and stitches them together into the output file
type Compiler struct {
	Parser
	MainFile   *os.File
	outputFile *os.File
	directory  string
}

// NewCompiler creates a new compiler with the given parameters
func NewCompiler(parser Parser, mainfile *os.File, outputfile *os.File, directory string) *Compiler {
	return &Compiler{
		parser,
		mainfile,
		outputfile,
		directory,
	}
}

// Start starts the compilation process
func (c Compiler) Start() {

}

func (c Compiler) processFile(file *os.File) {

}
