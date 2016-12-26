// +build appengine

package logberry

import (
	"appengine"
)

// AppEngineOutput is an OutputDriver that writes out log events in a more or
// less human readable form, using an App Engine context from a request.
type AppEngineOutput struct {
	root *Root
	ctx  appengine.Context
}

// NewAppEngineOutput creates a new AppEngineOutput logging to the given
// appengine context.
func NewAppEngineOutput(ctx appengine.Context) *AppEngineOutput {
	return &AppEngineOutput{ctx: ctx}
}

// Attach notifies the OutputDriver of its Root.  It should only be
// called by a Root.
func (d *AppEngineOutput) Attach(root *Root) {
	d.root = root
}

// Detach notifies the OutputDriver that it has been removed from its
// Root.  It should only be called by a root.
func (d *AppEngineOutput) Detach() {
	d.root = nil
}

// Event outputs a generated log entry, as called by a Root or a
// chaining OutputDriver.
func (d *AppEngineOutput) Event(evt *Event) {
	if d.ctx == nil {
		d.root.InternalError(NewError("no AppEngine context attached to output driver"))
		return
	}

	var logFunc func(string, ...interface{})

	switch evt.Event {
	case WARNING:
		logFunc = d.ctx.Warningf

	case ERROR:
		logFunc = d.ctx.Errorf

	default:
		logFunc = d.ctx.Infof
	}

	logFunc("%v %v (%v:%v): %v %v", evt.Event, evt.Component, evt.TaskID,
		evt.ParentID, evt.Message, evt.Data.String())
}
