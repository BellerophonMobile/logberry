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
	intensity int
}

var ReportTerminalStyles = [...]terminalstyle {
	{ RED,    true,  HIGH_INTENSITY },               // error
	{ RED,    true,  HIGH_INTENSITY },               // fatal
	{ YELLOW, true,  HIGH_INTENSITY },               // warning
	{ WHITE,  false, HIGH_INTENSITY },               // info
	{ BLUE,   false, LOW_INTENSITY  },               // configuration
	{ GREEN,  true,  HIGH_INTENSITY },               // start
	{ BLACK,  false, HIGH_INTENSITY },               // finish
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
	if o.Color {
		color = ReportTerminalStyles[class].color

		// Terminal commands produce no text, so don't include in writsofar
		o.printf("\x1b[%dm", ReportTerminalStyles[class].intensity+color)
		if ReportTerminalStyles[class].bold {
			o.printf("\x1b[1m")
		}
	}

	// Write the timestamp, component, and message

	writsofar += o.printf("%v %v ",
		o.timestamp(), context.Label)

	switch class {
	case ERROR:
		writsofar += o.printf("[ERROR] ")
	case FATAL:
		writsofar += o.printf("[FATAL ERROR] ")
	}

	if context.class == TASK) {

}

	writsofar += o.printf("%v ", msg)

	// Space out and then write the data fields

	for writsofar < 72 {
		writsofar += o.printf(" ")
	}

	if o.Color {
		if color != BLACK {
			o.printf("\x1b[0;%dm", LOW_INTENSITY+color)
		} else {
			o.printf("\x1b[0;%dm", HIGH_INTENSITY+color)
		}
	}

	o.printf("%v %v %s", EventClassText[class], context.ID, bytes)

	if o.Color {
		o.printf("\x1b[0m")
	}
	o.printf("\n")

}
