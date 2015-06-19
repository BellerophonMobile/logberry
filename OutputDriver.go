package logberry

// OutputDrivers are registered to Roots and receive log events to
// export, e.g., writing to disk, screen, or sending to a server.  To
// do so, an OutputDriver is created and then passed to the
// AddOutputDriver function of a Root.  That Root will then call the
// OutputDriver's Attach() function to notify it of its context.  It
// is an error with unspecified behavior to add an OutputDriver to
// more than one Root simultaneously.
type OutputDriver interface {

	Attach(root Root)
	Detach()

	Event(event *Event)

}
