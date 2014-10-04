package logberry

import (
	"os"
	"syscall"
	"io"
	"time"
	"fmt"
	"encoding/json"

	"github.com/BellerophonMobile/logberry/terminal"
)


//----------------------------------------------------------------------
//----------------------------------------------------------------------
type TextOutput struct {
	root *Root
	writer io.Writer

	Start time.Time
	DifferentialTime bool

	Color bool
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
	LOW_INTENSITY int = 30
)

type terminalstyle struct {
	color int
	bold bool
}

var TERMINAL_STYLES = [...]terminalstyle {
	{ RED, true },               // error
	{ RED, true },               // fatal
	{ YELLOW, true },            // warning
	{ WHITE, false },            // info
	{ BLUE, false },             // configuration
	{ GREEN, true },             // ready
	{ GREEN, false },            // instantiate
	{ BLACK, false },            // finalize
	{ WHITE, false },            // task
	{ WHITE, false },            // task_start
	{ WHITE, false },            // task_finish
	{ WHITE, false },            // resource
	{ WHITE, false },            // service
	{ WHITE, false },            // query
	{ WHITE, false },            // assert
	{ WHITE, false },            // calculate
	{ WHITE, false },            // read
	{ WHITE, false },            // write
	{ WHITE, false },            // connect
	{ WHITE, false },            // disconnect
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
		writer: w,
		Start: time.Now(),
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

	if (o.DifferentialTime) {
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

	o.root.InternalError(WrapError(err, "Could not write text output"))

}

// Convenience function to handle the potential write error.
func (o *TextOutput) printf(msg string, a ...interface{}) int {
	n,e := fmt.Fprintf(o.writer, msg, a...)
	if e != nil {
		o.internalerror(e)
	}
	return n
}


//----------------------------------------------------------------------
//----------------------------------------------------------------------
func (o *TextOutput) Report(context *Context,
                            class EventClass,
                            msg string,
                            data *D) {

	if class < 0 || class >= eventclass_sentinel {
		o.internalerror(NewError("EventClass out of range", class, context))
		return
	}

	// Marshal the data first in case there's an error
	var bytes []byte
	if data == nil {
		bytes = []byte("{}")
	} else {
		var err error
		bytes, err = json.Marshal(data)
		if err != nil {
			o.internalerror(WrapError(err, "Could not marshal log entry fields",
				context))
			return
		}
	}

	var writsofar int

	// Set the color
	var color int = WHITE
	var bold bool = false
	if o.Color {
		color = TERMINAL_STYLES[class].color
		bold = TERMINAL_STYLES[class].bold

		writsofar += o.printf("\x1b[%dm", HIGH_INTENSITY+color)
		if bold {
			writsofar += o.printf("\x1b[1m")
		}
	}

	// Write the timestamp, component, and message
	writsofar += o.printf("%v %v %v ",
		o.timestamp(), context.Root.Tag, context.Label)

	switch class {
	case ERROR:
		writsofar += o.printf("[ERROR] ")
	case FATAL:
		writsofar += o.printf("[FATAL ERROR] ")
	}

	writsofar += o.printf("%v ", msg)

	// Space out and then write the data fields

	for writsofar < 72 {
		writsofar += o.printf(" ")
	}

	if o.Color {
		o.printf("\x1b[0;%dm", LOW_INTENSITY+color)
	}

	o.printf("%v %s", EventClassText[class], bytes)

	if o.Color {
		o.printf("\x1b[0m")
	}
	o.printf("\n")

}
