package logberry

import (
	"os"
	"path"
)


var program string
var toplevel string = "main"

func init() {
	program = path.Base(os.Args[0])
}


//----------------------------------------------------------------------
//----------------------------------------------------------------------
func SetProgram(label string) {
	program = label
}

func SetTopLevel(label string) {
	toplevel = label
}


//----------------------------------------------------------------------
//----------------------------------------------------------------------
func Build(build BuildMetadata) {
	LogPrimitive(toplevel, CONFIGURATION, "Build", DBuild(build))
}

func CommandLine() {
	LogPrimitive(toplevel, CONFIGURATION, "Command line",
		&D{"Args": os.Args})
}

func Configuration(data ...interface{}) {
	LogPrimitive(toplevel, CONFIGURATION, "Configuration", DAggregate(data))
}

func Ready(msg string, data ...interface{}) {
	LogPrimitive(toplevel, READY, msg, DAggregate(data))
}

//------------------------------------------------------
func Info(msg string, data ...interface{}) {
	LogPrimitive(toplevel, INFO, msg, DAggregate(data))
}

func Warning(msg string, data ...interface{}) {
	LogPrimitive(toplevel, WARNING, msg, DAggregate(data))
}

func Error(msg string, err error, data ...interface{}) error {
	LogPrimitive(toplevel, ERROR, msg,
		DAggregate(data).Set("Error", err.Error()))
	return WrapError(err, msg)
}

// Failure is the same as Error but doesn't take an error object.
func Failure(msg string, data ...interface{}) error {
	LogPrimitive(toplevel, ERROR, msg, DAggregate(data))
	return NewError(msg)
}

// Only the top level should invoke fatal, not components.
func Fatal(msg string, err error, data ...interface{}) {
	LogPrimitive(toplevel, ERROR, msg,
		DAggregate(data).Set("Error", err.Error()))
	os.Exit(1)
}
