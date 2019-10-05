package moonlang

import (
	"errors"
	"regexp"
	"strings"
)

// MoonscriptLanguageService is a container struct that encapsulates a bunch of methods that
// take in a line of code and do text processing to see if the line matches certain properties
// based on the language this service provides.
type MoonscriptLanguageService struct {
}

// detects if a line of code contains a single line comment
var reSingleLineComment = regexp.MustCompile(`--.*`)

// knowing that a line contains an import statement, extracts a string containing comma-separated import
// symbols, as well as the relative import path to the file these symbols reside
var reImportExtract = regexp.MustCompile(`import\s+(.+)\s+from\s+require\s+"(\w+)"`)

// given the comma-separated import symbols string, break it down and find only the symbols within
var reExtractImportSymbols = regexp.MustCompile(`\b\w+\b`)

// StripUnimportant returns a new line which is the result of stripping all the unimportant or non-usable
// characters from it. This includes stripping away unneeded whitespace, comments, and any text that comes after comments
func (ls MoonscriptLanguageService) StripUnimportant(line string) string {
	withoutComments := reSingleLineComment.ReplaceAllString(line, "")
	return strings.TrimRight(withoutComments, " \t\n\r")
}

// IsLineImport Determines whether a line of code contains an import statement. In moonscript, this is the 'require' token.
func (ls MoonscriptLanguageService) IsLineImport(line string) bool {
	// a line is considered an attempt to import if it has the token "require"
	matched, _ := regexp.MatchString(`\brequire\b`, line)
	return matched
}

// GetImportData will extract a slice of imported symbols and the relative path to the file that is being imported given
// a line of code with an import statement
func (ls MoonscriptLanguageService) GetImportData(line string) ([]string, string, error) {
	matchInfo := reImportExtract.FindStringSubmatch(line)

	if len(matchInfo) != 3 {
		// this means the line does not match the template
		// import <symbols> from require "<path>"
		return nil, "", errors.New(`import line does not match the template: 'import {symbols} from require "{importpath}"'`)
	}

	importSymbols := reExtractImportSymbols.FindAllString(matchInfo[1], -1)
	requiredFilename := matchInfo[2] + ".moon"

	return importSymbols, requiredFilename, nil
}
