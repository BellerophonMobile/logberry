package logberry

import (
	"log"

	"fmt"
	"io"
	"os"
	"syscall"
	"time"
	//	"encoding/json"
	"bytes"

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

var ContextEventTerminalStyles = [...]TerminalStyle{
	{RED, true, HIGH_INTENSITY},    // error
	{RED, true, HIGH_INTENSITY},    // fatal
	{YELLOW, true, HIGH_INTENSITY}, // warning
	{WHITE, false, HIGH_INTENSITY}, // info
	{BLUE, false, LOW_INTENSITY},   // configuration
	{GREEN, true, HIGH_INTENSITY},  // start
	{BLACK, false, HIGH_INTENSITY}, // finish
	{WHITE, false, HIGH_INTENSITY}, // success
}

func init() {

	//-- Check that labels are defined for the enumerations
	if len(ContextEventTerminalStyles) != int(contexteventclasssentinel) {
		log.Fatal("Fatal internal error: " +
			"len(ContextEventTerminalStyles) != |ContextEventClass|")
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
	event ContextEventClass,
	msg string,
	data *D,
	style *TerminalStyle) {

	// Marshal the data first in case there's an error
	var bytes []byte = keyrender(data)
	/*
		if data == nil {
			bytes = []byte("{}")
		} else {
			var err error
			bytes, err = json.Marshal(data)
			if err != nil {
				o.internalerror(WrapError(err, "Could not marshal log entry fields", context.GetUID()))
				return
			}
		}
	*/

	var writsofar int

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

	switch event {
	case ERROR:
		writsofar += o.printf("[ERROR] ")
	case FATAL:
		writsofar += o.printf("[FATAL ERROR] ")
	}

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

	writsofar += o.printf("%7v %v %-2v ", ContextEventClassText[event], cxttype, context.GetUID())

	for writsofar < o.DataOffset {
		writsofar += o.printf(" ")
	}

	o.printf("%s", bytes)

	if o.Color {
		o.printf("\x1b[0m")
	}
	o.printf("\n")

	// end contextevent
}

//----------------------------------------------------------------------
func (o *TextOutput) ComponentEvent(component *Component,
	event ContextEventClass,
	msg string,
	data *D) {

	if event < 0 || event >= contexteventclasssentinel {
		o.internalerror(NewError("ContextEventClass out of range for component event",
			component.GetUID(), event))
		return
	}

	o.contextevent("cmpt", component, event, msg, data,
		&ContextEventTerminalStyles[event])

	// end ComponentEvent
}

//----------------------------------------------------------------------
func (o *TextOutput) TaskEvent(task *Task,
	event ContextEventClass) {

	var msg string = task.Activity
	var style *TerminalStyle

	switch event {
	case START:
		msg += " start"
		// if task.Timed { msg = "@Start " + msg }
		style = &TerminalStyle{WHITE, false, HIGH_INTENSITY}

	case FINISH:
		msg += " completed"
		style = &TerminalStyle{WHITE, false, HIGH_INTENSITY}

	case SUCCESS:
		msg += " success"
		// if task.Timed { msg = "@Success " + msg } else { }
		style = &TerminalStyle{GREEN, false, HIGH_INTENSITY}

	case ERROR:
		msg += " failed"
		// if task.Timed { msg = "@Failed " + msg } else { }
		style = &TerminalStyle{RED, true, HIGH_INTENSITY}

	default:
		o.internalerror(NewError("ContextEventClass out of range for task progress",
			task.GetUID(), event))
		return

	}

	o.contextevent("task", task, event, msg, task.Data, style)

	// end TaskEvent
}

//----------------------------------------------------------------------
func (o *TextOutput) TaskProgress(task *Task,
	event ContextEventClass,
	msg string,
	data *D) {

	if event < 0 || event >= contexteventclasssentinel {
		o.internalerror(NewError("ContextEventClass out of range for task progress",
			task.GetUID(), event))
		return
	}

	o.contextevent("task", task, event, msg, data,
		&ContextEventTerminalStyles[event])

	// end TaskProgress
}
