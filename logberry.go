package logberry

import (
	"sync/atomic"

	"os"
	"path"
	"log"
)


//----------------------------------------------------------------------
//----------------------------------------------------------------------
type ContextEventClass int
const (
	ERROR ContextEventClass = iota
	FATAL
	WARNING
	INFO
	CONFIGURATION
	START
	FINISH
	SUCCESS
	contexteventclasssentinel
)

var ContextEventClassText = [...]string {
	"error",
	"fatal",
	"warning",
	"info",
	"configuration",
	"start",
	"finish",
	"success",
}


type ComponentClass int
const (
	COMPONENT ComponentClass = iota
	INSTANCE
	componentclasssentinel
)

var ComponentClassText = [...]string {
	"component",
	"instance",
}


type ActivityClass int
const (
	APPLICATION ActivityClass = iota
	CALCULATION
	RESOURCE
	SERVICE
/*
	CONNECT
	DISCONNECT
	SECURE
	UNSECURE
	ATTEST
	VERIFY
	RENDER
*/
	activityclasssentinel
)

var ActivityClassText = [...]string {
	"app",
	"calculation",
	"resource",
	"service",
}


type Context interface {
	GetUID() uint64
	GetLabel() string
	GetParent() Context
	GetRoot() *Root

	Component(label string, data ...interface{}) *Component

	Task(activity string, data ...interface{}) *Task
	LongTask(activity string, data ...interface{}) *Task
}


//------------------------------------------------------
var Std *Root
var Main *Component

var numcontexts uint64


//------------------------------------------------------
func init() {

	//-- Check that labels are defined for the enumerations
	if len(ContextEventClassText) != int(contexteventclasssentinel) {
		log.Fatal("Fatal internal error: " +
			"len(ContextEventClassText) != |ContextEventClass|")
	}

	if len(ActivityClassText) != int(activityclasssentinel) {
		log.Fatal("Fatal internal error: " +
			"len(ActivityClassText) != |ActivityClass|")
	}

	if len(ComponentClassText) != int(componentclasssentinel) {
		log.Fatal("Fatal internal error: " +
			"len(ComponentClassText) != |ComponentClass|")
	}


	//-- Construct the standard default root
	Std = NewRoot(path.Base(os.Args[0]))
	Std.AddOutputDriver(NewStdOutput())

	//-- Construct the standard default context
	Main = Std.NewComponent("main")

	// end init
}


func newcontextuid() uint64 {
	return atomic.AddUint64(&numcontexts, 1)-1
}
