package logberry

import (
	"os"
	"os/user"
	"path"
	"path/filepath"
	"strings"
	"sync/atomic"
)

// Task represents a particular component, function, or activity.  In
// general a Task is meant to be used within a single thread of
// execution, and the calling code is responsible for managing any
// concurrent manipulation.
type Task struct {
	uid uint64

	root Root

	parent *Task

	component string

	activity string

	data D
}

var numtasks uint64

func newtaskuid() uint64 {
	// We have seen this atomic call cause problems on ARM...
	return atomic.AddUint64(&numtasks, 1) - 1
}

func newtask(parent *Task, activity string, data []interface{}) *Task {

	t := &Task{
		uid:    newtaskuid(),
		parent: parent,

		activity: activity,

		data: DAggregate(data),
	}

	if parent != nil {
		t.root = parent.root
		t.component = parent.component
	} else {
		t.root = Std
	}

	x.root.event(x, BEGIN, x.activity+" begin", d)

	return t

}

// Task creates a new sub-task of the host task.  Parameter activity
// should be a short natural language description of the work that the
// Task represents, without any terminating punctuation.  Any data
// given will be associated with the Task and reported with its
// events.  This call does not produce a log event.  Use Begin to
// indicate the start of a long running task.
func (x *Task) Task(activity string, data ...interface{}) *Task {
	return newtask(x, activity, data)
}

// Component creates a new Task object representing a set of related
// sub-functionality under the host Task, rather than a directed,
// tightly scoped line of computation.  Parameter component should be
// a short lowercase string identifying the class, module, or other
// component that this Task represents.  The activity text of this
// Task is set to be "Component " + component.  Any data given will be
// associated with the Task and reported with its events.  This call
// does produce a log event marking the instantiation.
func (x *Task) Component(component string, data ...interface{}) *Task {
	return newtask(x, component, data).SetComponent(component).Begin()
}

// AddData incorporates the given data into that associated and
// reported with this Task.  The rules for this construction are
// explained in AggregateFrom.  This call does not generate a log
// event.  The host Task is passed through as the return.  Among other
// things, this function is useful to silently accumulate data into
// the Task as it proceeds, to be reported when it concludes.
func (x *Task) AddData(data ...interface{}) *Task {
	x.data.CopyFrom(data)
	return x
}

// Event generates a user-specified log event.  Parameter event tags
// the class of the event, generally a short lowercase whitespace-free
// identifier.  A human-oriented text message is given as the msg
// parameter.  This should generally be static, short, use sentence
// capitalization but no terminating punctuation, and not itself
// include any data, which is better left to the structured data.  The
// variadic data parameter is aggregated as a D and reporting with the
// event, as is the data permanently associated with the Task.  The
// given data is not associated to the Task permanently.
func (x *Task) Event(event string, msg string, data ...interface{}) {
	x.root.event(x, event, msg, DAggregate(data).CopyFrom(x.data))
}

// Info generates an informational log event.  A human-oriented text
// message is given as the msg parameter.  This should generally be
// static, short, use sentence capitalization but no terminating
// punctuation, and not itself include any data, which is better left
// to the structured data.  The variadic data parameter is aggregated
// as a D and reporting with the event, as is the data permanently
// associated with the Task.  The given data is not associated to the
// Task permanently.
func (x *Task) Info(msg string, data ...interface{}) {
	x.root.event(x, INFO, msg, DAggregate(data).CopyFrom(x.data))
}

// Warning generates a warning log event reporting that a fault was
// encountered but the task is proceeding acceptably.  This should
// generally be static, short, use sentence capitalization but no
// terminating punctuation, and not itself include any data, which is
// better left to the structured data.  The variadic data parameter is
// aggregated as a D and reporting with the event, as is the data
// permanently associated with the Task.  The given data is not
// associated to the Task permanently.
func (x *Task) Warning(msg string, data ...interface{}) {
	d := DAggregate(data)
	d.CopyFrom(x.data)

	x.root.event(x, WARNING, msg, d)
}

// Ready generates a ready log event reporting that the activity or
// component the Task represents is initialized and prepared to begin.
// The variadic data parameter is aggregated as a D and reporting with
// the event, as is the data permanently associated with the Task.
// The given data is not associated to the Task permanently.
func (x *Task) Ready(data ...interface{}) {
	x.root.event(x, READY, x.activity+" ready",
		DAggregate(data).CopyFrom(x.data))
}

