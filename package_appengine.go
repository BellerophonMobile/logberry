// +build appengine

package logberry

import (
	"appengine"
)

var Std *Root
var Main *Task

// NewRootTask creates a new root, defaulting to an AppEngineOutput driver,
// logging to a given task.  This is necessary, since any output an App Engine
// program does needs to be done using a context from a request.
func NewRootTask(ctx appengine.Context) (*Root, *Task) {
	root := NewRoot(24)
	root.AddOutputDriver(NewAppEngineOutput(ctx))

	return root, &Task{
		component: "main",
		activity:  "Component main",
		root:      root,
	}
}
