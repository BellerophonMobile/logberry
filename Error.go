package logberry

import (
	"runtime"
	"bytes"
	"fmt"
)


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

func (e *Error) Locate(skip int) {
	_,file,line,ok := runtime.Caller(skip+1)
	if ok {
		e.File = file
		e.Line = line
	}
}

func NewError(msg string, data ...interface{}) *Error {
	e := newerror(msg, data)
	e.Locate(1)
	return e
}

func WrapError(msg string, err error, data ...interface{}) *Error {
	e := wraperror(msg, err, data)
	e.Locate(1)
	return e
}


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

func (e *Error) String() string {
	return e.Error()
}
