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
func (task *Task) Success(data ...interface{}) {

	task.Data.AggregateInto(data)

	// If this was a longrunning task, compute duration and alter fields
	if task.Long {
		duration := time.Now().Sub(task.Start)

		task.Class = TASK_FINISH
		task.Msg = "@Finish " + task.Msg
		task.Data.Set("Duration", duration)
	}

	LogPrimitive(task.Component.Component, task.Class, task.Msg, task.Data)

	// end Success
}

func (task *Task) Error(err error, data ...interface{}) error {
	task.Data.AggregateInto(data)
	task.Data.Set("Error", err)

	if task.Long {
		task.Msg = "@Error " + task.Msg
	}

	task.Msg += " failed"

	LogPrimitive(task.Component.Component, ERROR, task.Msg, task.Data)

	return WrapError(err, task.Msg)

	// end Error
}

func (task *Task) Failure(msg string, data ...interface{}) error {
	task.Data.AggregateInto(data)

	if task.Long {
		task.Msg = "@Error " + task.Msg
	}

	task.Msg += " failed"

	LogPrimitive(task.Component.Component, ERROR, task.Msg, task.Data)

	return NewError(msg)

	// end Failure
}
