package apperr

import "runtime"

type StackFrame struct {
	File           string
	Line           int
	Function       string
	ProgramCounter uintptr
}

func newStackFrames(callers []uintptr) []StackFrame {
	res := make([]StackFrame, 0, len(callers))
	frames := runtime.CallersFrames(callers)
	for {
		frame, more := frames.Next()
		res = append(res, StackFrame{
			File:           frame.File,
			Line:           frame.Line,
			Function:       frame.Function,
			ProgramCounter: frame.PC,
		})
		if !more {
			break
		}
	}
	return res
}
