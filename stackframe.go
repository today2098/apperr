package apperr

import "runtime"

// StackFrame contains information for stacktrace.
type StackFrame struct {
	File           string  // The file name
	Line           int     // The line number
	Function       string  // The function name
	ProgramCounter uintptr // The program counter
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
