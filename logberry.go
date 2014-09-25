package logberry

import (
	"os"
	"path"
)


//----------------------------------------------------------------------
//----------------------------------------------------------------------
var program string

func init() {
	program = path.Base(os.Args[0])
}


//----------------------------------------------------------------------
//----------------------------------------------------------------------

type StatementClass int

const (
	ERROR StatementClass = iota
	FATAL
	METADATA
	INFO
)

var CLASSTEXT = [...]string {
	"error",
	"fatal",
	"metadata",
	"info",
};


type Data map[string]interface{}


//------------------------------------------------------
var outputdrivers = []OutputDriver{}


//----------------------------------------------------------------------
//----------------------------------------------------------------------
func AddOutput(driver OutputDriver) error {

	outputdrivers = append(outputdrivers, driver)
	return nil

}


//----------------------------------------------------------------------
//----------------------------------------------------------------------
func logprimitive(component string,
	                class StatementClass,
	                msg string,
                  data interface{}) error {

	var accumerror *logerror = nil

	for _,driver := range(outputdrivers) {
		err := driver.log(component, class, msg, data)
		if err != nil {
			if accumerror == nil {
				accumerror = NewError("Error outputting to log(s)")
			}
			accumerror.AddError(err)
		}
	}

	return accumerror
}


//----------------------------------------------------------------------
//----------------------------------------------------------------------
func SetProgram(label string) {
	program = label
}

func Build(build BuildMetadata) error {
	return logprimitive(program, METADATA, "Build", build)
}

func CommandLine() error {
	return logprimitive(program, METADATA, "Command line",
		&Data{"Args": os.Args})
}

func Configuration(data interface{}) error {
	return logprimitive(program, METADATA, "Configuration", data)
}

func FatalError(msg string, err error) error {
	logprimitive(program, FATAL, msg,
		&Data{ "Error": err.Error() })
	os.Exit(1)
	return nil
}

//----------------------------------------------------------------------
//----------------------------------------------------------------------
type ComponentLog struct {
	Component string
}

func NewComponentLog(component string) *ComponentLog {
	return &ComponentLog{
		Component: component,
	}
}

func (log *ComponentLog) Build(build BuildMetadata) error {
	return logprimitive(log.Component, METADATA, "Build", build)
}

func (log *ComponentLog) Info(msg string, data interface{}) error {
	return logprimitive(log.Component, INFO, msg, data)
}

func (log *ComponentLog) Error(msg string, err error) error {
	return logprimitive(log.Component, ERROR, msg,
		&Data{ "Error": err.Error() })
}
