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

// NewSourceFile creates a new sourcefile from the filepath
func NewSourceFile(filepath string) *SourceFile {
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

// Lines returns an array comprised of all the lines in the sourcefile's code
func (s *SourceFile) Lines() []string {
	return strings.Split(s.code, "\n")
}
