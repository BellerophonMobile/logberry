package logberry

import (
	"log"

	"bytes"
	"fmt"
	"io"
	"os"
	"syscall"
	"time"

	"reflect"

	"github.com/BellerophonMobile/logberry/terminal"
)

//----------------------------------------------------------------------
//----------------------------------------------------------------------
type TextOutput struct {
	root   *Root
	writer io.Writer

	Start            time.Time
	DifferentialTime bool

	Color bool

	IDOffset   int
	DataOffset int
}

const (
	BLACK int = iota
	RED
	GREEN
	YELLOW
	BLUE
	MAGENTA
	CYAN
	WHITE
)

const (
	HIGH_INTENSITY int = 90
	LOW_INTENSITY  int = 30
)

type TerminalStyle struct {
	color     int
	bold      bool
	intensity int
}

var ComponentEventTerminalStyles = [...]TerminalStyle{
	{BLACK, false, HIGH_INTENSITY},  // start
	{BLACK, false, HIGH_INTENSITY},  // finish
	{BLUE, false, LOW_INTENSITY},    // configuration
	{GREEN, true, HIGH_INTENSITY},   // ready
	{WHITE, false, HIGH_INTENSITY},  // info
	{YELLOW, false, HIGH_INTENSITY}, // warning
	{RED, true, HIGH_INTENSITY},     // error
	{RED, true, HIGH_INTENSITY},     // fatal
}

var TaskEventTerminalStyles = [...]TerminalStyle{
	{WHITE, false, HIGH_INTENSITY},  // begin
	{WHITE, false, HIGH_INTENSITY},  // end
	{WHITE, false, LOW_INTENSITY},   // info
	{YELLOW, false, HIGH_INTENSITY}, // warning
	{RED, true, HIGH_INTENSITY},     // error
}

func init() {

	//-- Check that labels are defined for the enumerations
	if len(ComponentEventTerminalStyles) != int(componenteventclasssentinel) {
		log.Fatal("Fatal internal error: " +
			"len(ComponentEventTerminalStyles) != |ComponentEventClass|")
	}

	if len(TaskEventTerminalStyles) != int(taskeventclasssentinel) {
		log.Fatal("Fatal internal error: " +
			"len(TaskEventTerminalStyles) != |TaskEventClass|")
	}

}

//----------------------------------------------------------------------
//----------------------------------------------------------------------
func NewStdOutput() *TextOutput {
	t := NewTextOutput(os.Stdout)
	t.Color = terminal.IsTerminal(syscall.Stdout)
	return t
}

func NewErrOutput() *TextOutput {
	t := NewTextOutput(os.Stderr)
	t.Color = terminal.IsTerminal(syscall.Stderr)
	return t
}

func NewTextOutput(w io.Writer) *TextOutput {
	return &TextOutput{
		writer:     w,
		Start:      time.Now(),
		IDOffset:   84,
		DataOffset: 100,
	}
}

//------------------------------------------------------
func (o *TextOutput) Attach(root *Root) {
	o.root = root
}

func (o *TextOutput) Detach() {
	o.root = nil
}

//----------------------------------------------------------------------
//----------------------------------------------------------------------
func (o *TextOutput) timestamp() string {

	if o.DifferentialTime {
		return time.Since(o.Start).String()
	}

	return time.Now().Format(time.RFC3339)

	// end timestamp
}

func (o *TextOutput) internalerror(err error) {

	if o.Color {
		fmt.Fprintf(o.writer, "\x1b[%d;1m", HIGH_INTENSITY+RED)
	}

	fmt.Fprintf(o.writer, "%v [LOG ERROR] %v\n",
		o.timestamp(), err.Error())

	if o.Color {
		fmt.Fprintf(o.writer, "\x1b[0m")
	}

	fmt.Fprintf(o.writer, "\n")

	o.root.InternalError(WrapError(err, "Could not output log entry"))

}

// Convenience function to handle the potential write error.
func (o *TextOutput) printf(msg string, a ...interface{}) int {
	n, e := fmt.Fprintf(o.writer, msg, a...)
	if e != nil {
		o.internalerror(e)
	}
	return n
}

