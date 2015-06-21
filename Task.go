package logberry

import (
	"sync/atomic"
	"os"
	"os/user"
	"path/filepath"
	"path"
	"strings"
	"time"
)

// Task represents a particular component, function, or activity.  In
// general a Task is meant to be used within a single thread of
// execution, and the calling code is responsible for managing any
// concurrent manipulation.
type Task struct {
	uid    uint64

	root   Root

	parent *Task

	component  string

	activity string

	timed bool
	start time.Time

	data D

	mute bool
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
	return newtask(x, "Component " + component, data).SetComponent(component).Begin()
}

// GetUID returns the unique identifier for this Task.
func (x *Task) GetUID() uint64 {
	return x.uid
}

// GetRoot returns the Root for this Task.
func (x *Task) GetRoot() Root {
	return x.root
}

// GetParent returns the Context containing this Task.
func (x *Task) GetParent() *Task {
	return x.parent
}

// GetComponent returns the component label of this Task.
func (x *Task) GetComponent() string {
	return x.component
}

// SetComponent assigns the component label of this Task.
func (x *Task) SetComponent(c string) *Task{
	x.component = c
	return x
}

// GetActivity returns the activity text for this Task.
func (x *Task) GetActivity() string {
	return x.activity
}

// SetActivity sets the activity text for this Task.
func (x *Task) SetActivity(a string) *Task {
	x.activity = a
	return x
}

// GetTimed returns whether or not the Task is being timed.
func (x *Task) GetTimed() bool {
	return x.timed
}

// GetStart returns the timepoint at which this task began, indicated
// by calling Time.  Zero is returned if it is not being timed.
func (x *Task) GetStart() time.Time {
	return x.start
}

// Time indicates that this Task should be timed, starting now.  It
// does not generate a log event.
func (x *Task) Time() *Task {
	x.timed = true
	x.start = time.Now()
	return x
}

// Clock returns how long this Task has been running up to this point
// in time.  Zero is returned if the task in not being timed.  No log
// event is generated by this, but the duration is added to the Task's
// data.  Any previously set duration is overridden.
func (x *Task) Clock() time.Duration {

	if !x.timed {
		return 0
	}

	d := time.Now().Sub(x.start)
	x.data.Set("Duration", d)
	return d

}

// AddData incorporates the given key/value pair into the data
// associated and reported with this Task.  This call does not
// generate a log event.  The host Task is passed through as the
// return.  Among other things, this function is useful to silently
// accumulate data into the Task as it proceeds, to be reported when
// it concludes.
func (x *Task) AddData(k string, v interface{}) *Task {
	x.data.Set(k, v)
	return x
}

// AggregateData incorporates the given data into that associated and
// reported with this Task.  The rules for this construction are
// explained in AggregateFrom.  This call does not generate a log
// event.  The host Task is passed through as the return.  Among other
// things, this function is useful to silently accumulate data into
// the Task as it proceeds, to be reported when it concludes.
func (x *Task) AggregateData(data ...interface{}) *Task {
	x.data.CopyFrom(data)
	return x
}


// Mute indicates that this Task must not generate log events.  This
// is useful when using the Task merely to organize subtasks or
// generate informative error objects.
func (x *Task) Mute() *Task {
	x.mute = true
	return x
}

// Unmute indicates that this Task should no longer be muted.
func (x *Task) Unmute() *Task {
	x.mute = false
	return x
}

// IsMute returns true iff the Task is muted.
func (x *Task) IsMute() bool {
	return x.mute
}


// BuildMetadata generates a configuration log event reporting the
// build configuration, as captured by the passed object.  A utility
// script to generate such metadata automatically is in the util/
// folder of the Logberry repository.
func (x *Task) BuildMetadata(build *BuildMetadata) {
	x.root.event(x, CONFIGURATION, "Build metadata", DBuild(build))
}

// BuildSignature generates a configuration log event reporting build
// configuration, as captured by the given string.  A utility script
// to generate such metadata automatically is in the util/ folder of
// the Logberry repository.  It can be useful to use this string
// rather than a BuildMetadata object so that it can be passed in
// through the standard go tools command line, i.e., via linker flags.
func (x *Task) BuildSignature(build string) {
	x.root.event(x, CONFIGURATION, "Build signature", D{"Signature": build})
}

// Configuration generates a configuration log event reporting
// parameters or other initialization data.  The variadic data
// parameter is aggregated as a D and reporting with the event, as is
// the data permanently associated with the Task.  The given data is
// not associated to the Task permanently.
func (x *Task) Configuration(data ...interface{}) {
	d := DAggregate(append(data, x.data))
	x.root.event(x, CONFIGURATION, "Configuration", d)
}

// CommandLine generates a configuration log event reporting the
// command line used to execute the currently executing process.
func (x *Task) CommandLine() {

	hostname, err := os.Hostname()
	if err != nil {
		x.root.internalerror(WrapError("Could not retrieve hostname", err))
		return
	}

	u, err := user.Current()
	if err != nil {
		x.root.internalerror(WrapError("Could not retrieve user info", err))
		return
	}

	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		x.root.internalerror(WrapError("Could not retrieve program path", err))
		return
	}

	prog := path.Base(os.Args[0])

	d := D{
		"Host":    hostname,
		"User":    u.Username,
		"Path":    dir,
		"Program": prog,
		"Args":    os.Args[1:],
	}
	d.CopyFrom(x.data)
	
	x.root.event(x, CONFIGURATION, "Command line", d)

}

// Environment generates a configuration log event reporting the
// current operating system host environment variables of the
// currently executing process.
func (x *Task) Environment() {

	d := D{}
	for _, e := range os.Environ() {
		pair := strings.Split(e, "=")
		d.Set(pair[0], pair[1])
	}
	d.CopyFrom(x.data)

	x.root.event(x, CONFIGURATION, "Environment", d)

}

