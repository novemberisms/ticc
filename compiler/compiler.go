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
	alreadyImportedFiles map[string]*SourceFile
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
		make(map[string]*SourceFile),
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
	c.alreadyImportedFiles[sourcefile.path] = sourcefile
}

func (c *Compiler) _popFile() *SourceFile {
	file := c.fileStack.Pop()
	return file
}

// TODO: make use of this and caching of the already imported files to crosscheck exports and imports to make sure they match
// when importing a file
func (c Compiler) _shouldProcessNewFile(path string) bool {
	_, exists := c.alreadyImportedFiles[path]
	return !exists
}

func (c *Compiler) _getCachedFile(path string) *SourceFile {
	return c.alreadyImportedFiles[path]
}

func (c *Compiler) _processFile() error {
	currentFile := c.fileStack.Peek()

	for i, line := range currentFile.lines() {
		lineNumber := i + 1
		langService := c.LangService

		// TODO: make an exception for the prelude in the main file
		leanLine := langService.StripUnimportant(line)

		// completely empty strings should be ignored
		if len(strings.TrimSpace(leanLine)) == 0 {
			continue
		}

		if langService.IsLineImport(leanLine) {
			err := c._handleImport(leanLine, lineNumber, currentFile)
			if err != nil {
				return fmt.Errorf("Error processing file '%s' (line %d):\n%w", currentFile.path, lineNumber, err)
			}
			continue
		}

		// if control reaches here, then it the current line is
		// just a normal line that should be copied into the output

		if langService.IsExportDeclaration(leanLine) {
			exportedSymbols := langService.GetExportDeclarations(leanLine)
			currentFile.addExportedSymbols(exportedSymbols)
		}

		c._writeLine(leanLine)
	}

	return nil
}

func (c *Compiler) _handleImport(line string, lineNumber int, currentFile *SourceFile) error {
	importedSymbols, partialRequirePath, err := c.LangService.GetImportData(line)

	if err != nil {
		return err
	}

	currentFile.addImportedSymbols(importedSymbols)

	requirePath := path.Join(c.directory, partialRequirePath)

	if c._shouldProcessNewFile(requirePath) {
		requiredFile, err := newSourceFile(requirePath)
		if err != nil {
			return fmt.Errorf("Error trying to import file '%s':\n%w", requirePath, err)
		}

		c._pushFile(requiredFile)
		err = c._processFile()
		if err != nil {
			return err
		}
		poppedRequiredFile := c._popFile()

		err = _validateImportExportSymbols(importedSymbols, poppedRequiredFile)
		if err != nil {
			return fmt.Errorf("Error trying to import file '%s':\n%w", requirePath, err)
		}
	} else {
		requiredFile := c._getCachedFile(requirePath)
		err := _validateImportExportSymbols(importedSymbols, requiredFile)
		if err != nil {
			return fmt.Errorf("Error trying to import file '%s':\n%w", requirePath, err)
		}
	}

	return nil
}

func _validateImportExportSymbols(importedSymbols []string, requiredFile *SourceFile) error {

	for _, symbol := range importedSymbols {
		found := false
		// find the symbol in the file's exported symbols
		for _, exported := range requiredFile.exportedSymbols {
			if symbol == exported {
				found = true
				break
			}
		}

		if !found {
			return fmt.Errorf("file does not export symbol '%s'", symbol)
		}
	}
	return nil
}
