package wrenparser

import "strings"

// WrenParser is the parser for the Wren programming language
type WrenParser struct{}

// IsLineImport tells whether the given line is an import directive in the target language
func (wp WrenParser) IsLineImport(line string) bool {
	return false
}

// StripUnimportant will remove all trailing and leading whitespace and erase all comments
func (wp WrenParser) StripUnimportant(line string) string {
	noSpace := strings.TrimSpace(line)

	return noSpace
}