// Process generates a configuration log event reporting identifiers
// for the currently executing process.
func (x *Task) Process() {

	hostname, err := os.Hostname()
	if err != nil {
		x.root.internalerror(WrapError("Could not retrieve hostname", err))
		return
	}

	wd, err := os.Getwd()
	if err != nil {
		x.root.internalerror(WrapError("Could not retrieve working dir", err))
		return
	}

	u, err := user.Current()
	if err != nil {
		x.root.internalerror(WrapError("Could not retrieve user info", err))
		return
	}

	d := D{
		"Host": hostname,
		"WD":   wd,
		"UID":  u.Uid,
		"User": u.Username,
		"PID":  os.Getpid(),
	}
	d.CopyFrom(x.data)
	
	x.root.event(x, CONFIGURATION, "Process", d)

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

	if !x.mute {
		x.root.event(x, WARNING, msg, d)
	}
}


// Ready generates a ready log event reporting that the activity or
// component the Task represents is initialized and prepared to begin.
// The variadic data parameter is aggregated as a D and reporting with
// the event, as is the data permanently associated with the Task.
// The given data is not associated to the Task permanently.
func (x *Task) Ready(data ...interface{}) {
	x.root.event(x, READY, x.activity + " ready",
		DAggregate(data).CopyFrom(x.data))
}

// Stopped generates a stopped log event reporting that the activity
// or component the Task represents has paused or shutdown.  The
// variadic data parameter is aggregated as a D and reporting with the
// event, as is the data permanently associated with the Task.  The
// given data is not associated to the Task permanently.
func (x *Task) Stopped(data ...interface{}) {
	x.root.event(x, STOPPED, x.activity + " stopped",
		DAggregate(data).CopyFrom(x.data))
}


// Begin generates a begin log event reporting that the Task has been
// instantiated.  This is useful to report the start of a long-running
// activity, and is invoked when a component Task is created.  The
// host Task is passed through as the return.  The variadic data
// parameter is aggregated as a D and reporting with the event, as is
// the data permanently associated with the Task.  The given data is
// not associated to the Task permanently.
func (x *Task) Begin(data ...interface{}) *Task {
	x.Time()
	d := DAggregate(data)
	d.CopyFrom(x.data)

	if !x.mute {
		x.root.event(x, BEGIN, x.activity + " begin", d)
	}
	
	return x
}

// End generates an end log event reporting that the component the
// Task represents has been finalized.  If the Task is being timed it
// will be clocked and the duration reported.  Continuing to use the
// Task will not cause an error but is discouraged.  The variadic data
// parameter is aggregated as a D and reporting with the event, as is
// the data permanently associated with the Task.  The given data is
// not associated to the Task permanently.
func (x *Task) End(data ...interface{}) {

	x.Clock()
	d := DAggregate(data)
	d.CopyFrom(x.data)

	if !x.mute {
		x.root.event(x, END, x.activity + " end", d)
	}
	
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

	x.Clock()

	d := DAggregate(data)
	d.CopyFrom(x.data)

	if !x.mute {
		x.root.event(x, SUCCESS, x.activity + " success", d)
	}
	
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

	x.Clock()

	m := x.activity + " failed"

	e := wraperror(m, err, data)
	e.Locate(1)

	x.data.Set("Error", err)
	
	if !x.mute {
		x.root.event(x, ERROR, m, x.data)
	}
	
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

	x.Clock()

	m := x.activity + " failed"

	suberr := wraperror(msg, err, nil)
	
	e := wraperror(m, suberr, data)
	e.Locate(1)

	x.data.Set("Error", err)
	
	if !x.mute {
		x.root.event(x, ERROR, m, x.data)
	}
	
	return e

}

// Fatal is the same as Error except it terminates the program.  In
// general its use is discouraged outside of trivial programs.  For
// example, when using a BackgroundRoot the event may not be output
// because of the immediate cessation.  The variadic data parameter is
// aggregated as a D and reported as data embedded in the generated
// task error.  The data permanently associated with the Task is
// reported with the event.  Fatal does not respect Task muting.  The
// reported source code position of the generated task error is
// adjusted to be the event invocation.
func (x *Task) Fatal(err error, data ...interface{}) error {

	// This is all copied in so that the Locate call is correct...
	x.Clock()

	m := x.activity + " failed"

	e := wraperror(m, err, data)
	e.Locate(1)

	x.data.Set("Error", err)

	x.root.event(x, ERROR, m, x.data)

	os.Exit(-1)
	return e
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

	x.Clock()

	m := x.activity + " failed"

	e := wraperror(m, err, nil)
	e.Locate(1)

	x.data.Set("Error", err)
	
	if !x.mute {
		x.root.event(x, ERROR, m, x.data)
	}
	
	return e

}

// Die is the same as Failure except it terminates the program.  In
// general its use is discouraged outside of trivial programs.  For
// example, when using a BackgroundRoot the event may not be output
// because of the immediate cessation.  The variadic data parameter is
// aggregated as a D and reported as data embedded in the generated
// task error.  The data permanently associated with the Task is
// reported with the event.  Die does not respect Task muting.  The
// reported source code position of the generated task error is
// adjusted to be the event invocation.
func (x *Task) Die(msg string, data ...interface{}) error {

	// This is all copied in so that the Locate() call is correct...
	err := newerror(msg, data)
	err.Locate(1)

	x.Clock()

	m := x.activity + " failed"

	e := wraperror(m, err, nil)
	e.Locate(1)

	x.data.Set("Error", err)
	
	x.root.event(x, ERROR, m, x.data)

	os.Exit(-1)
	return nil

}
