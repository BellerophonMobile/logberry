package logberry

import (
	"log"
//	"time"
)


//----------------------------------------------------------------------
//----------------------------------------------------------------------
type StatementClass int

const (
	ERROR StatementClass = iota
	FATAL
	WARNING
	INFO
	CONFIGURATION
	READY
	INSTANTIATE
	FINALIZE
	TASK
	TASK_START
	TASK_FINISH
	RESOURCE
	SERVICE
	UNKNOWN // This is the sentinel, must be last!
)

var STATEMENT_CLASS_TEXT = [...]string {
	"error",
	"fatal",
	"warning",
	"info",
	"configuration",
	"ready",
	"instantiate",
	"finalize",
	"task",
	"task_start",
	"task_finish",
	"resource",
	"service",
	"unknown",
};


//------------------------------------------------------
type ErrorListener interface {
	LoggingError(err error)
}


//------------------------------------------------------
var outputdrivers = []OutputDriver{}
var errorlisteners = []ErrorListener{}

func init() {
	if len(STATEMENT_CLASS_TEXT) != int(UNKNOWN) + 1 {
		log.Fatal("Fatal internal error: " +
			"len(STATEMENT_CLASS_TEXT) != |StatementClass|")
	}

	AddOutputDriver(NewStdOutput())
}


//----------------------------------------------------------------------
//----------------------------------------------------------------------
func AddOutputDriver(driver OutputDriver) OutputDriver {
	outputdrivers = append(outputdrivers, driver)
	return driver
}

func SetOutputDriver(driver OutputDriver) OutputDriver {
	ClearOutputDrivers()
	return AddOutputDriver(driver)
}

func ClearOutputDrivers() {
	outputdrivers = []OutputDriver{}
}

func AddErrorListener(listener ErrorListener) {
	errorlisteners = append(errorlisteners, listener)
}


//----------------------------------------------------------------------
//----------------------------------------------------------------------
/*
 * Report an error that occurred in logging itself.
 * Public so OutputDrivers in other packages can utilize.
 */
func LoggingError(err error) {
	for _,listener := range(errorlisteners) {
		listener.LoggingError(err)
	}
	// end logerror
}

/*
 * Internal multiplexer out to all active OutputDrivers.
 */
func LogPrimitive(component string,
	                class StatementClass,
	                msg string,
                  data *D) {

	for _,driver := range(outputdrivers) {
		driver.Log(component, class, msg, data)
	}

	// end LogPrimitive
}
