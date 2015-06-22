package logberry

import (
	"time"
	"sync"
)

// An ImmediateRoot executes logging immediately upon event
// generation, in the same goroutine.
type ImmediateRoot struct {
	outputdrivers  []OutputDriver
	errorlisteners []ErrorListener

	outputmx sync.Mutex
}

// NewImmediateRoot creates a new ImmediateRoot.
func NewImmediateRoot() *ImmediateRoot {
	return &ImmediateRoot{
		outputdrivers:  make([]OutputDriver, 0, 1),
		errorlisteners: make([]ErrorListener, 0),
	}
}


// CleaOutputDriver removes all of the currently registered outputs.
func (x *ImmediateRoot) ClearOutputDrivers() Root {

	old := x.outputdrivers

	x.outputdrivers = make([]OutputDriver, 0, 1)

	// Must detach after clearing so the OutputDrivers won't
	// receive output after being detached.
	for _, o := range old {
		o.Detach()
	}

	return x
}

// AddOutputDriver includes the given additional output in those to
// which this ImmediateRoot forwards events.  This is not thread safe with
// event generation, drivers are assumed to be attached at startup.
func (x *ImmediateRoot) AddOutputDriver(driver OutputDriver) Root {

	// Must attach first so that the OutputDriver won't receive output
	// until it knows its root.
	driver.Attach(x) 
	x.outputdrivers = append(x.outputdrivers, driver)
	return x

}

// SetOutputDriver makes the given driver the only output for this
// root.  It is identical to calling x.ClearOutputDrivers() and then
// x.AddOutputDriver(driver).
func (x *ImmediateRoot) SetOutputDriver(driver OutputDriver) Root {
	x.ClearOutputDrivers()
	x.AddOutputDriver(driver)
	return x
}

// ClearErrorListeners removes all of the registered listeners.
func (x *ImmediateRoot) ClearErrorListeners() Root {
	x.errorlisteners = make([]ErrorListener, 0)
	return x
}

// AddErrorListener includes the given listener among those to which
// internal logging errors are reported.
func (x *ImmediateRoot) AddErrorListener(listener ErrorListener) Root {
	x.errorlisteners = append(x.errorlisteners, listener)
	return x
}

// SetErrorListener makes the given listener the only one for this
// ImmediateRoot.  It is identical to calling x.ClearErrorListeners() and then
// x.AddErrorListener(listener).
func (x *ImmediateRoot) SetErrorListener(listener ErrorListener) Root {
	x.ClearErrorListeners()
	x.AddErrorListener(listener)
	return x
}


// Task creates a new top level Task under this ImmediateRoot.
func (x *ImmediateRoot) Task(activity string, data ...interface{}) *Task {
	t := newtask(nil, activity, data)
	t.root = x
	return t
}

// Component creates a new top level Task under this ImmediateRoot.
func (x *ImmediateRoot) Component(component string, data ...interface{}) *Task {
	t := newtask(nil, "Component " + component, data)
	t.SetComponent(component)
	t.root = x
	return t
}


// internalerror reports an internal logging error.  It is generally
// to be used only by OutputDrivers.
func (x *ImmediateRoot) internalerror(err error) {

	x.outputmx.Lock()
	defer x.outputmx.Unlock()

	for _, listener := range x.errorlisteners {
		listener.Error(err)
	}
	// end logerror
}

// event indicates something to report, a log entry to make.  It is
// generally to be used by Tasks.
func (x *ImmediateRoot ) event(task *Task, event string, message string, data D) *Event {

	x.outputmx.Lock()
	defer x.outputmx.Unlock()

	e := &Event{
		TaskID: task.uid,
		Component: task.component,
		Event: event,
		Message: message,
		Data: data,
		
		Timestamp: time.Now(),
	}

	if task.parent != nil {
		e.ParentID = task.parent.uid
	}

	for _, output := range x.outputdrivers {
		output.Event(e)
	}
		
	return e

	// end event
}
