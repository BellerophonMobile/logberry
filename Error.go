package logberry

import (
	"runtime"
	"bytes"
	"fmt"
)

// Error is used to report a fault, capturing a human-oriented message
// describing the problem, structured data providing identifying
// details, the source code file name and line number location at
// which this Error was generated, and if appropriate a preceding
// error that caused this higher level fault.
type Error struct {
	Message string
	Data D

	File string
	Line int

	Cause error
}


func newerror(msg string, data []interface{}) *Error {	
	e := &Error{
		Message:msg,
		Data: DAggregate(data),
	}
	return e
}

func wraperror(msg string, err error, data []interface{}) *Error {
	e := newerror(msg, data)

	if _,ok := err.(*Error); ok {
		e.Cause = err
	} else {
		e.Cause = newerror(err.Error(), nil)
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
// as that point where the Locate call is made.  In general it should
// not be necessary to invoke this manually.
func (e *Error) Locate(skip int) {
	_,file,line,ok := runtime.Caller(skip+1)
	if ok {
		e.File = file
		e.Line = line
	}
}

// Error implements the standard Go error interface, returning a
// human-oriented text string serialization of the Error.
func (e *Error) Error() string {

	var buffer = new(bytes.Buffer)

	buffer.WriteString(e.Message)

	if e.File != "" {
		fmt.Fprintf(buffer, " [%v:%v]", e.File, e.Line)
	}
	
	if len(e.Data) > 0 {
		fmt.Fprintf(buffer, " %v", e.Data.String())
	}

	if e.Cause != nil {
		fmt.Fprintf(buffer, ":%v", e.Error())
	}

	return buffer.String()
}

// String returns a human-oriented text string serialization of the
// Error.
func (e *Error) String() string {
	return e.Error()
}
