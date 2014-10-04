package logberry


type Root struct {
	outputdrivers []OutputDriver
	errorlisteners []ErrorListener

	Tag string
}

type OutputDriver interface {
	Attach(root *Root)
	Detach()

	Report(context *Context,
	  class EventClass,
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
		outputdrivers: make([]OutputDriver, 0, 1),
		errorlisteners: make([]ErrorListener, 0),

		Tag: tag,
	}
}


//----------------------------------------------------------------------
//----------------------------------------------------------------------
func (x *Root) ClearOutputDrivers() {
	for _,o := range(x.outputdrivers) {
		o.Detach()
	}
	x.outputdrivers = make([]OutputDriver, 0, 1)
}

func (x *Root) AddOutputDriver(driver OutputDriver) *Root {
	x.outputdrivers = append(x.outputdrivers, driver)
	driver.Attach(x)
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
func (x *Root) NewContext(class ContextClass, label string) *Context {
	c := NewContext(nil, class, label)
	c.Root = x
	return c
}


//----------------------------------------------------------------------
//----------------------------------------------------------------------
// Report an error that occurred in logging itself.
func (x *Root) InternalError(err error) {
	for _,listener := range(x.errorlisteners) {
		listener.InternalError(err)
	}
	// end logerror
}

/*
 * Internal multiplexer out to all active OutputDrivers.
 */
func (x *Root) Report(context *Context,
	                    class EventClass,
	                    msg string,
	                    data *D) {

	for _,driver := range(x.outputdrivers) {
		driver.Report(context, class, msg, data)
	}

	// end Report
}
