package logberry

import (
	"time"
	"sync"
)

// BackgroundRoots push events to OutputDrivers in a thread safe and
// receipt ordered fashion but in a separate, dedicated goroutine.
// This may be useful for logging outputs that may take some time,
// e.g., pushing to a logging server.
type BackgroundRoot struct {
	outputdrivers  []OutputDriver
	errorlisteners []ErrorListener
	events chan *Event
	wg sync.WaitGroup
}

// NewBackgroundRoot creates a new BackgroundRoot.  The buffer
// parameter indicates the size of the channel buffer connecting event
// generation to outputs.  The goroutine that creates the
// BackgroundRoot should defer a call to Stop() to ensure that all
// events are pushed.
func NewBackgroundRoot(buffer int) Root {
	r := &BackgroundRoot{
		outputdrivers:  make([]OutputDriver, 0, 1),
		errorlisteners: make([]ErrorListener, 0),
		events: make(chan *Event, buffer),
	}

	r.wg.Add(1)
	go r.run()

	return r
}

// Stop shuts down the BackgroundRoot.  Its internal channel is closed, and
// generated log events no longer forwarded to output drivers.
func (x *BackgroundRoot) Stop() {
	close(x.events)
	x.wg.Wait()
}

func (x *BackgroundRoot) run() {

	for {
		e, more := <- x.events
		if !more { break }
		
		for _, driver := range x.outputdrivers {
			driver.Event(e)
		}

	}

	x.wg.Done()

}

// CleaOutputDrivers removes all of the currently registered outputs.
func (x *BackgroundRoot) ClearOutputDrivers() Root {

	old := x.outputdrivers

	x.outputdrivers = make([]OutputDriver, 0, 1)

	// Must detach after clearing so the OutputDrivers won't
	// receive output after being detached.
	for _, o := range old {
		o.Detach()
	}

	return x

}

// AddOutputDrivers includes the given additional output in those to
// which this BackgroundRoot forwards events.  This is not thread safe with
// event generation, drivers are assumed to be attached at startup.
func (x *BackgroundRoot) AddOutputDriver(driver OutputDriver) Root {

	// Must attach first so that the OutputDriver won't receive output
	// until it knows its root.
	driver.Attach(x) 
	x.outputdrivers = append(x.outputdrivers, driver)
	return x

}

// SetOutputDriver makes the given driver the only output for this
// root.  It is identical to calling x.ClearOutputDrivers() and then
// x.AddOutputDriver(driver).
func (x *BackgroundRoot) SetOutputDriver(driver OutputDriver) Root {
	x.ClearOutputDrivers()
	x.AddOutputDriver(driver)
	return x
}

// ClearErrorListeners removes all of the registered elisteners.
func (x *BackgroundRoot) ClearErrorListeners() Root {
	x.errorlisteners = make([]ErrorListener, 0)
	return x
}

// AddErrorListener includes the given listener among those to which
// internal logging errors are reported.
func (x *BackgroundRoot) AddErrorListener(listener ErrorListener) Root {
	x.errorlisteners = append(x.errorlisteners, listener)
	return x
}

// SetErrorListener makes the given listener the only one for this
// BackgroundRoot.  It is identical to calling x.ClearErrorListeners() and then
// x.AddErrorListener(listener).
func (x *BackgroundRoot) SetErrorListener(listener ErrorListener) Root {
	x.ClearErrorListeners()
	x.AddErrorListener(listener)
	return x
}


// NewTask creates a new top level Task under this BackgroundRoot.
func (x *BackgroundRoot) Task(activity string, data ...interface{}) *Task {
	t := newtask(nil, activity, data)
	t.root = x
	return t
}

// NewTask creates a new top level Task under this BackgroundRoot.
func (x *BackgroundRoot) Component(component string, data ...interface{}) *Task {
	t := newtask(nil, component, data)
	t.SetComponent(component)
	t.root = x
	return t
}


// internalerror reports an internal logging error.  It is generally
// to be used only by OutputDrivers.
func (x *BackgroundRoot) internalerror(err error) {
	for _, listener := range x.errorlisteners {
		listener.Error(err)
	}
	// end logerror
}

// event indicates something to report, a log entry to make.  It is
// generally to be used by Tasks.
func (x *BackgroundRoot ) event(task *Task, event string, message string, data D) *Event {

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