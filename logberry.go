package logberry

import (
	"sync/atomic"

	"log"
	"os"
	"path"
)

//----------------------------------------------------------------------
//----------------------------------------------------------------------
type ComponentEventClass int

const (
	COMPONENT_START ComponentEventClass = iota
	COMPONENT_FINISH
	COMPONENT_CONFIGURATION
	COMPONENT_READY
	COMPONENT_INFO
	COMPONENT_WARNING
	COMPONENT_ERROR
	COMPONENT_FATAL
	componenteventclasssentinel
)

var ComponentEventClassText = [...]string{
	"start",
	"finish",
	"config",
	"ready",
	"info",
	"warning",
	"error",
	"fatal",
}

type ComponentClass int

const (
	COMPONENT ComponentClass = iota
	INSTANCE
	componentclasssentinel
)

var ComponentClassText = [...]string{
	"component",
	"instance",
}

type TaskEventClass int

const (
	TASK_BEGIN TaskEventClass = iota
	TASK_END
	TASK_INFO
	TASK_WARNING
	TASK_ERROR
	taskeventclasssentinel
)

var TaskEventClassText = [...]string{
	"begin",
	"end",
	"info",
	"warning",
	"error",
}

type TaskClass int

const (
	APPLICATION TaskClass = iota
	CALCULATION
	RESOURCE
	SERVICE
	CONNECT
	DISCONNECT
	/*
		SECURE
		UNSECURE
		ATTEST
		VERIFY
		RENDER
	*/
	taskclasssentinel
)

var TaskClassText = [...]string{
	"app",
	"calculation",
	"resource",
	"service",
	"connect",
	"disconnect",
}

type Context interface {
	GetUID() uint64
	GetLabel() string
	GetParent() Context
	GetRoot() *Root

	Component(label string, data ...interface{}) *Component

	Task(activity string, data ...interface{}) *Task
	LongTask(activity string, data ...interface{}) *Task

	IsHighlighted() bool

	CalculationTask(activity string, calculation interface{}, data ...interface{}) *Task
	LongCalculationTask(activity string, calculation interface{}, data ...interface{}) *Task

	ResourceTask(activity string, resource interface{}, data ...interface{}) *Task
	LongResourceTask(activity string, resource interface{}, data ...interface{}) *Task

	ServiceTask(activity string, service interface{}, query interface{}, data ...interface{}) *Task
	LongServiceTask(activity string, service interface{}, query interface{}, data ...interface{}) *Task
}

type highlightmarker int

var HIGHLIGHT highlightmarker = 0xDEADBEEF

//------------------------------------------------------
var Std *Root
var Main *Component

var numcontexts uint64

//------------------------------------------------------
func init() {

	//-- Check that labels are defined for the enumerations
	if len(ComponentEventClassText) != int(componenteventclasssentinel) {
		log.Fatal("Fatal internal error: " +
			"len(ComponentEventClassText) != |ComponentEventClass|")
	}

	if len(ComponentClassText) != int(componentclasssentinel) {
		log.Fatal("Fatal internal error: " +
			"len(ComponentClassText) != |ComponentClass|")
	}

	if len(TaskEventClassText) != int(taskeventclasssentinel) {
		log.Fatal("Fatal internal error: " +
			"len(TaskEventClassText) != |TaskEventClass|")
	}

	if len(TaskClassText) != int(taskclasssentinel) {
		log.Fatal("Fatal internal error: " +
			"len(TaskClassText) != |TaskClass|")
	}

	//-- Construct the standard default root
	Std = NewRoot(path.Base(os.Args[0]))
	Std.AddOutputDriver(NewStdOutput())

	//-- Construct the standard default context
	Main = Std.NewComponent("main")

	// end init
}

func newcontextuid() uint64 {
	return atomic.AddUint64(&numcontexts, 1) - 1
}

func InvalidComponentEventClass(event ComponentEventClass) bool {
	return (event < 0 || event >= componenteventclasssentinel)
}

func InvalidComponentClass(class ComponentClass) bool {
	return (class < 0 || class >= componentclasssentinel)
}

func InvalidTaskEventClass(event TaskEventClass) bool {
	return (event < 0 || event >= taskeventclasssentinel)
}

func InvalidTaskClass(class TaskClass) bool {
	return (class < 0 || class >= taskclasssentinel)
}
