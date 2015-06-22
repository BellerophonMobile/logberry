package logberry

// An OutputDriver is registered to Roots and receives log events to
// export, e.g., writing to disk, screen, or sending to a server.  To
// do so, an OutputDriver is created and then passed to the
// AddOutputDriver function of a Root.  That Root will then call the
// OutputDriver's Attach() function to notify it of its context.
// Unless specifically noted otherwise by the implementation, it is an
// error with unspecified behavior to add an OutputDriver instance to
// more than one Root simultaneously.
type OutputDriver interface {

	Attach(root Root)
	Detach()

	Event(event *Event)

}
