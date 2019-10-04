package compiler

import (
	"log"
	"os"
)

// A LangService is one of the language implementations that will read a line of code and determine certain
// qualities about it based on unique properties of the language itself
type LangService interface {
	// whether a given line is an import directive in the language
	IsLineImport(line string) bool
	// remove all unnecessary whitespace and any comments
	StripUnimportant(line string) string
	// given a line, fetch a list of imported symbols and the path of the file to import them from
	GetImportData(line string) ([]string, string)
}

// Compiler is the central control struct that reads input files and stitches them together into the output file
type Compiler struct {
	LangService
	outputFile *os.File
	directory  string
	fileStack  *FileStack
}

// NewCompiler creates a new compiler with the given parameters
func NewCompiler(langservice LangService, mainfile string, outputfile *os.File, directory string) *Compiler {
	mainSourceFile := newSourceFile(mainfile)

	fileStack := NewFileStack(1)
	fileStack.Push(mainSourceFile)

	return &Compiler{
		langservice,
		outputfile,
		directory,
		fileStack,
	}
}

// Start starts the compilation process
func (c *Compiler) Start() {
	c.processFile()
}

func (c Compiler) write(values ...string) {
	for _, s := range values {
		c.outputFile.WriteString(s)
	}
}

func (c Compiler) writeLine(lines ...string) {
	for _, line := range lines {
		c.outputFile.WriteString(line + "\n")
	}
}

func (c *Compiler) processFile() {
	currentFile := c.fileStack.Peek()
	for _, line := range currentFile.lines() {
		langService := c.LangService
		leanLine := langService.StripUnimportant(line)
		if langService.IsLineImport(leanLine) {
			importedSymbols, requirePath := langService.GetImportData(leanLine)
			requiredFile := newSourceFile(requirePath)
			c.fileStack.Push(requiredFile)

		} else {
			c.writeLine(leanLine)
		}
	}
}

func checkError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
