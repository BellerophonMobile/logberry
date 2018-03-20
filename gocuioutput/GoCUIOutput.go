package gocuioutput

import (
	"github.com/BellerophonMobile/logberry"
	"github.com/BellerophonMobile/gocui"
)

type GoCUIOutput struct {
	root   *logberry.Root
	gui    *gocui.Gui
	output logberry.OutputDriver
}

func New(g *gocui.Gui, output logberry.OutputDriver) *GoCUIOutput {
	return &GoCUIOutput{
		gui: g,
		output: output,
	}
}

// Attach notifies the OutputDriver of its Root.  It should only be
// called by a Root.
func (o *GoCUIOutput) Attach(root *logberry.Root) {
	o.root = root
}

// Detach notifies the OutputDriver that it has been removed from its
// Root.  It should only be called by a root.
func (o *GoCUIOutput) Detach() {
	o.root = nil
}

// Event outputs a generated log entry, as called by a Root or a
// chaining OutputDriver.
func (o *GoCUIOutput) Event(event *logberry.Event) {

	o.gui.Update(func(g *gocui.Gui) error {
		o.output.Event(event)
		return nil
	})

	// end Event
}
