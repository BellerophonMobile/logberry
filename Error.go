package logberry

import (
	"bytes"
	"fmt"
	"runtime"
)

type Position struct {
	File string
	Line int
}

// Error captures structured information about a fault.
type Error struct {

	// An optional identifier for differentiating classes of errors
	Code string

	// Human-oriented description of the fault
	Message string

	// Inputs, parameters, and other data associated with the fault
	Data EventData

	// The source code file and line number where the error occurred
	Source Position

	// Optional link to a preceding error underlying the fault
	Cause error `logberry:"quiet"`

	// Whether or not this error has already been reported
	Reported bool `logberry:"quiet"`

}

func newerror(msg string, data []interface{}) *Error {
	e := &Error{
		Message: msg,
		Data: Copy(data),
	}
	return e
}

func wraperror(msg string, err error, data []interface{}) *Error {
	e := newerror(msg, data)
	e.Cause = err

	if le,ok := err.(Error); ok {
		e.Code = le.Code
	}
	
	return e
}

// NewError generates a new Error capturing the given human-oriented
// message and optionally structured data associated with this fault.
// The source code position to be reported by this Error is the point
// at which NewError was called.
func NewError(msg string, data ...interface{}) *Error {
	e := newerror(msg, data)
	e.Locate(1)
	return e
}

// WrapError generates a new Error capturing the given human-oriented
// message, a preceding error which caused this higher level fault,
// and optionally structured data associated with this fault.  The
// source code position to be reported by this Error is the point at
// which WrapError was called.
func WrapError(msg string, err error, data ...interface{}) *Error {
	e := wraperror(msg, err, data)
	e.Locate(1)
	return e
}

// Locate sets the source code position to be reported with this error
// as that point where the Locate call is made.  It should not
// generally be necessary to invoke this manually when using Logberry.
func (e *Error) Locate(skip int) {
	_, file, line, ok := runtime.Caller(skip + 1)
	if ok {
		e.Source.File = file
		e.Source.Line = line
	}
}

// SetCode associates the error with a particular error class string.
func (e *Error) SetCode(code string) {
	e.Code = code
}

// IsError checks if the given error is a Logberry Error tagged with
// any of the given codes, returning true if so and false otherwise.
func IsError(e error, code ...string) bool {

	err, ok := e.(*Error)
	if !ok {
		return false
	}

	for _, c := range code {
		if err.Code == c {
			return true
		}
	}

	return false

}

// Error returns a human-oriented serialization of the error.  It does
// not report the wrapped cause, if any.  That must be retrieved and
// reported manually.
func (e *Error) Error() string {

	var buffer = new(bytes.Buffer)

	buffer.WriteString(e.Message)

	if e.Code != "" {
		fmt.Fprintf(buffer, " <%v>", e.Code)
	}
	
	if e.Source.File != "" {
		fmt.Fprintf(buffer, " [%v:%v]", e.Source.File, e.Source.Line)
	}

	if len(e.Data) > 0 {
		fmt.Fprintf(buffer, " %v", e.Data.String())
	}

		if e.Cause != nil {
			fmt.Fprintf(buffer, "\n%v", e.Cause.Error())
		}

	return buffer.String()
}

// String returns a human-oriented serialization of the error.  It is
// the same as Error().
func (e *Error) String() string {
	return e.Error()
}
