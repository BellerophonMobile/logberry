package logberry

import (
	"time"
)


type Task struct {
	UID uint64
	Parent Context
	Root *Root
	Label string

	Activity string
	Class ActivityClass

	Timed bool
	Start time.Time

	Data *D
}

const (
	LONG = true
	SHORT = false
)


//----------------------------------------------------------------------
//----------------------------------------------------------------------
func newtask(parent Context, long bool, activity string, data []interface{}) *Task {

	var class = APPLICATION
	if data != nil && len(data) > 0 {
		if ac,ok := data[0].(ActivityClass); ok {
			class = ac
			data = data[1:]
		}
	}

	t := &Task {
		UID: newcontextuid(),
		Parent: parent,
		Root: parent.GetRoot(),
		Label: parent.GetLabel(),

		Activity: activity,
		Class: class,

		Data: DAggregate(data),
	}

	if t.Class < 0 || t.Class >= activityclasssentinel {
		t.Root.InternalError(NewError("ActivityClass out of range", t.UID, t.Class))
		t.Data.Set(t.Root.FieldPrefix+"Class", t.Class)
	} else {
		t.Data.Set(t.Root.FieldPrefix+"Class", ActivityClassText[t.Class])
	}

	t.Data.Set(t.Root.FieldPrefix+"Parent", t.Parent.GetUID())

	if long {
		t.Timed = true
		t.Root.TaskEvent(t, START, t.Data)
		t.Start = time.Now()
	}

	return t

}

func calculationtask(parent Context, long bool, activity string, calculation interface{}, data ...interface{}) *Task {
	return newtask(parent, long, activity, 
		append([]interface{} {
		CALCULATION,
		&D{
			parent.GetRoot().FieldPrefix+"Calculation": calculation,
	  } },
		data...))
}

func resourcetask(parent Context, long bool, activity string, resource interface{}, data ...interface{}) *Task {
	return newtask(parent, long, activity,
		append([]interface{} {
		RESOURCE,
		&D {
			parent.GetRoot().FieldPrefix+"Resource": resource,
		} },
		data...))
}

func servicetask(parent Context, long bool, activity string, service interface{}, query interface{}, data ...interface{}) *Task {
	return newtask(parent, long, activity,
		append([]interface{} { 
		SERVICE,
		&D{
		  parent.GetRoot().FieldPrefix+"Service": service,
		  parent.GetRoot().FieldPrefix+"Query": query,
	  } },
		data...))
}

//----------------------------------------------------------------------
func (x *Task) Component(label string, data ...interface{}) *Component {
	return newcomponent(x, label, data...)
}

func (x *Task) Task(activity string, data ...interface{}) *Task {
	return newtask(x, false, activity, data)
}
func (x *Task) LongTask(activity string, data ...interface{}) *Task {
	return newtask(x, true, activity, data)
}

func (x *Task) CalculationTask(activity string, calculation interface{}, data ...interface{}) *Task {
	return calculationtask(x, false, activity, calculation, data...)
}
func (x *Task) LongCalculationTask(activity string, calculation interface{}, data ...interface{}) *Task {
	return calculationtask(x, true, activity, calculation, data...)
}

func (x *Task) ResourceTask(activity string, resource interface{}, data ...interface{}) *Task {
	return resourcetask(x, false, activity, resource, data...)
}
func (x *Task) LongResourceTask(activity string, resource interface{}, data ...interface{}) *Task {
	return resourcetask(x, true, activity, resource, data...)
}

func (x *Task) ServiceTask(activity string, service interface{}, query interface{}, data ...interface{}) *Task {
	return servicetask(x, false, activity, service, query, data...)
}
func (x *Task) LongServiceTask(activity string, service interface{}, query interface{}, data ...interface{}) *Task {
	return servicetask(x, true, activity, service, query, data...)
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

func (x *Task) Time() {
  x.Timed = true
  x.Start = time.Now()
}

func (x *Task) Clock() time.Duration {

  if !x.Timed {
    return 0
  }

  d := time.Now().Sub(x.Start)
  x.Data.Set(x.Root.FieldPrefix+"Duration", d)
  return d

}

func (x *Task) SetData(k string, v interface{}) {
  (*x.Data)[k] = v
}

func (x *Task) AggregateData(data ...interface{}) {
  x.Data.AggregateFrom(data)
}


//----------------------------------------------------------------------
//----------------------------------------------------------------------
func (x *Task) Finished(data ...interface{}) {

	x.Clock()
	x.Data.AggregateFrom(data)
	x.Root.TaskEvent(x, FINISH, x.Data)

}

func (x *Task) Success(data ...interface{}) {

	x.Clock()
	x.Data.AggregateFrom(data)
	x.Root.TaskEvent(x, SUCCESS, x.Data)

}

func (x *Task) Error(err error, data ...interface{}) error {

	// Note that this can't just throw err into the data blob because
	// the standard errors package error interface doesn't expose much,
	// even the message, so the marshalers don't get anything in that
	// common case.  Hence the reduction to a string via Error().

	x.Clock()
	x.Data.AggregateFrom(data)
	x.Data.Set(x.Root.FieldPrefix+"Error", err.Error())

	x.Root.TaskEvent(x, ERROR, x.Data)

	return WrapError(err, x.Activity + " failed")

}

// Failure is the same as Error but doesn't take an error object.
func (x *Task) Failure(msg string, data ...interface{}) error {

	x.Clock()
	x.Data.AggregateFrom(data)
	x.Data.Set(x.Root.FieldPrefix+"Error", msg)
	x.Root.TaskEvent(x, ERROR, x.Data)

	return WrapError(NewError(msg), x.Activity + " failed")

}


//----------------------------------------------------------------------
//----------------------------------------------------------------------
/*
func (x *Task) Info(msg string, data ...interface{}) {
	x.Root.TaskEvent(x, INFO, msg, DAggregate(data))
}

func (x *Task) Warning(msg string, data ...interface{}) {
	x.Root.TaskEvent(x, WARNING, msg, DAggregate(data))
}
*/
