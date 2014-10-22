package logberry

import (
	"sync/atomic"

	"log"
	"os"
	"path"
)

// ComponentEventClass captures a simple enumeration of types of
// events that may happen in a component's lifecycle.
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
	componenteventclass_sentinel
)

// ComponentEventClassText is an array mapping the ComponentEventClass
// enumeration into text labels.
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

// ComponentClass captures a simple enumeration of types of
// Components.  These are either components, to be taken as either
// classes or just clusters of related functionality, or instances,
// specific instantiations of objects.
type ComponentClass int

const (
	COMPONENT ComponentClass = iota
	INSTANCE
	componentclass_sentinel
)

// ComponentClassText is an array mapping the ComponentClass
// enumeration into text labels.
var ComponentClassText = [...]string{
	"component",
	"instance",
}

// TaskEventClass captures a simple enumeration of types of major
// events that may happen in a task's lifecycle.
type TaskEventClass int

const (
	TASK_BEGIN TaskEventClass = iota
	TASK_END
	TASK_INFO
	TASK_WARNING
	TASK_ERROR
	taskeventclass_sentinel
)

// TaskEventClassText is an array mapping the TaskEventClass
// enumeration into text labels.
var TaskEventClassText = [...]string{
	"begin",
	"end",
	"info",
	"warning",
	"error",
}

// Context is an interface for objects representing entities that
// generate logging events.  It encompasses aspects common to both
// Components and Tassks.
type Context interface {
	GetUID() uint64
	GetLabel() string
	GetParent() Context
	GetRoot() *Root

	IsMute() bool
	IsHighlighted() bool

	Component(label string, data ...interface{}) *Component
	Task(activity string, data ...interface{}) *Task
}

// Std is the default Root created at startup.
var Std *Root

// Main is the default Component created at startup, roughly intended
// to represent main program execution.
var Main *Component

var numcontexts uint64

func init() {

	//-- Check that labels are defined for the enumerations
	if len(ComponentEventClassText) != int(componenteventclass_sentinel) {
		log.Fatal("Fatal internal error: " +
			"len(ComponentEventClassText) != |ComponentEventClass|")
	}

	if len(ComponentClassText) != int(componentclass_sentinel) {
		log.Fatal("Fatal internal error: " +
			"len(ComponentClassText) != |ComponentClass|")
	}

	if len(TaskEventClassText) != int(taskeventclass_sentinel) {
		log.Fatal("Fatal internal error: " +
			"len(TaskEventClassText) != |TaskEventClass|")
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

// InvalidComponentEventClass returns true if event is within the
// known enumeration of component events.
func InvalidComponentEventClass(event ComponentEventClass) bool {
	return (event < 0 || event >= componenteventclass_sentinel)
}

// InvalidComponentClass returns true if class is within the known
// enumeration of component classes.
func InvalidComponentClass(class ComponentClass) bool {
	return (class < 0 || class >= componentclass_sentinel)
}

// InvalidTaskEventClass returns true if event is within the known
// enumeration of task events.
func InvalidTaskEventClass(event TaskEventClass) bool {
	return (event < 0 || event >= taskeventclass_sentinel)
}
