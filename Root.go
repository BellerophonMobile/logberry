package logberry

import (
	"time"
)

// Roots interface between event generation and output.  Tasks are
// created under a root for generating log events.  OutputDrivers are
// attached to Roots for receiving those events.  All attached
// OutputDrivers receive each event, in a thread safe and receipt
// ordered fashion.  Internal logging errors, e.g., failures to write
// to disk, may be captured via attached ErrorListeners.
type Root struct {
	outputdrivers  []OutputDriver
	errorlisteners []ErrorListener
	events chan *Event
}

// NewRoot creates a new Root.  The buffer parameter indicates whether
// or not the channel between these is buffered or not.  Small
// programs should generally not buffer (input 0) so that all output
// is processed before exit.  Larger programs may want to buffer so
// that logging may be managed on a separate goroutine, e.g., to
// forward to a log server.
func NewRoot(buffer int) *Root {
	r := &Root{
		outputdrivers:  make([]OutputDriver, 0, 1),
		errorlisteners: make([]ErrorListener, 0),
		events: make(chan *Event, buffer),
	}

	go r.run()

	return r
}

// Stop shuts down the Root.  Its internal channel is closed, and
// generated log events no longer forwarded to output drivers.
func (x *Root) Stop() {
	close(x.events)
}

func (x *Root) run() {

	for {
		e, more := <- x.events
		if !more { break }
		
		for _, driver := range x.outputdrivers {
			driver.Event(e)
		}

	}
	
}

// CleaOutputDrivers removes all of the currently registered outputs.
func (x *Root) ClearOutputDrivers() {

	old := x.outputdrivers

	x.outputdrivers = make([]OutputDriver, 0, 1)

	// Must detach after clearing so the OutputDrivers won't
	// receive output after being detached.
	for _, o := range old {
		o.Detach()
	}
	
}

// AddOutputDrivers includes the given additional output in those to
// which this Root forwards events.  This is not thread safe with
// event generation, drivers are assumed to be attached at startup.
func (x *Root) AddOutputDriver(driver OutputDriver) *Root {

	// Must attach first so that the OutputDriver won't receive output
	// until it knows its root.
	driver.Attach(x) 
	x.outputdrivers = append(x.outputdrivers, driver)
	return x

}

// SetOutputDriver makes the given driver the only output for this
// root.  It is identical to calling x.ClearOutputDrivers() and then
// x.AddOutputDriver(driver).
func (x *Root) SetOutputDriver(driver OutputDriver) *Root {
	x.ClearOutputDrivers()
	x.AddOutputDriver(driver)
	return x
}

// ClearErrorListeners removes all of the registered elisteners.
func (x *Root) ClearErrorListeners() *Root {
	x.errorlisteners = make([]ErrorListener, 0)
	return x
}

// AddErrorListener includes the given listener among those to which
// internal logging errors are reported.
func (x *Root) AddErrorListener(listener ErrorListener) *Root {
	x.errorlisteners = append(x.errorlisteners, listener)
	return x
}

// SetErrorListener makes the given listener the only one for this
// Root.  It is identical to calling x.ClearErrorListeners() and then
// x.AddErrorListener(listener).
func (x *Root) SetErrorListener(listener ErrorListener) *Root {
	x.ClearErrorListeners()
	x.AddErrorListener(listener)
	return x
}


// NewTask creates a new top level Task under this Root.
func (x *Root) Task(activity string, data ...interface{}) *Task {
	t := newtask(nil, activity, data)
	t.root = x
	return t
}

// NewTask creates a new top level Task under this Root.
func (x *Root) Component(component string, data ...interface{}) *Task {
	t := newtask(nil, component, data)
	t.SetComponent(component)
	t.root = x
	return t
}


// internalerror reports an internal logging error.  It is generally
// to be used only by OutputDrivers.
func (x *Root) internalerror(err error) {
	for _, listener := range x.errorlisteners {
		listener.Error(err)
	}
	// end logerror
}

// event indicates something to report, a log entry to make.  It is
// generally to be used by Tasks.
func (x *Root) event(task *Task, event string, message string, data D) *Event {

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
	
  x.events <- e
	return e

	// end event
}