// Stopped generates a stopped log event reporting that the activity
// or component the Task represents has paused or shutdown.  The
// variadic data parameter is aggregated as a D and reporting with the
// event, as is the data permanently associated with the Task.  The
// given data is not associated to the Task permanently.
func (x *Task) Stopped(data ...interface{}) {
	x.root.event(x, STOPPED, x.activity+" stopped",
		DAggregate(data).CopyFrom(x.data))
}

// End generates an end log event reporting that the component the
// Task represents has been finalized.  If the Task is being timed it
// will be clocked and the duration reported.  Continuing to use the
// Task will not cause an error but is discouraged.  The variadic data
// parameter is aggregated as a D and reporting with the event, as is
// the data permanently associated with the Task.  The given data is
// not associated to the Task permanently.
func (x *Task) End(data ...interface{}) {

	d := DAggregate(data)
	d.CopyFrom(x.data)

	x.root.event(x, END, x.activity+" end", d)

}

// Success generates a success log event reporting that the activity
// the Task represents has concluded successfully.  If the Task is
// being timed it will be clocked and the duration reported.  It
// always returns nil.  Continuing to use the Task will not cause an
// error but is discouraged.  The variadic data parameter is
// aggregated as a D and reporting with the event, as is the data
// permanently associated with the Task.  The given data is not
// associated to the Task permanently.
func (x *Task) Success(data ...interface{}) error {

	d := DAggregate(data)
	d.CopyFrom(x.data)

	x.root.event(x, SUCCESS, x.activity+" success", d)

	return nil

}

// Error generates an error log event reporting an unrecoverable fault
// in an activity or component.  If the Task is being timed it will be
// clocked and the duration reported.  An error is returned wrapping
// the original error with a message reporting that the Task's
// activity has failed.  Continuing to use the Task will not cause an
// error but is discouraged.  The variadic data parameter is
// aggregated as a D and reported as data embedded in the generated
// task error.  The data permanently associated with the Task is
// reported with the event.  The reported source code position of the
// generated task error is adjusted to be the event invocation.
func (x *Task) Error(err error, data ...interface{}) error {

	m := x.activity + " failed"

	x.data.Set("Error", err)

	e := wraperror(m, err, data)
	e.Locate(1)
	e.Data.CopyFrom(x.data)

	x.root.event(x, ERROR, m, x.data)

	return e

}

// WrapError generates an error log event reporting an unrecoverable
// fault in an activity or component.  It is similar to Error but
// useful for taking a causal error, wrapping it in an additional
// message, and then throwing a task error.  If the Task is being
// timed it will be clocked and the duration reported.  An error is
// returned reporting that the activity or component represented by
// the Task has failed, wrapping a causal error with the given
// message, which in turn wraps the given root error.  Continuing to
// use the Task will not cause an error but is discouraged.  The
// variadic data parameter is aggregated as a D and reported as data
// embedded in the generated task error.  The data permanently
// associated with the Task is reported with the event.  The reported
// source code position of the generated task error is adjusted to be
// the event invocation.
func (x *Task) WrapError(msg string, err error, data ...interface{}) error {

	m := x.activity + " failed"

	suberr := wraperror(msg, err, nil)
	suberr.Locate(1)
	x.data.Set("Error", suberr)

	x.root.event(x, ERROR, m, x.data)

	return suberr

}

// Failure generates an error log event reporting an unrecoverable
// fault.  Failure and Error are essentially the same, the difference
// being that Failure is useful to both report and generate a fault
// detected directly by the calling code.  Error in contrast takes an
// underlying error, typically as returned from another function or
// component.  If the Task is being timed it will be clocked and the
// duration reported.  An error is returned reporting that the
// activity or component represented by the Task has failed due to the
// underlying cause given in the message.  Continuing to use the Task
// will not cause an error but is discouraged.  The variadic data
// parameter is aggregated as a D and reported as data embedded in the
// generated task error.  The data permanently associated with the
// Task is reported with the event.  The reported source code position
// of the generated task error is adjusted to be the event invocation.
func (x *Task) Failure(msg string, data ...interface{}) error {

	err := newerror(msg, data)
	err.Locate(1)

	m := x.activity + " failed"

	x.data.Set("Error", err)

	e := wraperror(m, err, nil)
	e.Locate(1)
	e.Data.CopyFrom(x.data)

	x.root.event(x, ERROR, m, x.data)

	return e

}
