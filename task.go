package logberry

import (
	"time"
)

type Task struct {
	UID    uint64
	Parent Context
	Root   *Root
	Label  string

	Activity string

	Timed bool
	Start time.Time

	Data *D

	mute bool
	highlight bool
}

func newtask(parent Context, activity string, data []interface{}) *Task {

	t := &Task{
		UID:    newcontextuid(),
		Parent: parent,
		Root:   parent.GetRoot(),
		Label:  parent.GetLabel(),

		Activity: activity,

		Data: DAggregate(data),
	}

	t.Data.Set(t.Root.FieldPrefix+"Parent", t.Parent.GetUID())

	return t

}

//----------------------------------------------------------------------
func (x *Task) Component(label string, data ...interface{}) *Component {
	return newcomponent(x, label, data...)
}

func (x *Task) Task(activity string, data ...interface{}) *Task {
	return newtask(x, activity, data)
}

//----------------------------------------------------------------------
func (x *Task) GetLabel() string {
	return x.Label
}

func (x *Task) GetUID() uint64 {
	return x.UID
}

func (x *Task) GetParent() Context {
	return x.Parent
}

func (x *Task) GetRoot() *Root {
	return x.Root
}

func (x *Task) Time() *Task {
	x.Timed = true
	x.Start = time.Now()
	return x
}

func (x *Task) Clock() time.Duration {

	if !x.Timed {
		return 0
	}

	d := time.Now().Sub(x.Start)
	x.Data.Set(x.Root.FieldPrefix+"Duration", d)
	return d

}

func (x *Task) AddData(k string, v interface{}) *D {
	(*x.Data)[k] = v
	return x.Data
}

func (x *Task) AggregateData(data ...interface{}) *D {
	x.Data.AggregateFrom(data)
	return x.Data
}


func (x *Task) Mute() *Task {
	x.mute = true
	return x
}
func (x *Task) Unmute() *Task {
	x.mute = false
	return x
}
func (x *Task) IsMute() bool {
	return x.mute
}

func (x *Task) Highlight() *Task {
	x.highlight = true
	return x
}

func (x *Task) ClearHighlight() *Task {
	x.highlight = false
	return x
}

func (x *Task) IsHighlighted() bool {
	return x.highlight
}

func (x *Task) Calculation(calculation interface{}) *Task {
	x.AddData("Calculation", calculation)
	return x
}

func (x *Task) File(file interface{}) *Task {
	x.AddData("File", file)
	return x
}

func (x *Task) Resource(resource interface{}) *Task {
	x.AddData("Resource", resource)
	return x
}

func (x *Task) Service(service interface{}) *Task {
	x.AddData("Service", service)
	return x
}

func (x *Task) User(user interface{}) *Task {
	x.AddData("User", user)
	return x
}

func (x *Task) Endpoint(endpoint interface{}) *Task {
	x.AddData("Endpoint", endpoint)
	return x
}

// Always returns nil.
func (x *Task) Success(data ...interface{}) error {

	x.Clock()
	x.Data.AggregateFrom(data)
	x.Root.TaskEvent(x, TASK_END)

	return nil

}

func (x *Task) Terminated(msg string, data ...interface{}) error {
	x.Clock()
	x.Data.AggregateFrom(data)
	x.Data.Set(x.Root.FieldPrefix+"Warning", msg)
	x.Root.TaskEvent(x, TASK_WARNING)
	return nil
}

func (x *Task) Error(err error, data ...interface{}) error {

	// Note that this can't just throw err into the data blob because
	// the standard errors package error interface doesn't expose much,
	// even the message, so the marshalers don't get anything in that
	// common case.  Hence the reduction to a string via Error().

	x.Clock()
	x.Data.AggregateFrom(data)
	x.Data.Set(x.Root.FieldPrefix+"Error", err.Error())

	x.Root.TaskEvent(x, TASK_ERROR)

	return WrapError(err, x.Activity+" failed")

}

// Failure is the same as Error but doesn't take an error object.
func (x *Task) Failure(msg string, data ...interface{}) error {

	x.Clock()
	x.Data.AggregateFrom(data)
	x.Data.Set(x.Root.FieldPrefix+"Error", msg)
	x.Root.TaskEvent(x, TASK_ERROR)

	return WrapError(NewError(msg), x.Activity+" failed")

}

func (x *Task) Begin(data ...interface{}) *Task {
	x.Data.AggregateFrom(data)
	x.Root.TaskEvent(x, TASK_BEGIN)
	return x
}

// Unlike the terminal events, this does not accumulate the given data
// into the Task.  However, you may replicate that behavior
// (aggregating & reporting all of the accumulated data so far) by:
//
//   foo.Info("Status report", foo.AggregateData("mushi", "sushi"))
//
// foo.AddData() may be used similarly.
func (x *Task) Info(msg string, data ...interface{}) {
	x.Root.TaskProgress(x, TASK_INFO, msg, DAggregate(data))
}

// Unlike the terminal events, this does not accumulate the given data
// into the Task.  However, you may replicate that behavior
// (aggregating & reporting all of the accumulated data so far) by:
//
//   foo.Info("Status report", foo.AggregateData("mushi", "sushi"))
//
// foo.AddData() may be used similarly.
func (x *Task) Warning(msg string, data ...interface{}) {
	x.Root.TaskProgress(x, TASK_WARNING, msg, DAggregate(data))
}
