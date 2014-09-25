package logberry

import (
	"os"
	"syscall"
	"io"
	"time"
	"fmt"
	"encoding/json"
)


//----------------------------------------------------------------------
//----------------------------------------------------------------------
type OutputDriver interface {
  log(component string,
      class StatementClass,
      msg string,
      data interface{}) error
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
	                                 err error) error {

	fmt.Fprintf(o.writer, "{ \"Class\": \"%s\", \"Component\": \"%s\", \"Msg\": \"%s\" }\n",
		CLASSTEXT[ERROR],
		component,
		err.Error())

	// end criticalerror
	return err
}


func (o *JSONOutput) log(component string,
                         class StatementClass,
                         msg string,
                         data interface{}) error {

	var entry = make(map[string]interface{})
	entry["Class"] = CLASSTEXT[class]
	entry["Component"] = component
	entry["Time"] = o.timestamp()
	entry["Msg"] = msg
	entry["Data"] = data

	var bytes []byte
	var err error
	bytes, err = json.Marshal(entry)
	if err != nil {
		wrapped := WrapError("Could not marshal log entry", err)
		return o.criticalerror(wrapped.Error(), err)
		return wrapped
	}

	o.writer.Write(bytes)
	o.writer.Write([]byte("\n"))

	return nil
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

	HIGH int = 90
	LOW int = 30
	)


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
	                                 err error) error {

	if o.Color {
		fmt.Fprintf(o.writer, "\x1b[%d;1m", HIGH+RED)
	}


	fmt.Fprintf(o.writer, "%v %v [ERROR] %v\n",
		o.timestamp(), component, err.Error())

	if o.Color {
		fmt.Fprintf(o.writer, "\x1b[0m")
	}

	o.writer.Write([]byte("\n"))

	// end criticalerror
	return err
}


func (o *TextOutput) log(component string,
                         class StatementClass,
                         msg string,
                         data interface{}) error {

	// Marshal the data first in case there's an error
	var bytes []byte
	if data == nil {
		bytes = []byte("{}")
	} else {
		var err error
		bytes, err = json.Marshal(data)
		if err != nil {
			wrapped := WrapError("Could not marshal log entry fields", err)
			return o.criticalerror(wrapped.Error(), err)
			return wrapped
		}
	}

	var writsofar int = 28 + len(component) + len(msg)

	// Set the color
	var color int = WHITE
	var bold bool = false
	if o.Color {
		switch class {
		case ERROR: fallthrough;
		case FATAL:
			color = RED; bold = true
		case METADATA: color = BLUE
		case INFO: color = WHITE
		}

		fmt.Fprintf(o.writer, "\x1b[%dm", HIGH+color)
		if bold {
			fmt.Fprintf(o.writer, "\x1b[1m")
		}
	}

	// Write the timestamp, component, and message
	o.writer.Write([]byte(o.timestamp()))
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

	o.writer.Write(bytes)

	if o.Color {
		o.writer.Write([]byte("\x1b[0m"))
	}
	o.writer.Write([]byte("\n"))


	return nil
}
