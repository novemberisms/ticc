package compiler

import (
	"fmt"
	"os"
	"path"
	"strings"
)

// A LangService is one of the language implementations that will read a line of code and determine certain
// qualities about it based on unique properties of the language itself
type LangService interface {
	// whether a given line is an import directive in the language
	IsLineImport(line string) bool
	// remove all unnecessary whitespace and any comments
	StripUnimportant(line string) string
	// given a line, fetch a list of imported symbols and the relative path of the file to import them from (with extension)
	GetImportData(line string) ([]string, string, error)
	// whether a given line declares a new export for the file
	IsExportDeclaration(line string) bool
	// fetch a list of exported symbols for a given line, assuming it really is an export declaration
	GetExportDeclarations(line string) []string
}

// Compiler is the central control struct that reads input files and stitches them together into the output file
type Compiler struct {
	LangService
	outputFile           *os.File
	directory            string
	fileStack            *FileStack
	alreadyImportedFiles map[string]SourceFile
}

// NewCompiler creates a new compiler with the given parameters
func NewCompiler(
	langservice LangService,
	mainfile string,
	outputfile *os.File,
	directory string,
) *Compiler {
	// the main file is guaranteed to exist
	mainSourceFile, _ := newSourceFile(mainfile)

	fileStack := NewFileStack(1)
	fileStack.Push(mainSourceFile)

	return &Compiler{
		langservice,
		outputfile,
		directory,
		fileStack,
		make(map[string]SourceFile),
	}
}

// Start starts the compilation process
func (c *Compiler) Start() {
	err := c._processFile()
	if err != nil {
		fmt.Println(err.Error())
	}
}

func (c Compiler) _write(values ...string) {
	for _, s := range values {
		c.outputFile.WriteString(s)
	}
}

func (c Compiler) _writeLine(lines ...string) {
	for _, line := range lines {
		c.outputFile.WriteString(line + "\n")
	}
}

func (c *Compiler) _pushFile(sourcefile *SourceFile) {
	c.fileStack.Push(sourcefile)
	c.alreadyImportedFiles[sourcefile.path] = *sourcefile
}

// TODO: make use of this and caching of the already imported files to crosscheck exports and imports to make sure they match
// when importing a file
func (c Compiler) _shouldProcessNewFile(path string) bool {
	_, exists := c.alreadyImportedFiles[path]
	return !exists
}

func (c *Compiler) _processFile() error {
	currentFile := c.fileStack.Peek()

	for lineNumber, line := range currentFile.lines() {
		langService := c.LangService

		// TODO: make an exception for the prelude in the main file
		leanLine := langService.StripUnimportant(line)

		// completely empty strings should be ignored
		if len(strings.TrimSpace(leanLine)) == 0 {
			continue
		}

		if langService.IsLineImport(leanLine) {

			importedSymbols, partialRequirePath, err := langService.GetImportData(leanLine)

			if err != nil {
				return fmt.Errorf(
					"Error processing file %s (line %d): %w",
					currentFile.path,
					lineNumber,
					err,
				)
			}

			currentFile.addImportedSymbols(importedSymbols)

			requirePath := path.Join(c.directory, partialRequirePath)

			if !c._shouldProcessNewFile(requirePath) {
				// skip files that were already imported earlier in the program
				continue
			}

			requiredFile, err := newSourceFile(requirePath)

			if err != nil {
				return fmt.Errorf("Error trying to import file '%s': %w", requirePath, err)
			}

			c.fileStack.Push(requiredFile)

			err = c._processFile()

			if err != nil {
				return err
			}

		} else if langService.IsExportDeclaration(leanLine) {

			exportedSymbols := langService.GetExportDeclarations(leanLine)

			currentFile.addExportedSymbols(exportedSymbols)

			c._writeLine(leanLine)

		} else {
			c._writeLine(leanLine)
		}
	}

	return nil
}
