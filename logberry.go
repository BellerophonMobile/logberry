package logberry

import (
	"os"
	"path"
	"log"
)


//----------------------------------------------------------------------
//----------------------------------------------------------------------
type ContextClass int
const (
	COMPONENT ContextClass = iota
	INSTANCE
	TASK
	contextclass_sentinel
)

var ContextClassText = [...]string {
	"component",
	"instance",
	"task",
}


type EventClass int
const (
	ERROR EventClass = iota
	FATAL
	WARNING
	INFO
	CONFIGURATION
	READY
	INSTANTIATE
	FINALIZE
	TASK_START
	TASK_FINISH
	RESOURCE
	SERVICE
	QUERY
	ASSERT
	CALCULATE
	READ
	WRITE
	CONNECT
	DISCONNECT
	eventclass_sentinel
)

var EventClassText = [...]string {
	"error",
	"fatal",
	"warning",
	"info",
	"configuration",
	"ready",
	"instantiate",
	"finalize",
	"task_start",
	"task_finish",
	"resource",
	"service",
	"query",
	"assert",
  "calculate",
  "read",
  "write",
  "connect",
  "disconnect",
}


//------------------------------------------------------
var Std *Root
var Main *Context


//------------------------------------------------------
func init() {

	//-- Check that labels are defined for all the enumerations
	if len(ContextClassText) != int(contextclass_sentinel) {
		log.Fatal("Fatal internal error: " +
			"len(ContextClassText) != |ContextClass|")
	}

	if len(EventClassText) != int(eventclass_sentinel) {
		log.Fatal("Fatal internal error: " +
			"len(EventClassText) != |EventClass|")
	}

	//-- Construct the standard default root
	Std = NewRoot(path.Base(os.Args[0]))
	Std.AddOutputDriver(NewStdOutput())

	//-- Construct the standard default context
	Main = Std.NewContext(COMPONENT, "main")

	// end init
}
