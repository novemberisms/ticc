package compiler

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/novemberisms/stack"
)

// A Parser is one of the language implementations that will read a line of code and determine certain
// qualities about it based on unique properties of the language itself
type Parser interface {
	IsLineImport(line string) bool
}

// Compiler is the central control struct that reads input files and stitches them together into the output file
type Compiler struct {
	Parser
	outputFile *os.File
	directory  string
	fileStack  *stack.Stack
}

// NewCompiler creates a new compiler with the given parameters
func NewCompiler(parser Parser, mainfile string, outputfile *os.File, directory string) *Compiler {
	mainSourceFile := newSourceFile(mainfile)

	fileStack := stack.NewStack(0)
	fileStack.Push(mainSourceFile)

	return &Compiler{
		parser,
		outputfile,
		directory,
		fileStack,
	}
}

// Start starts the compilation process
func (c *Compiler) Start() {
	main := c.peekSourceFile()
	fmt.Println(main.code)
	c.writeLine(strings.Split(main.code, "\n")...)
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

func (c Compiler) peekSourceFile() *SourceFile {
	return c.fileStack.Peek().(*SourceFile)
}

func checkError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
