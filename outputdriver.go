package logberry

import (
	"os"
	"syscall"
	"io"
	"time"
	"fmt"
	"encoding/json"
	"log"
)

func init() {
	if len(TERMINAL_STYLES) != int(UNKNOWN) + 1 {
		log.Fatal("Fatal internal error: len(TERMINAL_STYLES) != |StatementClass|")
	}

	// end init
}

//----------------------------------------------------------------------
//----------------------------------------------------------------------
type OutputDriver interface {
  Log(component string,
      class StatementClass,
      msg string,
      data interface{})
}


//----------------------------------------------------------------------
//----------------------------------------------------------------------
type JSONOutput struct {
	writer io.Writer

	Start time.Time
	DifferentialTime bool
}

//----------------------------------------------------------------------
func NewJSONOutput(w io.Writer) *JSONOutput {
	return &JSONOutput{
		writer: w,
		Start: time.Now(),
		DifferentialTime: false,
	}
}


//------------------------------------------------------
func (o *JSONOutput) timestamp() string {
	if (o.DifferentialTime) {
		return time.Since(o.Start).String()
	}

	return time.Now().Format(time.RFC3339)
	// end timestamp
}


func (o *JSONOutput) criticalerror(component string,
	                                 err error) {
 
	fmt.Fprintf(o.writer, "{ \"Class\": \"%s\", \"Component\": \"%s\", \"Msg\": \"%s\" }\n",
		STATEMENT_CLASS_TEXT[ERROR],
		component,
		err.Error())

	LoggingError(err)
}


func (o *JSONOutput) Log(component string,
                         class StatementClass,
                         msg string,
                         data interface{}) {

	var entry = make(map[string]interface{})
	entry["Class"] = STATEMENT_CLASS_TEXT[class]
	entry["Component"] = component
	entry["Time"] = o.timestamp()
	entry["Msg"] = msg
	entry["Data"] = data

	var bytes []byte
	var err error
	bytes, err = json.Marshal(entry)
	if err != nil {
		o.criticalerror(component, WrapError(err, "Could not marshal log entry"))
		return
	}

	o.writer.Write(bytes)
	o.writer.Write([]byte("\n"))

	// end JSONOutput::Log
}


//----------------------------------------------------------------------
//----------------------------------------------------------------------
type TextOutput struct {
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
	HIGH int = 90
	LOW int = 30
)

type terminalstyle struct {
	color int
	bold bool
}

var TERMINAL_STYLES = [...]terminalstyle {
	{ RED, true},               // error
	{ RED, true},               // fatal
	{ YELLOW, true},            // warning
	{ WHITE, false},            // info
	{ BLUE, false},             // configuration
	{ GREEN, false},            // instantiate
	{ BLACK, false},            // finalize
	{ WHITE, false},            // task_start
	{ WHITE, false},            // task_finish
	{ WHITE, false},            // resource
	{ WHITE, false},            // service
	{ RED, true},               // unknown
}


//----------------------------------------------------------------------
func NewStdOutput() *TextOutput {
	t := NewTextOutput(os.Stdout)
	t.Color = IsTerminal(syscall.Stdout)
	return t
}

func NewErrOutput() *TextOutput {
	t := NewTextOutput(os.Stderr)
	t.Color = IsTerminal(syscall.Stderr)
	return t
}

func NewTextOutput(w io.Writer) *TextOutput {
	return &TextOutput{
		writer: w,
		Start: time.Now(),
		DifferentialTime: false,
		Color: false,
	}
}


//------------------------------------------------------
func (o *TextOutput) timestamp() string {

	if (o.DifferentialTime) {
		return time.Since(o.Start).String()
	}

	return time.Now().Format(time.RFC3339)

	// end timestamp
}


func (o *TextOutput) criticalerror(component string,
	                                 err error) {

	if o.Color {
		fmt.Fprintf(o.writer, "\x1b[%d;1m", HIGH+RED)
	}

	fmt.Fprintf(o.writer, "%v %v [ERROR] %v\n",
		o.timestamp(), component, err.Error())

	if o.Color {
		fmt.Fprintf(o.writer, "\x1b[0m")
	}

	o.writer.Write([]byte("\n"))

	LoggingError(err)
}


func (o *TextOutput) Log(component string,
                         class StatementClass,
                         msg string,
                         data interface{}) {

	if class < 0 || class > UNKNOWN {
		o.criticalerror(component, NewError("Class", class, "out of range"))
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
			o.criticalerror(component,
				WrapError(err, "Could not marshal log entry fields"))
			return
		}
	}

	var writsofar int = 28 + len(program) + len(component) + len(msg)

	// Set the color
	var color int = WHITE
	var bold bool = false
	if o.Color {
		color = TERMINAL_STYLES[class].color
		bold = TERMINAL_STYLES[class].bold

		fmt.Fprintf(o.writer, "\x1b[%dm", HIGH+color)
		if bold {
			fmt.Fprintf(o.writer, "\x1b[1m")
		}
	}

	// Write the timestamp, component, and message
	o.writer.Write([]byte(o.timestamp()))
	o.writer.Write([]byte(" "))
	o.writer.Write([]byte(program))
	o.writer.Write([]byte(" "))
	o.writer.Write([]byte(component))
	o.writer.Write([]byte(" "))

	switch class {
	case ERROR:
		o.writer.Write([]byte("[ERROR] "))
		writsofar += 8
	case FATAL:
		o.writer.Write([]byte("[FATAL ERROR] "))
		writsofar += 14
	}

	o.writer.Write([]byte(msg))
	o.writer.Write([]byte(" "))

	// Space out and then write the data fields

	if class == ERROR {
	}

	for ; writsofar < 72; writsofar++ {
		o.writer.Write([]byte(" "))
	}

	if o.Color {
		fmt.Fprintf(o.writer, "\x1b[0;%dm", LOW+color)
	}

	o.writer.Write([]byte(STATEMENT_CLASS_TEXT[class]))
	o.writer.Write([]byte(" "))
	o.writer.Write(bytes)

	if o.Color {
		o.writer.Write([]byte("\x1b[0m"))
	}
	o.writer.Write([]byte("\n"))

}
