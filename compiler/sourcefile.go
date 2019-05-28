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

func newSourceFile(filename string) *SourceFile {
	// get the code
	codeBytes, err := ioutil.ReadFile(filename)
	checkError(err)
	code := string(codeBytes)

	return &SourceFile{
		code: code,
	}
}

func (s *SourceFile) lines() []string {
	return strings.Split(s.code, "\n")
}
