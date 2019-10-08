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
	// extract the prelude from the main file
	ExtractPrelude(mainFileCode string) string

	IsLineMacro(line string) bool
	GetMacroType(line string) MacroType
	GetMacroArgs(line string) []string
	GetMacroStringDeclaration(line string) (string, string, error)

	SubstituteDefines(line string, defines map[string]string) string
}

// Compiler is the central control struct that reads input files and stitches them together into the output file
type Compiler struct {
	LangService
	outputFile           *os.File
	directory            string
	fileStack            *FileStack
	alreadyImportedFiles map[string]*SourceFile
	defines              map[string]string
}

// NewCompiler creates a new compiler with the given parameters
func NewCompiler(
	langservice LangService,
	mainfile string,
	outputfilename string,
	directory string,
) *Compiler {
	// the main file is guaranteed to exist
	mainSourceFile, _ := newSourceFile(mainfile)

	fileStack := NewFileStack(1)
	fileStack.Push(mainSourceFile)

	outputFile, _ := os.Create(outputfilename)

	return &Compiler{
		langservice,
		outputFile,
		directory,
		fileStack,
		make(map[string]*SourceFile),
		make(map[string]string),
	}
}

// Start starts the compilation process
func (c *Compiler) Start() error {
	defer func() {
		c.outputFile.Close()
		c.outputFile = nil
	}()
	if err := c._writePrelude(); err != nil {
		return err
	}
	if err := c._processFile(); err != nil {
		return err
	}
	return nil
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

func (c *Compiler) _writePrelude() error {
	mainFile := c.fileStack.Peek()
	prelude := c.LangService.ExtractPrelude(mainFile.code)

	c._write(prelude)
	return nil
}

func (c *Compiler) _processFile() error {
	currentFile := c.fileStack.Peek()

	for i, line := range currentFile.lines() {
		lineNumber := i + 1
		langService := c.LangService

		// check if the line is a macro
		if langService.IsLineMacro(line) {
			macroType := langService.GetMacroType(line)
			err := c.handleMacro(macroType, line)
			if err != nil {
				return fmt.Errorf("Error processing file '%s' (line %d):\n%w", currentFile.path, lineNumber, err)
			}
			continue
		}

		line = langService.StripUnimportant(line)
		// NOTE that it is important that we substitute the defines AFTER we strip the unimportant spaces.
		// This allows us to easily preserve any player-facing strings from mangling by putting them in a
		// #string define
		line = langService.SubstituteDefines(line, c.defines)

		// completely empty strings should be ignored
		if len(strings.TrimSpace(line)) == 0 {
			continue
		}

		if langService.IsLineImport(line) {
			err := c._handleImport(line, lineNumber, currentFile)
			if err != nil {
				return fmt.Errorf("Error processing file '%s' (line %d):\n%w", currentFile.path, lineNumber, err)
			}
			continue
		}

		// if control reaches here, then it the current line is
		// just a normal line that should be copied into the output

		if langService.IsExportDeclaration(line) {
			exportedSymbols := langService.GetExportDeclarations(line)
			currentFile.addExportedSymbols(exportedSymbols)
		}

		c._writeLine(line)
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
