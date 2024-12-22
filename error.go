package apperr

import (
	"errors"
	"fmt"
	"io"
	"runtime"
)

// The maximum depth of stackframes on any Error.
const maxStackDepth = 32

// Error represents an error with status code, body and stacktrace.
type Error struct {
	StatusCode int
	Body       Body
	prefix     string
	cause      error
	stack      []uintptr
	frames     []StackFrame
}

var (
	_ error         = (*Error)(nil)
	_ fmt.Formatter = (*Error)(nil)
)

// New creates a new Error from status code and body.
func New(code int, body Body) *Error {
	return &Error{
		StatusCode: code,
		Body:       body,
	}
}

// Error returns a message about status code, body and the wrapped error.
func (e *Error) Error() string {
	if e.cause == nil {
		return fmt.Sprintf("apperr(%d): %v", e.StatusCode, e.Body)
	}
	return fmt.Sprintf("apperr(%d): %v; %s%v", e.StatusCode, e.Body, e.prefix, e.cause)
}

// Unwrap returns the wrapped error.
func (e *Error) Unwrap() error {
	return e.cause
}

// Is reports whether e matches target.
func (e *Error) Is(target error) bool {
	var appErr *Error
	if errors.As(target, &appErr) && appErr.StatusCode == e.StatusCode {
		if e.Body == nil || appErr.Body == nil {
			return e.Body == appErr.Body
		}
		return e.Body.Is(appErr.Body)
	}
	return false
}

// Wrap creates a new Error that copies status code and body and wraps err.
func (e *Error) Wrap(err any) *Error {
	res := &Error{
		StatusCode: e.StatusCode,
		Body:       e.Body.Copy(),
	}

	switch err := err.(type) {
	case nil:
		res.prefix = ""
		res.cause = nil
	case *Error:
		res.prefix = err.prefix
		res.cause = err.cause
	case error:
		res.prefix = ""
		res.cause = err
	default:
		res.prefix = ""
		res.cause = fmt.Errorf("%v", err)
	}

	res.callers(1)

	return res
}

// WrapPrefix creates a new Error that copies status code and body and wraps err with prefix.
func (e *Error) WrapPrefix(err any, prefix string) *Error {
	res := &Error{
		StatusCode: e.StatusCode,
		Body:       e.Body.Copy(),
	}

	switch err := err.(type) {
	case nil:
		res.prefix = ""
		res.cause = nil
	case *Error:
		res.prefix = fmt.Sprintf("%s: %s", prefix, err.prefix)
		res.cause = err.cause
	case error:
		res.prefix = prefix + ": "
		res.cause = err
	default:
		res.prefix = prefix + ": "
		res.cause = fmt.Errorf("%v", err)
	}

	res.callers(1)

	return res
}

// StackFrames returns an array of StackFrame.
func (e *Error) StackFrames() []StackFrame {
	if e.frames == nil {
		e.frames = newStackFrames(e.stack)
	}
	return e.frames
}

// Format prints the detail about e (api for fmt package).
func (e *Error) Format(s fmt.State, verb rune) {
	switch verb {
	case 'v':
		if s.Flag('+') {
			io.WriteString(s, e.Error()) // TODO: Print stackframes.
			return
		}
		fallthrough
	case 's':
		io.WriteString(s, e.Error())
	case 'q':
		fmt.Fprintf(s, "%q", e.Error())
	}
}

func (e *Error) callers(skip int) {
	stack := make([]uintptr, maxStackDepth)
	len := runtime.Callers(2+skip, stack)
	e.stack = stack[:len]
	e.frames = nil
}