package logberry

import (
	"os"
	"path"
	"log"
	"time"
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
	INSTANTIATE
	FINALIZE
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
	"instantiate",
	"finalize",
	"task_start",
	"task_finish",
	"resource",
	"service",
	"unknown",
};

//------------------------------------------------------
type Data map[string]interface{}

type ErrorListener interface {
	LoggingError(err error)
}


//------------------------------------------------------
var outputdrivers = []OutputDriver{}
var errorlisteners = []ErrorListener{}

var program string
var toplevel string = "main"

func init() {

	if len(STATEMENT_CLASS_TEXT) != int(UNKNOWN) + 1 {
		log.Fatal("Fatal internal error: len(STATEMENT_CLASS_TEXT) != |StatementClass|")
	}

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

func AddOutput(driver OutputDriver) {
	outputdrivers = append(outputdrivers, driver)
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
func logprimitive(component string,
	                class StatementClass,
	                msg string,
                  data interface{}) {

	for _,driver := range(outputdrivers) {
		driver.Log(component, class, msg, data)
	}

	// end logprimitive
}


//----------------------------------------------------------------------
//----------------------------------------------------------------------
func Build(build BuildMetadata) {
	logprimitive(toplevel, CONFIGURATION, "Build", build)
}

func CommandLine() {
	logprimitive(toplevel, CONFIGURATION, "Command line", &Data{"Args": os.Args})
}

func Configuration(data interface{}) {
	logprimitive(toplevel, CONFIGURATION, "Configuration", data)
}

/*
 * Only the top level should invoke fatal, not components
 */
func Fatal(msg string, err error) {
	logprimitive(toplevel, FATAL, msg,
		&Data{ "Error": err.Error() })
	os.Exit(1)
}


//----------------------------------------------------------------------
//----------------------------------------------------------------------
type ComponentLog struct {
	Component string
}

func NewComponent(component string, data interface{}) *ComponentLog {
	logprimitive(component, INSTANTIATE, "Instantiate", data)
	return &ComponentLog{
		Component: component,
	}
}

func (log *ComponentLog) Finalize() {
	logprimitive(log.Component, FINALIZE, "Finalize", nil)
}

func (log *ComponentLog) Build(build BuildMetadata) {
	logprimitive(log.Component, CONFIGURATION, "Build", build)
}

func (log *ComponentLog) Info(msg string, data interface{}) {
	logprimitive(log.Component, INFO, msg, data)
}

func (log *ComponentLog) Error(msg string, err error) error {
	// Note that this can't/shouldn't just throw err into the data blob
	// because the standard errors package error doesn't expose
	// anything, even the message.  So you basically have to reduce to a
	// string via Error().
	e := WrapError(err, msg)
	logprimitive(log.Component, ERROR, msg, &Data{ "Error": err.Error() })
	return e
}

func (log *ComponentLog) Failure(msg string) error {
	e := NewError(msg)
	logprimitive(log.Component, ERROR, msg, nil)
	return e
}

func (log *ComponentLog) Warning(msg string) {
	logprimitive(log.Component, WARNING, msg, nil)
}

func (log *ComponentLog) Resource(msg string, resource interface{}) {
	logprimitive(log.Component, RESOURCE, msg, &Data{"Resource": resource})
}


//----------------------------------------------------------------------
type Task struct {
	Component *ComponentLog
	Class StatementClass
	Msg string
	Data interface{}

	Long bool
	Start time.Time

}


func (task *Task) Error(err error) error {
	if task.Long {
		task.Msg = "@Error " + task.Msg
	}

	task.Msg += " failed"
	e := WrapError(err, task.Msg)
	logprimitive(task.Component.Component, ERROR, task.Msg,
		&Data{ "Error": err.Error() })
	return e
}

func (task *Task) Failure(msg string) error {
	if task.Long {
		task.Msg = "@Error " + task.Msg
	}

	task.Msg += " failed"
	e := WrapError(NewError(msg), task.Msg)
	logprimitive(task.Component.Component, ERROR, task.Msg,
		&Data{ "Error": msg })
	return e
}

func (task *Task) Success() {

	// If this was a longrunning task, compute duration and alter fields
	if task.Long {
		duration := time.Now().Sub(task.Start)

		task.Class = TASK_FINISH
		task.Msg = "@Finish " + task.Msg
		task.Data = &Data{"Data": task.Data, "Duration": duration}
	}

	logprimitive(task.Component.Component, task.Class, task.Msg, task.Data)
}


//------------------------------------------------------
func (log *ComponentLog) Task(msg string, resource interface{}) *Task {
	return &Task{
		Component: log,
		Class: RESOURCE,
		Msg: msg,
		Data: &Data{"Resource": resource},
	}
}

func (log *ComponentLog) LongTask(msg string, data interface{}) *Task {
	logprimitive(log.Component, TASK_START, "@Start " + msg, data)

	return &Task{
		Component: log,
		Class: TASK_START,
		Msg: msg,
		Data: data,

		Long: true,
		Start: time.Now(),
	}
}

func (log *ComponentLog) ResourceTask(msg string, resource interface{}) *Task {
	return &Task{
		Component: log,
		Class: RESOURCE,
		Msg: msg,
		Data: &Data{"Resource": resource},
	}
}

func (log *ComponentLog) ServiceTask(msg string, service interface{}) *Task {
	return &Task{
		Component: log,
		Class: SERVICE,
		Msg: msg,
		Data: &Data{"Service": service},
	}
}
