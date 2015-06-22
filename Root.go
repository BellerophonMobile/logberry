package logberry

// A Root interfaces between event generation and output.  Tasks are
// created under a root for generating log events.  OutputDrivers are
// attached to Roots for receiving those events.  All attached
// OutputDrivers receive each event, in a thread safe and receipt
// ordered fashion.  Internal logging errors, e.g., failures to write
// to disk, may be captured via attached ErrorListeners.
type Root interface {
	ClearOutputDrivers() Root
	AddOutputDriver(driver OutputDriver) Root
	SetOutputDriver(driver OutputDriver) Root

	ClearErrorListeners() Root
	AddErrorListener(listener ErrorListener) Root
	SetErrorListener(listener ErrorListener) Root

	Task(activity string, data ...interface{}) *Task
	Component(component string, data ...interface{}) *Task

	internalerror(err error)

	event(task *Task, event string, message string, data D) *Event
}
