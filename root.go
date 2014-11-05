package logberry

type Root struct {
	outputdrivers  []OutputDriver
	errorlisteners []ErrorListener

	Tag string

	FieldPrefix string
}

type OutputDriver interface {
	Attach(root *Root)
	Detach()

	ComponentEvent(component *Component,
		event ComponentEventClass,
		msg string,
		data *D)

	TaskEvent(task *Task,
		event TaskEventClass)

	TaskProgress(task *Task,
		event TaskEventClass,
		msg string,
		data *D)

	//	Action(action Action)
}

type ErrorListener interface {
	InternalError(err error)
}

//----------------------------------------------------------------------
//----------------------------------------------------------------------
func NewRoot(tag string) *Root {
	return &Root{
		outputdrivers:  make([]OutputDriver, 0, 1),
		errorlisteners: make([]ErrorListener, 0),

		Tag: tag,
	}
}

//----------------------------------------------------------------------
//----------------------------------------------------------------------
func (x *Root) ClearOutputDrivers() {
	x.outputdrivers = make([]OutputDriver, 0, 1)
	for _, o := range x.outputdrivers {
		o.Detach() // Must be after clearing so the OutputDrivers won't
		// receive output after being detached.
	}
}

func (x *Root) AddOutputDriver(driver OutputDriver) *Root {
	driver.Attach(x) // Must be first so that the OutputDriver won't
	// receive output until it knows its root.
	x.outputdrivers = append(x.outputdrivers, driver)
	return x
}

// Is identical to calling x.ClearOutputDrivers() and then
// x.AddOutputDriver(driver).
func (x *Root) SetOutputDriver(driver OutputDriver) *Root {
	x.ClearOutputDrivers()
	x.AddOutputDriver(driver)
	return x
}

//----------------------------------------------------------------------
//----------------------------------------------------------------------
func (x *Root) ClearErrorListeners() *Root {
	x.errorlisteners = make([]ErrorListener, 0)
	return x
}

func (x *Root) AddErrorListener(listener ErrorListener) *Root {
	x.errorlisteners = append(x.errorlisteners, listener)
	return x
}

// Is identical to calling x.ClearErrorListeners() and then
// x.AddErrorListener(listener).
func (x *Root) SetErrorListener(listener ErrorListener) *Root {
	x.ClearErrorListeners()
	x.AddErrorListener(listener)
	return x
}

//----------------------------------------------------------------------
//----------------------------------------------------------------------
func (x *Root) NewComponent(label string, data ...interface{}) *Component {
	c := newcomponent(nil, label, DAggregate(data))
	c.root = x
	return c
}

//----------------------------------------------------------------------
//----------------------------------------------------------------------
// Report an error that occurred in logging itself.
func (x *Root) InternalError(err error) {
	for _, listener := range x.errorlisteners {
		listener.InternalError(err)
	}
	// end logerror
}

/*
 * Internal fan-out to all active OutputDrivers.
 */
func (x *Root) ComponentEvent(component *Component,
	event ComponentEventClass,
	msg string,
	data *D) {

	// Root doesn't check that event is within range because the output
	// drivers need to actually report the error anyway.

	for _, driver := range x.outputdrivers {
		driver.ComponentEvent(component, event, msg, data)
	}

	// end ComponentEvent
}

/*
 * Internal fan-out to all active OutputDrivers.
 */
func (x *Root) TaskEvent(task *Task,
	event TaskEventClass) {

	// Root doesn't check that event is within range because the output
	// drivers need to actually report the error anyway.

	for _, driver := range x.outputdrivers {
		driver.TaskEvent(task, event)
	}

	// end TaskEvent
}

/*
 * Internal fan-out to all active OutputDrivers.
 */
func (x *Root) TaskProgress(task *Task,
	event TaskEventClass,
	msg string,
	data *D) {

	// Root doesn't check that event is within range because the output
	// drivers need to actually report the error anyway.

	for _, driver := range x.outputdrivers {
		driver.TaskProgress(task, event, msg, data)
	}

	// end TaskEvent
}
