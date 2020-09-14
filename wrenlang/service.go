package wrenlang

import (
	"errors"
	"regexp"
	"strings"

	"github.com/novemberisms/ticc/compiler"
)

type WrenLanguageService struct {
}

// finds the start of a single line comment
var reSingleLineComment = regexp.MustCompile(`\/\/.*`)

// matches `import "<path>" for <symbols>`
var reImportExtract = regexp.MustCompile(`import\s+"([\w\/]+)"\s+for\s+(.+)`)

// matches a bare import statement `import "path"`
var reBareImport = regexp.MustCompile(`import\s+"([\w\/]+)"\s*$`)

// given a comma-separated string of values, extracts each word
var reExtractImportSymbols = regexp.MustCompile(`\b\w+\b`)

var reExportClassDeclaration = regexp.MustCompile(`class\s+([A-Z]\w*)`)

var reExportVarDeclaration = regexp.MustCompile(`var\s+([A-Z]\w*)`)

var reIsPreludeComment = regexp.MustCompile(`^\/\/\s*\w+\s*:`)

var reGetMacroType = regexp.MustCompile(`\/\/#\s*(\w+)`)
var reGetMacroArgs = regexp.MustCompile(`\/\/#\s*\w+\s+(.*)$`)
var reBreakDownMacroArgs = regexp.MustCompile(`\S+`)

var reGetMacroStringDeclarationArgs = regexp.MustCompile(`\/\/#\s*\w+\s+(\w+)\s+(.*)$`)

var reIdentifiers = regexp.MustCompile(`\w+`)

func (ls WrenLanguageService) StripUnimportant(line string) string {
	withoutComments := reSingleLineComment.ReplaceAllString(line, "")
	trimmed := strings.TrimSpace(withoutComments)
	return trimmed
}

func (ls WrenLanguageService) IsLineImport(line string) bool {
	matched, _ := regexp.MatchString(`\bimport\b`, line)
	return matched
}

// GetImportData will extract a slice of imported symbols and the relative path to the file that is being imported given
// a line of code with an import statement
func (ls WrenLanguageService) GetImportData(line string) (compiler.ImportData, error) {
	matchInfo := reImportExtract.FindStringSubmatch(line)

	if len(matchInfo) != 3 {
		// this means the line does not match the template
		// import "<path>" for <symbols>

		// check if the line is a bare imperative import statement
		// import "<path>"
		if bareImport := reBareImport.FindStringSubmatch(line); len(bareImport) > 0 {
			// return []string{}, bareImport[1] + ".wren", nil
			return compiler.ImportData{
				Path: bareImport[1] + ".wren",
			}, nil
		}

		return compiler.ImportData{}, errors.New(`import line does not match the templates: 'import {path} for {symbols}' or 'import {path}'`)
	}

	importData := compiler.ImportData{
		Path:    matchInfo[1] + ".wren",
		Symbols: reExtractImportSymbols.FindAllString(matchInfo[2], -1),
	}

	return importData, nil
}

func (ls WrenLanguageService) IsExportDeclaration(line string) bool {
	// in wren, the following are export declarations:
	// var Dude
	// class Entity
	// but because we cannot tell if these are top-level declarations,
	// the rule in ticc is that if it starts with a capital letter, it will
	// be recognized as an export.

	if reExportClassDeclaration.MatchString(line) {
		return true
	}

	if reExportVarDeclaration.MatchString(line) {
		return true
	}

	return false
}

func (ls WrenLanguageService) GetExportDeclarations(line string) []string {
	matchInfo := reExportClassDeclaration.FindStringSubmatch(line)

	if len(matchInfo) != 0 {
		return matchInfo[1:2]
	}

	matchInfo = reExportVarDeclaration.FindStringSubmatch(line)

	if len(matchInfo) != 0 {
		return matchInfo[1:2]
	}

	return []string{}
}

func (ls WrenLanguageService) ExtractPrelude(mainFileCode string) string {
	result := ""

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

func (ls WrenLanguageService) SubstituteDefines(line string, defines map[string]string) string {
	return reIdentifiers.ReplaceAllStringFunc(line, func(identifier string) string {
		replacement, isDefined := defines[identifier]
		if !isDefined {
			return identifier
		}
		return replacement
	})
}

func (ls WrenLanguageService) IsLineMacro(line string) bool {
	matches, _ := regexp.MatchString(`^\s*\/\/#`, line)
	return matches
}

func (ls WrenLanguageService) GetMacroType(line string) compiler.MacroType {
	matchInfo := reGetMacroType.FindStringSubmatch(line)

	symbol := matchInfo[1]

	switch strings.ToUpper(symbol) {
	case "DEFINE":
		return compiler.MacroTypeDefine
	case "STRING":
		return compiler.MacroTypeString
	case "IF":
		return compiler.MacroTypeIf
	case "ELSEIF":
		return compiler.MacroTypeElseIf
	case "ELSE":
		return compiler.MacroTypeElse
	case "ENDIF":
		return compiler.MacroTypeEndIf
	default:
		return compiler.MacroTypeUnknown
	}
}

func (ls WrenLanguageService) GetMacroArgs(line string) []string {
	matchInfo := reGetMacroArgs.FindStringSubmatch(line)

	if len(matchInfo) < 2 {
		return []string{}
	}

	fullArgs := matchInfo[1]

	return reBreakDownMacroArgs.FindAllString(fullArgs, -1)
}

func (ls WrenLanguageService) GetMacroStringDeclaration(line string) (string, string, error) {
	matchInfo := reGetMacroStringDeclarationArgs.FindStringSubmatch(line)

	// the first string in matchInfo is always the full matched text
	// the second string is the string name
	// the third string is the string contents
	if len(matchInfo) != 3 {
		return "", "", errors.New("invalid format for string macro. must be //#string [STRING_NAME] STRING CONTENTS")
	}

	return matchInfo[1], matchInfo[2], nil
}
