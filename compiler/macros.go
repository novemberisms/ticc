package compiler

import "errors"

// MacroType is an enum that specifies one of the available macro types supported by the ticc code compiler
type MacroType int
type conditionalMode int

const (
	conditionalNormalExecution conditionalMode = iota
	conditionalExecuteBlock
	conditionalWaitForElseIf
	conditionalWaitForEnd
)

const (
	// MacroTypeUnknown denotes a macro with an invalid name. Detecting this should raise errors.
	MacroTypeUnknown MacroType = iota
	// MacroTypeDefine denotes a #define macro. This is used just like in C.
	MacroTypeDefine
	// MacroTypeString denotes a #string macro, which is similar to a #define macro, but it takes
	// in the contents of a string and preserves punctuation and spacing.
	MacroTypeString
	// MacroTypeIf denotes the start of a conditional compilation block
	MacroTypeIf
	// MacroTypeElseIf denotes an else-if block in a conditional compilation block
	MacroTypeElseIf
	// MacroTypeElse denotes an else block in a conditional compilation block
	MacroTypeElse
	// MacroTypeEndIf denotes an endif marker that terminates the conditional compilation mode
	MacroTypeEndIf
)

func (c *Compiler) handleMacro(macroType MacroType, line string) error {

	if !c._shouldHandleMacro(macroType, line) {
		return nil
	}

	switch macroType {
	case MacroTypeDefine:
		args := c.LangService.GetMacroArgs(line)
		if len(args) == 0 {
			return errors.New("define macro must have at least 1 argument")
		} else if len(args) == 1 {
			c._newDefine(args[0], "true")
		} else if len(args) == 2 {
			c._newDefine(args[0], args[1])
		} else {
			return errors.New("too many arguments to define macro. did you mean to use a string macro?")
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

func (c *Compiler) _getConditionalMode() conditionalMode {
	if c.conditionStack.Len() == 0 {
		return conditionalNormalExecution
	}

	return c.conditionStack.Peek().(conditionalMode)
}

func (c *Compiler) _shouldHandleMacro(macroType MacroType, line string) bool {

	// if we're in a conditional block that's not supposed to be executed,
	// only process macros if they're an ELSEIF, ELSE, or ENDIF
	// (new IF blocks are not processed if it's within a block that should not be processed)
	switch c._getConditionalMode() {
	case conditionalWaitForElseIf:
		fallthrough
	case conditionalWaitForEnd:
		switch macroType {
		case MacroTypeIf:
			c.disabledNestedIfCount++
		case MacroTypeEndIf:
			c.disabledNestedIfCount--

		case MacroTypeElseIf:
			fallthrough
		case MacroTypeElse:

		}

		return false
	}

	return true
}
