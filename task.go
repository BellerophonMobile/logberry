package logberry

import (
	"time"
)

type Task struct {

	Component *ComponentLog
	Class StatementClass
	Msg string
	Data *D

	Long bool
	Start time.Time

}


//----------------------------------------------------------------------
//----------------------------------------------------------------------
func (x *ComponentLog) LongTask(msg string, data ...interface{}) *Task {

	d := DAggregate(data)

	LogPrimitive(x.Component, TASK_START, "@Start " + msg, d)

	return &Task{
		Component: x,
		Class: TASK_START,
		Msg: msg,
		Data: d,

		Long: true,
		Start: time.Now(),
	}

	// end LongTask
}

func (x *ComponentLog) Task(msg string, data ...interface{}) *Task {

	d := DAggregate(data)

	return &Task{
		Component: x,
		Class: TASK,
		Msg: msg,
		Data: d,
	}

	// end Task
}

func (x *ComponentLog) ResourceTask(msg string, resource interface{}) *Task {

	return &Task{
		Component: x,
		Class: RESOURCE,
		Msg: msg,
		Data: &D{"Resource": resource},
	}

	// end ResourceTask
}

func (x *ComponentLog) ServiceTask(msg string, service interface{}) *Task {

	return &Task{
		Component: x,
		Class: SERVICE,
		Msg: msg,
		Data: &D{"Service": service},
	}

	// end ServiceTask
}


//----------------------------------------------------------------------
//----------------------------------------------------------------------
func (task *Task) AddField(k string, v interface{}) {
	(*task.Data)[k] = v
}

func (task *Task) AddData(data ...interface{}) {
	task.Data.AggregateInto(data)
}


//----------------------------------------------------------------------
//----------------------------------------------------------------------
func (task *Task) Success(data ...interface{}) {

	task.Data.AggregateInto(data)

	var msg = task.Msg

	// If this was a long running task, compute duration and alter fields
	if task.Long {
		duration := time.Now().Sub(task.Start)

		task.Class = TASK_FINISH
		msg = "@Success " + task.Msg
		task.Data.Set("Duration", duration)
	}

	LogPrimitive(task.Component.Component, task.Class, msg, task.Data)

	// end Success
}

func (task *Task) Warning(w string, data ...interface{}) {
	task.Data.AggregateInto(data)
	task.Data.Set("Warning", w)

	var msg = task.Msg + " warning"
	if task.Long {
		msg = "@Warning " + msg
	}

	LogPrimitive(task.Component.Component, WARNING, msg, task.Data)

	// end Warning
}

func (task *Task) Error(err error, data ...interface{}) error {
	task.Data.AggregateInto(data)
	task.Data.Set("Error", err.Error())

	var msg = task.Msg + " failed"
	if task.Long {
		msg = "@Error " + msg
	}

	LogPrimitive(task.Component.Component, ERROR, msg, task.Data)

	return WrapError(err, task.Msg)

	// end Error
}

func (task *Task) Failure(f string, data ...interface{}) error {
	task.Data.AggregateInto(data)
	task.Data.Set("Error", f)

	var msg = task.Msg + " failed"
	if task.Long {
		msg = "@Error " + msg
	}

	LogPrimitive(task.Component.Component, ERROR, msg, task.Data)

	return NewError(f)

	// end Failure
}
