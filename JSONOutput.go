package logberry

import (
	"encoding/json"
	"io"
)


// JSONOutput is an OutputDriver that writes log events in JSON.
type JSONOutput struct {
	root   Root
	writer io.Writer
}

// NewJSONOutput creates a new JSONOutput targeted at the given
// Writer.  DifferentialTime defaults to false.
func NewJSONOutput(w io.Writer) *JSONOutput {
	return &JSONOutput{
		writer:           w,
	}
}

// Attach notifies the OutputDriver of its Root.  It should only be
// called by a Root.
func (x *JSONOutput) Attach(root Root) {
	x.root = root
}

// Detach notifies the OutputDriver that it has been removed from its
// Root.  It should only be called by a root.
func (x *JSONOutput) Detach() {
	x.root = nil
}


// Event outputs a generated log entry, as called by a Root or a
// chaining OutputDriver.
func (x *JSONOutput) Event(event *Event) {

	var bytes []byte
	var err error
	bytes, err = json.Marshal(event)
	if err != nil {
		x.root.internalerror(WrapError("Could not marshal log entry", err))
		return
	}

	_,e := x.writer.Write(bytes)
	if e != nil {
		x.root.internalerror(WrapError("Could not write entry", e))
		return
	}
	
	_,e = x.writer.Write([]byte("\n"))
	if e != nil {
		x.root.internalerror(WrapError("Could not write entry", e))
		return
	}

	// end Event
}
