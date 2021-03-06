package logberry

import (
	"sync/atomic"
)

// Task represents a particular component, function, or activity.  In
// general a Task is meant to be used within a single thread of
// execution, and the calling code is responsible for managing any
// concurrent manipulation.
type Task struct {
	uid uint64

	root *Root

	parent *Task

	component string

	activity string

	data EventDataMap
}

var numtasks uint64

func newtaskuid() uint64 {
	// We have seen this atomic call cause problems on ARM...
	return atomic.AddUint64(&numtasks, 1) - 1
}

func newtask(parent *Task, component string, activity string, data []interface{}) *Task {

	t := &Task{
		uid:      newtaskuid(),
		parent:   parent,
		activity: activity,
		data:     EventDataMap{},
	}

	if parent != nil {
		t.root = parent.root
		t.component = parent.component
	} else {
		t.root = Std
	}

	if component != "" {
		t.component = component
	}

	t.root.event(t, BEGIN, t.activity+" begin", Aggregate(data))

	return t

}

// Task creates a new sub-task.  Parameter activity should be a short
// natural language description of the work that the Task represents,
// without any terminating punctuation.
func (x *Task) Task(activity string, data ...interface{}) *Task {
	return newtask(x, "", activity, data)
}

// Component creates a new Task object representing related long-lived
// functionality, rather than a directed, tightly scoped line of
// computation.  Parameter component should be a short lowercase
// string identifying the class, module, or other component that this
// Task represents.  The activity text of this Task is set to be
// "Component " + component.
func (x *Task) Component(component string, data ...interface{}) *Task {
	return newtask(x, component, "Component "+component, data)
}

// AddData incorporates the given data into that associated and
// reported with this Task.  This call does not generate a log event.
// The host Task is passed through as the return.  Among other things,
// this function is useful to silently accumulate data into the Task
// as it proceeds, to be reported when it concludes.
func (x *Task) AddData(data ...interface{}) *Task {
	for _, v := range data {
		x.data.Aggregate(v)
	}
	return x
}

// Event generates a user-specified log event.  Parameter event tags
// the class of the event, generally a short lowercase whitespace-free
// identifier.  A human-oriented text message is given as the msg
// parameter.  This should generally be static, short, use sentence
// capitalization but no terminating punctuation, and not itself
// include any data, which is better left to the structured data.  The
// variadic data parameter is aggregated as a D and reported with the
// event, as is the data permanently associated with the Task.  The
// given data is not associated to the Task permanently.
func (x *Task) Event(event string, msg string, data ...interface{}) {
	x.root.event(x, event, msg, Aggregate(data).Aggregate(x.data))
}

// Info generates an informational log event.  A human-oriented text
// message is given as the msg parameter.  This should generally be
// static, short, use sentence capitalization but no terminating
// punctuation, and not itself include any data, which is better left
// to the structured data.  The variadic data parameter is aggregated
// as a D and reported with the event, as is the data permanently
// associated with the Task.  The given data is not associated to the
// Task permanently.
func (x *Task) Info(msg string, data ...interface{}) {
	x.root.event(x, INFO, msg, Aggregate(data).Aggregate(x.data))
}

// Warning generates a warning log event indicating that a fault was
// encountered but the task is proceeding acceptably.  This should
// generally be static, short, use sentence capitalization but no
// terminating punctuation, and not itself include any data, which is
// better left to the structured data.  The variadic data parameter is
// aggregated as a D and reported with the event, as is the data
// permanently associated with the Task.  The given data is not
// associated to the Task permanently.
func (x *Task) Warning(msg string, data ...interface{}) {
	x.root.event(x, WARNING, msg, Aggregate(data).Aggregate(x.data))
}

// Ready generates a ready log event reporting that the activity or
// component the Task represents is initialized and prepared to begin.
// The variadic data parameter is aggregated as a D and reported with
// the event, as is the data permanently associated with the Task.
// The given data is not associated to the Task permanently.
func (x *Task) Ready(data ...interface{}) {
	x.root.event(x, READY, x.activity+" ready", Aggregate(data).Aggregate(x.data))
}

// Stopped generates a stopped log event reporting that the activity
// or component the Task represents has paused or shutdown.  The
// variadic data parameter is aggregated as a D and reported with the
// event, as is the data permanently associated with the Task.  The
// given data is not associated to the Task permanently.
func (x *Task) Stopped(data ...interface{}) {
	x.root.event(x, STOPPED, x.activity+" stopped", Aggregate(data).Aggregate(x.data))
}

