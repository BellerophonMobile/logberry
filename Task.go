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

	mute      bool
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
		t.data.Set("Parent", t.parent.uid)
	} else {
		t.root = Std
	}

	return t

}


// Task creates a new task context as a child of this task.  Parameter
// activity should be a short natural language description of the work
// carried out by this activity, without any terminating punctuation.
// Any data given will be associated with the task and reported with
// its events.  This call does not produce a log event.  Use Begin()
// to indicate the start of a long running task.
//
// This is safe to call concurrently.
func (x *Task) SubTask(activity string, data ...interface{}) *Task {
	return newtask(x, activity, data)
}

func (x *Task) SubComponent(component string, data ...interface{}) *Task {
	return newtask(x, component, data).SetComponent(component)
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

func (x *Task) GetComponent() string {
	return x.component
}

func (x *Task) SetComponent(c string) *Task{
	x.component = c
	return x
}

// GetActivity returns the work description for this Task.
func (x *Task) GetActivity() string {
	return x.activity
}

func (x *Task) SetActivity(a string) *Task {
	x.activity = a
	return x
}

// GetTimed returns whether or not the Task is being timed.
func (x *Task) GetTimed() bool {
	return x.timed
}

// GetStart returns the timepoint at which this task began, indicated
// by calling Time().  Zero is returned if it is not being timed.
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
// in time.  Zero is returned if the task in not being timed (Time()
// was not called).  No log event is generated by this, but the
// duration is added to the Task's data.  Any previously set duration
// is overridden.
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
// explained in D.AggregateFrom().  This call does not generate a log
// event.  The host Task is passed through as the return.  Among other
// things, this function is useful to silently accumulate data into
// the Task as it proceeds, to be reported when it concludes.
func (x *Task) AggregateData(data ...interface{}) *Task {
	x.data.CopyFrom(data)
	return x
}

// Mute indicates that this Task should not generate log events except
// errors.  This is useful when using the Task merely to organize subtasks
// or purely to report failures.
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




// Calculation is a data modifier that adds the given expression or
// result to the Task under the field "Calculation".
//
// This call does not generate a log event.  The host Task is passed
// through as the return.
func (x *Task) Calculation(calculation interface{}) *Task {
	x.AddData("Calculation", calculation)
	return x
}

// File is a data modifier that adds the given filename or descriptor
// to the Task under the field "File".
//
// This call does not generate a log event.  The host Task is passed
// through as the return.
func (x *Task) File(file interface{}) *Task {
	x.AddData("File", file)
	return x
}

// Resource is a data modifier that adds the given resource identifier
// or description to the Task under the field "Resource".
//
// This call does not generate a log event.  The host Task is passed
// through as the return.
func (x *Task) Resource(resource interface{}) *Task {
	x.AddData("Resource", resource)
	return x
}

// Service is a data modifier that adds the given service identifier
// or description to the Task under the field "Service".
//
// This call does not generate a log event.  The host Task is passed
// through as the return.
func (x *Task) Service(service interface{}) *Task {
	x.AddData("Service", service)
	return x
}

// User is a data modifier that adds the given name or description to
// the Task under the field "User".
//
// This call does not generate a log event.  The host Task is passed
// through as the return.
func (x *Task) User(user interface{}) *Task {
	x.AddData("User", user)
	return x
}

// Endpoint is a data modifier that adds the given target identifier
// or description to the Task under the field "Endpoint".
//
// This call does not generate a log event.  The host Task is passed
// through as the return.
func (x *Task) Endpoint(endpoint interface{}) *Task {
	x.AddData("Endpoint", endpoint)
	return x
}

// BuildMetadata reports on the build configuration, as captured by
// the passed object.  A utility script to generate such metadata
// automatically from a git repository is in
// util/build-metadata-go.sh.
func (x *Task) BuildMetadata(build *BuildMetadata) {
	x.root.event(x, CONFIGURATION, "Build metadata", DBuild(build))
}

// BuildSignature reports on the build configuration, as captured by
// the passed string.  A utility script to generate such metadata
// automatically from a git repository is in
// util/build-signature-go.sh.  It can be useful to use this string
// rather than a BuildMetadata object so that it can be passed in
// through the standard go tools, i.e., via linker flags.
func (x *Task) BuildSignature(build string) {
	x.root.event(x, CONFIGURATION, "Build signature", D{"Signature": build})
}

// Configuration generates a configuration log event reporting
// parameters or other initialization data.
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


func (x *Task) End(data ...interface{}) error {

	x.Clock()

	d := DAggregate(data)
	d.CopyFrom(x.data)

	x.root.event(x, END, x.activity + " end", d)

	return nil

}

// Info generates an informational log event.
func (x *Task) Info(msg string, data ...interface{}) {
	x.root.event(x, INFO, msg, DAggregate(data).CopyFrom(x.data))
}


// Begin generates a begin log event reporting that a task has
// started.  This is useful to report the start of a long-running
// activity.  The given data is accumulated into the task using
// D.AggregateFrom().  The host Task is passed through as the return.
func (x *Task) Begin(data ...interface{}) *Task {
	x.Time()
	d := DAggregate(data)
	d.CopyFrom(x.data)
	x.root.event(x, BEGIN, x.activity + " begin", d)
	return x
}


// Warning generates a warning log event reporting that a fault was
// encountered but the Task is proceeding acceptably.
func (x *Task) Warning(msg string, data ...interface{}) {
	d := DAggregate(data)
	d.CopyFrom(x.data)
	x.root.event(x, WARNING, msg, d)
}

// Success reports that the Task has concluded successfully.  If the
// task is being timed it will be clocked and the duration reported.
// It always returns nil.  Continuing to use the Task will not cause
// an error but is discouraged.
func (x *Task) Success(data ...interface{}) error {

	x.Clock()

	d := DAggregate(data)
	d.CopyFrom(x.data)

	x.root.event(x, END, x.activity + " success", d)

	return nil

}

// Error reports an unrecoverable fault.  If the Task is being timed
// it will be clocked and the duration reported.  An error is returned
// wrapping the original error with a message reporting that the
// Task's activity has failed.  Continuing to use the Task will not
// cause an error but is discouraged.
func (x *Task) Error(err error, data ...interface{}) error {

	x.Clock()

	m := x.activity + " failed"

	e := wraperror(m, err, data)
	e.Locate(1)

	var d = D{}
	d.CopyFromD(x.data)
	d.Set("Error", err)
	
	x.root.event(x, ERROR, m, d)

	return e

}

// Failure reports an unrecoverable fault.  If the Task is being timed
// it will be clocked and the duration reported.  An error is returned
// reporting that the Task's activity has failed due to the underlying
// cause given in the message.  Continuing to use the Task will not
// cause an error but is discouraged.
//
// Failure and Error are essentially the same, the difference being
// that Failure is useful to both report and generate a fault detected
// directly by the calling code, rather than one caused by an
// underlying error returned from another function or component.
func (x *Task) Failure(msg string, data ...interface{}) error {
	e := newerror(msg, data)
	e.Locate(1)
	return x.Error(e)
}
