package logberry

import (
	"fmt"
)


type LogError struct {
	Msg string
	Src []error
}


func NewError(data ...interface{}) *LogError {
	return &LogError{fmt.Sprint(data...), []error{}}
}

func WrapError(err error, data ...interface{}) *LogError {
	return &LogError{fmt.Sprint(data...), []error{err}}
}


func (e *LogError) AddError(err error) {
	e.Src = append(e.Src, err)
}


func (e *LogError) Error() string {

	if len(e.Src) == 1 {
		return e.Msg + "---" + e.Src[0].Error()
	}

	if len(e.Src) != 0 {
		s := e.Msg + ", multiple errors---" + e.Src[0].Error()
		for i := 1; i < len(e.Src); i++ {
			s += "; " + e.Src[i].Error()
		}
		return s;
	}

	return e.Msg;
}

func (e *LogError) String() string {
	return e.Error()
}
