package logberry


type MultiplexerOutput struct {
	root *Root
	drivers []OutputDriver
}

//----------------------------------------------------------------------
//----------------------------------------------------------------------
// MultiplexerOutput is not necessary for basic multiple output.
// logberry's core supports multiple outputs.  But the multiplexer
// permits arrangements such as feeding the core output to a
// ThreadSafeOutput, which in turns feeds to a multiplexer, and then
// through that to several outputs.
func NewMultiplexerOutput() *MultiplexerOutput {
	return &MultiplexerOutput {
		drivers: make([]OutputDriver, 0),
	}
}

func (x *MultiplexerOutput) AddOutputDriver(out OutputDriver) {
	x.drivers = append(x.drivers, out)
}


//----------------------------------------------------------------------
func (x *MultiplexerOutput) Attach(root *Root) {
	x.root = root
	for _,out := range(x.drivers) {
		out.Attach(root)
	}
}

func (x *MultiplexerOutput) Detach() {
	for _,out := range(x.drivers) {
		out.Detach()
	}
	x.root = nil
}


//----------------------------------------------------------------------
func (x *MultiplexerOutput) ComponentEvent(component *Component,
  class ContextEventClass,
  msg string,
  data *D) {

	for _,out := range(x.drivers) {
		out.ComponentEvent(component, class, msg, data)
	}

}

func (x *MultiplexerOutput) TaskEvent(task *Task,
	event ContextEventClass) {

	for _,out := range(x.drivers) {
		out.TaskEvent(task, event)
	}

}
