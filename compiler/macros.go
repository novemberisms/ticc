package compiler

import "errors"

type MacroType string

const (
	MacroTypeUnknown MacroType = "unknown"
	MacroTypeDefine  MacroType = "define"
	MacroTypeIf      MacroType = "if"
	MacroTypeElseIf  MacroType = "elseif"
	MacroTypeElse    MacroType = "else"
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
	}

	return errors.New("unknown macro")
}

func (c *Compiler) _newDefine(from string, to string) {
	c.defines[from] = to
}
