package compiler

import (
	"io/ioutil"
	"strings"
)

// A SourceFile is a representation of a single source file in any of the supported languages
type SourceFile struct {
	code            string
	dependencies    []SourceFile
	exportedSymbols []string
	importedSymbols []string
}

// newSourceFile creates a new sourcefile from the filepath
func newSourceFile(filepath string) *SourceFile {
	// get the code
	codeBytes, err := ioutil.ReadFile(filepath)
	checkError(err)
	code := string(codeBytes)

	return &SourceFile{
		code:            code,
		dependencies:    make([]SourceFile, 0),
		exportedSymbols: make([]string, 0),
		importedSymbols: make([]string, 0),
	}
}

// lines returns an array comprised of all the lines in the sourcefile's code
func (s *SourceFile) lines() []string {
	return strings.Split(s.code, "\n")
}

// addImportedSymbols appends the given symbols to the sourcefile's slice of imported symbols
func (s *SourceFile) addImportedSymbols(symbols []string) {
	s.importedSymbols = append(s.importedSymbols, symbols...)
}

// addExportedSymbols appends the given symbols to the sourcefile's slice of exported symbols
func (s *SourceFile) addExportedSymbols(symbols []string) {
	s.exportedSymbols = append(s.exportedSymbols, symbols...)
}