// Finalized generates an end log event reporting that the component
// the Task represents has ceased.  It is generally intended to be
// used for components, while Success is used for discrete activities.
// Continuing to use the Task is discouraged.  The variadic data
// parameter is aggregated as a D and reported with the event, as is
// the data permanently associated with the Task.  The given data is
// not associated to the Task permanently.
func (x *Task) Finalized(data ...interface{}) {
	x.root.event(x, END, x.activity+" finalized", Aggregate(data).Aggregate(x.data))
}

// Success generates a success log event reporting that the activity
// the Task represents has concluded successfully.  It always returns
// nil.  Continuing to use the Task is discouraged.  The variadic data
// parameter is aggregated as a D and reported with the event, as is
// the data permanently associated with the Task.  The given data is
// not associated to the Task permanently.
func (x *Task) Success(data ...interface{}) error {
	x.root.event(x, SUCCESS, x.activity+" success", Aggregate(data).Aggregate(x.data))
	return nil
}

// Error generates an error log event reporting an unrecoverable fault
// in an activity or component.  If the given error is not a Logberry
// Error that has already been logged then it will be reported.  An
// error is returned wrapping the original error with a message
// reporting that the Task's activity has failed.  Continuing to use
// the Task is discouraged.  The variadic data parameter is aggregated
// as a D and embedded in the generated error.  It and the data
// permanently associated with the Task is reported with the event.
// The reported source code position of the generated task error is
// adjusted to be the event invocation.
func (x *Task) Error(cause error, data ...interface{}) *Error {

	m := x.activity + " failed"
	taskerr := wraperror(m, cause, data)
	taskerr.Locate(1) // Locate up the call stack

	taskerr.Reported = true

	d := Aggregate(data).Aggregate(D{"Source": taskerr.Source})

	dsub := d
	for cursor := cause; cursor != nil; {
		if ce, ok := cursor.(*Error); ok {
			if !ce.Reported {
				ce.Reported = true
				next := EventDataMap{}.Aggregate(ce)
				dsub["Cause"] = next
				dsub = next
				cursor = ce.Cause
			} else {
				break
			}
		} else {
			dsub["Cause"] = Copy(cursor)
			break
		}
	}

	d.Aggregate(x.data)

	x.root.event(x, ERROR, m, d)

	return taskerr

}

func (x *Task) WrapError(msg string, cause error, data ...interface{}) *Error {

	usererr := wraperror(msg, cause, nil)
	usererr.Reported = true
	usererr.Locate(1)

	m := x.activity + " failed"
	taskerr := wraperror(m, usererr, data)

	usererr.Reported = true
	taskerr.Reported = true

	d := Aggregate(data)

	dsub := Aggregate([]interface{}{usererr})
	d["Cause"] = dsub

	for cursor := cause; cursor != nil; {
		if ce, ok := cursor.(*Error); ok {
			if !ce.Reported {
				ce.Reported = true
				next := EventDataMap{}.Aggregate(ce)
				dsub["Cause"] = next
				dsub = next
				cursor = ce.Cause
			} else {
				break
			}
		} else {
			dsub["Cause"] = Copy(cursor)
			break
		}
	}

	d.Aggregate(x.data)

	x.root.event(x, ERROR, m, d)

	return taskerr

}

// Failure generates an error log event reporting an unrecoverable
// fault.  Failure and Error are essentially the same, the difference
// being that Failure is the first point of fault while Error takes an
// underlying error typically returned from another function or
// component.  An error is returned reporting that the activity or
// component represented by the Task has failed due to the underlying
// cause given in the message.  Continuing to use the Task is
// discouraged.  The variadic data parameter is aggregated as a D and
// embedded in the generated task error.  It and the data permanently
// associated with the Task is reported with the event.  The reported
// source code position of the generated task error is adjusted to be
// the event invocation.
func (x *Task) Failure(msg string, data ...interface{}) *Error {

	cause := newerror(msg, nil)
	cause.Locate(1)

	m := x.activity + " failed"
	taskerr := wraperror(m, cause, data)

	cause.Reported = true
	taskerr.Reported = true

	d := Aggregate(data)
	d["Cause"] = Copy(cause)
	d.Aggregate(x.data)

	x.root.event(x, ERROR, m, d)

	return taskerr

}
