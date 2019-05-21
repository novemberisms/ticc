package wrenparser

// WrenParser is the parser for the Wren programming language
type WrenParser struct{}

// IsLineImport tells whether the given line is an import directive in the target language
func (wp WrenParser) IsLineImport(line string) bool {
	return false
}
