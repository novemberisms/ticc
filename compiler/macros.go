package compiler

import "errors"

// MacroType is an enum that specifies one of the available macro types supported by the ticc code compiler
type MacroType string

const (
	// MacroTypeUnknown denotes a macro with an invalid name. Detecting this should raise errors.
	MacroTypeUnknown MacroType = "unknown"
	// MacroTypeDefine denotes a #define macro. This is used just like in C.
	MacroTypeDefine MacroType = "define"
	// MacroTypeString denotes a #string macro, which is similar to a #define macro, but it takes
	// in the contents of a string and preserves punctuation and spacing.
	MacroTypeString MacroType = "string"
	// MacroTypeIf denotes the start of a conditional compilation block
	MacroTypeIf MacroType = "if"
	// MacroTypeElseIf denotes an else-if block in a conditional compilation block
	MacroTypeElseIf MacroType = "elseif"
	// MacroTypeElse denotes an else block in a conditional compilation block
	MacroTypeElse MacroType = "else"
)

func (c *Compiler) handleMacro(macroType MacroType, line string) error {
	switch macroType {
	case MacroTypeDefine:
		args := c.LangService.GetMacroArgs(line)
		if len(args) == 0 {
			return errors.New("define macro must have at least 1 argument")
		}
		if len(args) == 1 {
			c._newDefine(args[0], "true")
		} else {
			c._newDefine(args[0], args[1])
		}
		return nil
	case MacroTypeString:
		name, contents, err := c.LangService.GetMacroStringDeclaration(line)
		if err != nil {
			return err
		}
		c._newDefine(name, contents)
		return nil
	}
	return errors.New("unknown macro")
}

func (c *Compiler) _newDefine(from string, to string) {
	c.defines[from] = to
}
