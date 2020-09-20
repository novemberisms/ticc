package compiler

import (
	"errors"
	"fmt"
)

// MacroType is an enum that specifies one of the available macro types supported by the ticc code compiler
type MacroType int
type conditionalMode int

const (
	// conditionalNormalExecution means we are not in an if block of any kind
	conditionalNormalExecution conditionalMode = iota
	// conditionalExecuteBlock means we are in an if/elseif/else block that should be executed
	conditionalExecuteBlock
	// conditionalWaitForElseIf means we are in an if block and are waiting for an elseif/else/end block
	// because the condition in the previous block did not evaluate to true
	conditionalWaitForElseIf
	// conditionalWaitForEnd means we are in an if block and are now waiting for the end block because we've
	// already executed one of the blocks
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

func (m MacroType) String() string {
	switch m {
	case MacroTypeDefine:
		return "define"
	case MacroTypeString:
		return "string"
	case MacroTypeIf:
		return "if"
	case MacroTypeElseIf:
		return "elseif"
	case MacroTypeElse:
		return "else"
	case MacroTypeEndIf:
		return "endif"
	case MacroTypeUnknown:
		fallthrough
	default:
		return "unknown"
	}
}

func (c *Compiler) handleMacro(macroType MacroType, line string) error {

	if !c._shouldHandleMacro(macroType, line) {
		fmt.Printf("skipping macro %v\n", macroType.String())
		return nil
	}
	fmt.Printf("handling macro %v\n", macroType.String())

	switch macroType {
	case MacroTypeDefine:
		return c._handleDefineMacro(line)
	case MacroTypeString:
		return c._handleStringMacro(line)
	case MacroTypeIf:
		return c._handleIfMacro(line)
	case MacroTypeElseIf:
		return c._handleElseIfMacro(line)
	case MacroTypeElse:
		return c._handleElseMacro(line)
	case MacroTypeEndIf:
		return c._handleEndIfMacro(line)
	}
	return errors.New("unknown macro")
}

func (c *Compiler) shouldProcessLine(line string) bool {
	switch c._getConditionalMode() {
	case conditionalNormalExecution:
		return true
	case conditionalExecuteBlock:
		return true
	case conditionalWaitForElseIf:
		return false
	case conditionalWaitForEnd:
		return false
	}
	return true
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
	case conditionalNormalExecution:
		return true

	case conditionalExecuteBlock:
		return true

	case conditionalWaitForElseIf:
		fallthrough
	case conditionalWaitForEnd:
		switch macroType {
		case MacroTypeIf:
			c.disabledNestedIfCount++
			return false

		case MacroTypeEndIf:
			if c.disabledNestedIfCount > 0 {
				c.disabledNestedIfCount--
				return false
			}
			return true

		case MacroTypeElseIf:
			fallthrough
		case MacroTypeElse:
			return c.disabledNestedIfCount == 0
		}

		return false
	}

	return true
}

func (c *Compiler) _handleDefineMacro(line string) error {
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
}

func (c *Compiler) _handleStringMacro(line string) error {
	name, contents, err := c.LangService.GetMacroStringDeclaration(line)

	if err != nil {
		return err
	}

	c._newDefine(name, contents)

	return nil
}

func (c *Compiler) _handleIfMacro(line string) error {

	condition, err := c._evaluateConditional(line)

	if err != nil {
		return err
	}

	if condition {
		fmt.Println("conditional evaluated to true")
		c.conditionStack.Push(conditionalExecuteBlock)
	} else {
		fmt.Println("conditional evaluated to false")
		c.conditionStack.Push(conditionalWaitForElseIf)
	}

	return nil
}

func (c *Compiler) _handleEndIfMacro(line string) error {
	if c.conditionStack.Len() == 0 {
		return errors.New("found ENDIF macro with no prior IF")
	}

	c.conditionStack.Pop()

	return nil
}

func (c *Compiler) _handleElseIfMacro(line string) error {
	switch c._getConditionalMode() {
	case conditionalNormalExecution:
		return errors.New("found ELSEIF macro with no prior IF")
	case conditionalExecuteBlock:
		// the previous block was already executed, so do not evaluate or execute this one or any
		// subsequent ones until ENDIF is encountered
		c.conditionStack.Pop()
		c.conditionStack.Push(conditionalWaitForEnd)
	case conditionalWaitForElseIf:
		// the previous block did not get executed, so try to see if we can evaluate this one
		// this now just behaves like a normal if macro
		c.conditionStack.Pop()
		return c._handleIfMacro(line)
	case conditionalWaitForEnd:
		// one of the previous blocks was already executed, so do not evaluate this one
		return nil
	}
	return nil
}

func (c *Compiler) _handleElseMacro(line string) error {
	switch c._getConditionalMode() {
	case conditionalNormalExecution:
		return errors.New("found ELSE macro with no prior IF")
	case conditionalExecuteBlock:
		// the previous block was already executed, so do not evaluate the ELSE
		c.conditionStack.Pop()
		c.conditionStack.Push(conditionalWaitForEnd)
	case conditionalWaitForElseIf:
		// no previous blocks have been executed, so do execute the ELSE
		c.conditionStack.Pop()
		c.conditionStack.Push(conditionalExecuteBlock)
	case conditionalWaitForEnd:
		// one of the previous blocks was already executed, so do not evaluate this one
		return nil
	}
	return nil
}

func (c *Compiler) _evaluateConditional(line string) (bool, error) {
	args := c.LangService.GetMacroArgs(line)

	if len(args) == 1 {
		// #IF SOMEVALUE
		// evaluate as long as SOMEVALUE is not a falsy value
		// falsy values are "false", 0, or if SOMEVALUE has not been defined

		expression := args[0]

		if expression == "false" || expression == "0" {
			return false, nil
		}

		definedValue, exists := c.defines[expression]

		if !exists {
			return false, nil
		}

		if definedValue == "false" || definedValue == "0" {
			return false, nil
		}

		return true, nil

	}

	if len(args) == 3 {
		// #IF A == B
		// #IF A != B

		evaluateExpression := func(expression string) string {
			if val, defined := c.defines[expression]; defined {
				return val
			}

			return expression
		}

		lhs := evaluateExpression(args[0])
		op := args[1]
		rhs := evaluateExpression(args[2])

		if op == "==" {
			return lhs == rhs, nil
		} else if op == "!=" {
			return lhs != rhs, nil
		}

	}

	return false, fmt.Errorf("cannot parse condition in macro: %s", line)
}
