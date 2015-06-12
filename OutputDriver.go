package logberry

// OutputDrivers receive log events and export them, e.g., writing to
// disk, screen, or sending to a server.  To do so, an OutputDriver is
// created and then passed to the Root.AddOutputDriver() function.
// That will then call the OutputDriver's Attach() function to notify
// it of its Root.  It is an error with unspecified behavior to add an
// OutputDriver to more than one Root simultaneously.
type OutputDriver interface {

	Attach(root Root)
	Detach()

	Event(event *Event)

}
