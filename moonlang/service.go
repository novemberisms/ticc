package moonlang

import (
	"errors"
	"regexp"
	"strings"

	"github.com/novemberisms/ticc/compiler"
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

var reRequireFile = regexp.MustCompile(`require\s*"(\w+)"`)

// given the comma-separated import symbols string, break it down and find only the symbols within
var reExtractImportSymbols = regexp.MustCompile(`\b\w+\b`)

// extracts an exported symbol of kind [identifier] = ...
var reExtractExportedSymbols = regexp.MustCompile(`^(\w+)\s*=`)

// extracts an exported symbol of kind class [identifier]
var reExtractExportedClass = regexp.MustCompile(`^class\s+(\w+)`)

// determines if the given line matches the structure needed to be a prelude comment
var reIsPreludeComment = regexp.MustCompile(`^--\s*\w+\s*:`)

var reGetMacroType = regexp.MustCompile(`--#\s*(\w+)`)
var reGetMacroArgs = regexp.MustCompile(`--#\s*\w+\s+(.*)$`)
var reBreakDownMacroArgs = regexp.MustCompile(`\S+`)

var reIdentifiers = regexp.MustCompile(`[\w_]+`)

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

	// check if the line is a bare require without any imports (like 'require "defines"')
	if requiredFile := reRequireFile.FindStringSubmatch(line); len(requiredFile) > 0 {
		return []string{}, requiredFile[1] + ".moon", nil
	}

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

// IsExportDeclaration determines if a line contains a global declaration that should be available to other files
// importing this one
func (ls MoonscriptLanguageService) IsExportDeclaration(line string) bool {
	// in moonscript, the following are export declarations
	// * class Entity
	// * Entity = ...
	// and they must all have zero leading indentation

	// if there is any leading space, then it isn't a top-level statement and cannot be an export
	if match, _ := regexp.MatchString(line, `^\s`); match == true {
		return false
	}

	if reExtractExportedSymbols.MatchString(line) {
		return true
	}

	if reExtractExportedClass.MatchString(line) {
		return true
	}

	return false
}

// GetExportDeclarations extracts a list of exported symbols from the line
func (ls MoonscriptLanguageService) GetExportDeclarations(line string) []string {

	matchInfo := reExtractExportedSymbols.FindStringSubmatch(line)

	if len(matchInfo) != 0 {
		return matchInfo[1:2]
	}

	matchInfo = reExtractExportedClass.FindStringSubmatch(line)

	if len(matchInfo) != 0 {
		return matchInfo[1:2]
	}

	return []string{}
}

func (ls MoonscriptLanguageService) ExtractPrelude(mainFileCode string) string {
	result := ""
	// split the code into lines
	lines := strings.Split(mainFileCode, "\n")
	for _, line := range lines {
		if reIsPreludeComment.MatchString(line) {
			result = result + line + "\n"
		} else {
			break
		}
	}
	return result
}

func (ls MoonscriptLanguageService) SubstituteDefines(line string, defines map[string]string) string {
	return reIdentifiers.ReplaceAllStringFunc(line, func(identifier string) string {
		replacement, isDefined := defines[identifier]
		if !isDefined {
			return identifier
		}
		return replacement
	})
}

func (ls MoonscriptLanguageService) IsLineMacro(line string) bool {
	matches, _ := regexp.MatchString(`^\s*--#`, line)
	return matches
}

func (ls MoonscriptLanguageService) GetMacroType(line string) compiler.MacroType {
	matchInfo := reGetMacroType.FindStringSubmatch(line)

	symbol := matchInfo[1]

	switch strings.ToUpper(symbol) {
	case "DEFINE":
		return compiler.MacroTypeDefine
	case "IF":
		return compiler.MacroTypeIf
	case "ELSEIF":
		return compiler.MacroTypeElseIf
	case "ELSE":
		return compiler.MacroTypeElse
	default:
		return compiler.MacroTypeUnknown
	}
}

func (ls MoonscriptLanguageService) GetMacroArgs(line string) []string {
	matchInfo := reGetMacroArgs.FindStringSubmatch(line)

	if len(matchInfo) < 2 {
		return []string{}
	}

	fullArgs := matchInfo[1]

	return reBreakDownMacroArgs.FindAllString(fullArgs, -1)
}