func keyrenderrecurse(bytes *bytes.Buffer, wrap bool, data interface{}) {

	var val = reflect.ValueOf(data)

	switch val.Kind() {

	case reflect.Interface:
		fallthrough
	case reflect.Ptr:
		keyrenderrecurse(bytes, wrap, val.Elem().Interface())

	case reflect.Map:
		var vals = val.MapKeys()
		if wrap {
			fmt.Fprint(bytes, "{ ")
		}
		for _, k := range vals {
			fmt.Fprintf(bytes, "%s=", k.Interface())
			keyrenderrecurse(bytes, true, val.MapIndex(k).Interface())
		}
		if wrap {
			fmt.Fprint(bytes, "}")
		}

	case reflect.Struct:
		var vtype = val.Type()

		if wrap {
			fmt.Fprint(bytes, "{ ")
		}

		for i := 0; i < val.NumField(); i++ {
			var f = val.Field(i)
			fmt.Fprintf(bytes, "%s=", vtype.Field(i).Name)
			keyrenderrecurse(bytes, true, f.Interface())
		}

		if wrap {
			fmt.Fprint(bytes, "}")
		}

	case reflect.String:
		fmt.Fprintf(bytes, "%q", val.String())

	default:
		fmt.Fprintf(bytes, "%v", val.Interface())

		// end switch type
	}

	fmt.Fprint(bytes, " ")
}

func keyrender(data interface{}) []byte {

	var bytes = new(bytes.Buffer)
	keyrenderrecurse(bytes, false, data)
	return bytes.Bytes()

}

//----------------------------------------------------------------------
//----------------------------------------------------------------------
func (o *TextOutput) contextevent(cxttype string,
	context Context,
	event string,
	msg string,
	data *D,
	style *TerminalStyle) {

	var writsofar int

	if context.IsHighlighted() {
		style.bold = true
	}

	// Set the color
	var color int = WHITE
	if o.Color {
		color = style.color

		// Terminal commands produce no text, so don't include in writsofar
		o.printf("\x1b[%dm", style.intensity+color)
		if style.bold {
			o.printf("\x1b[1m")
		}
	}

	// Write the timestamp, root tag, label, and message

	writsofar += o.printf("%v %v %v ",
		o.timestamp(), context.GetRoot().Tag, context.GetLabel())

	writsofar += o.printf("%v ", msg)

	// Space out and then write the data fields

	for writsofar < o.IDOffset {
		writsofar += o.printf(" ")
	}

	if o.Color {
		if color != BLACK {
			o.printf("\x1b[0;%dm", LOW_INTENSITY+color)
		} else {
			o.printf("\x1b[0;%dm", HIGH_INTENSITY+color)
		}
	}

	writsofar += o.printf("%7v %v %-2v ", event, cxttype, context.GetUID())

	for writsofar < o.DataOffset {
		writsofar += o.printf(" ")
	}

	o.printf("%s", keyrender(data))

	if o.Color {
		o.printf("\x1b[0m")
	}
	o.printf("\n")

	// end contextevent
}

//----------------------------------------------------------------------
func (o *TextOutput) ComponentEvent(component *Component,
	event ComponentEventClass,
	msg string,
	data *D) {

	if InvalidComponentEventClass(event) {
		o.internalerror(NewError("ComponentEventClass out of range",
			component.GetUID(), event))
		return
	}

	var style = &ComponentEventTerminalStyles[event]

	o.contextevent("cmpt", component, ComponentEventClassText[event], msg, data, style)

	// end ComponentEvent
}

//----------------------------------------------------------------------
func (o *TextOutput) TaskEvent(task *Task,
	event TaskEventClass) {

	if InvalidTaskEventClass(event) {
		o.internalerror(NewError("TaskEventClass out of range for TaskProgress()",
			task.GetUID(), event))
		return
	}

	var style = &TaskEventTerminalStyles[event]

	var msg string = task.Activity

	switch event {
	case TASK_BEGIN:
		msg += " start"

	case TASK_END:
		msg += " success"
		if task.highlight {
			style = &TerminalStyle{GREEN, true, HIGH_INTENSITY}
		}

	case TASK_ERROR:
		msg += " failure"

	default:
		o.internalerror(NewError("TaskEventClass out of range for TaskEvent()",
			task.GetUID(), event))
		return

	}

	o.contextevent("task", task, TaskEventClassText[event], msg, task.Data, style)

	// end TaskEvent
}

//----------------------------------------------------------------------
func (o *TextOutput) TaskProgress(task *Task,
	event TaskEventClass,
	msg string,
	data *D) {

	if InvalidTaskEventClass(event) {
		o.internalerror(NewError("TaskEventClass out of range for TaskProgress()",
			task.GetUID(), event))
		return
	}

	var style = &TaskEventTerminalStyles[event]

	o.contextevent("task", task, TaskEventClassText[event], msg, data, style)

	// end TaskProgress
}
