package compiler

import (
	"github.com/novemberisms/stack"
)

// FileStack is a stack of SourceFiles
type FileStack struct {
	stack.Stack
}

// Push pushes a SourceFile onto the stack
func (f *FileStack) Push(s *SourceFile) {
	f.Stack.Push(s)
}

// Pop pops a SourceFile from the stack
func (f *FileStack) Pop() *SourceFile {
	val := f.Stack.Pop()
	if val == nil {
		return nil
	}
	return val.(*SourceFile)
}

// Peek gets the top SourceFile
func (f FileStack) Peek() *SourceFile {
	val := f.Stack.Peek()
	if val == nil {
		return nil
	}
	return val.(*SourceFile)
}

// NewFileStack creates a new FileStack
func NewFileStack(cap int) *FileStack {
	return &FileStack{*stack.NewStack(cap)}
}
