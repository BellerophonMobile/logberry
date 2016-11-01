package logberry

type FilterOutput struct {
	root   *Root	
	label string
	outputdriver  OutputDriver
}

func NewFilterOutput(bytes []byte, label string, output OutputDriver) (*FilterOutput,error) {
	f := &FilterOutput{
		label: label,
		outputdriver:  output,
	}

	return f,nil

}


// Attach notifies the OutputDriver of its Root.  It should only be
// called by a Root.
func (x *FilterOutput) Attach(root *Root) {
	x.root = root
}

// Detach notifies the OutputDriver that it has been removed from its
// Root.  It should only be called by a root.
func (x *FilterOutput) Detach() {
	x.root = nil
}

// Event outputs a generated log entry, as called by a Root or a
// chaining OutputDriver.
func (x *FilterOutput) Event(event *Event) {
	x.outputdriver.Event(event)
}
