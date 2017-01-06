// +build !appengine

package logberry

import (
	"fmt"
	"io"
	"os"
	"syscall"
	"time"

	"github.com/BellerophonMobile/logberry/terminal"
)

// TextOutput is an OutputDriver that writes out log events in a
// structured but more or less human readable form.  It has the
// following public properties:
//
//   Program                    String label of the executing program.
//
//   Color                      Set to true/false to enable/disable
//                              outputting terminal color codes as
//                              part of formatting log entries.
//                              Defaults to false except for when
//                              constructed via NewStdOutput and
//                              NewErrOutput as below, in which case
//                              it defaults to true iff the underlying
//                              streams are terminals.
//
//   IDOffset                   The column at which to start printing
//                              identifying information.
//
//   DataOffset                 The column at which to start printing
//                              event data.
//
// The default offsets are designed to wrap well on either 80 column
// or very wide terminals, generally putting each event on one or two
// lines respectively.
type TextOutput struct {
	root   *Root
	writer io.Writer

	Program string

	Color bool

	IDOffset   int
	DataOffset int
}

const (
	black int = iota
	red
	green
	yellow
	blue
	magenta
	cyan
	white
)

const (
	highintensity int = 90
	lowintensity  int = 30
)

type terminalstyle struct {
	color     int
	bold      bool
	intensity int
}

var defaultstyle = terminalstyle{cyan, false, highintensity} // default

var eventstyles = map[string]terminalstyle{
	BEGIN:         {black, false, highintensity},   // begin
	END:           {black, false, highintensity},   // end
	CONFIGURATION: {blue, false, lowintensity},     // configuration
	READY:         {green, true, highintensity},    // ready
	STOPPED:       {magenta, false, highintensity}, // stopped
	INFO:          {white, false, highintensity},   // info
	SUCCESS:       {white, false, highintensity},   // end
	WARNING:       {yellow, false, highintensity},  // warning
	ERROR:         {red, true, highintensity},      // error
}

// NewStdOutput creates a new TextOutput attached to stdout.
func NewStdOutput(program string) *TextOutput {
	t := NewTextOutput(os.Stdout, program)
	t.Color = terminal.IsTerminal(syscall.Stdout)
	return t
}

// NewErrOutput creates a new TextOutput attached to stderr.
func NewErrOutput(program string) *TextOutput {
	t := NewTextOutput(os.Stderr, program)
	t.Color = terminal.IsTerminal(syscall.Stderr)
	return t
}

// NewTextOutput creates a new TextOutput attached to the given writer.
func NewTextOutput(w io.Writer, program string) *TextOutput {
	return &TextOutput{
		writer:     w,
		Program:    program,
		IDOffset:   80,
		DataOffset: 100,
	}
}

// Attach notifies the OutputDriver of its Root.  It should only be
// called by a Root.
func (o *TextOutput) Attach(root *Root) {
	o.root = root
}

// Detach notifies the OutputDriver that it has been removed from its
// Root.  It should only be called by a root.
func (o *TextOutput) Detach() {
	o.root = nil
}

// Event outputs a generated log entry, as called by a Root or a
// chaining OutputDriver.
func (o *TextOutput) Event(event *Event) {

	style, ok := eventstyles[event.Event]
	if !ok {
		style = defaultstyle
	}
	//	if event.highlight {
	//		style.bold = true
	//	}

	var writsofar int // Track characters written so far, to space;
	// terminal commands produce no text, so don't
	// include in writsofar

	// Set the color
	var color = style.color

	// Write the timestamp, program tag, and component

	if o.Color {
		var c = highintensity
		if color != black {
			c = lowintensity
		}

		_, e := fmt.Fprintf(o.writer, "\x1b[%dm", c+color)
		if e != nil {
			o.root.InternalError(WrapError("Could write entry", e))
			return
		}
	}

	n, e := fmt.Fprintf(o.writer, "%v %v %-12v ",
		event.Timestamp.Format(time.RFC3339), o.Program, event.Component)
	if e != nil {
		o.root.InternalError(WrapError("Could write entry", e))
		return
	}
	writsofar += n

	// Write the message
	if o.Color {
		_, e := fmt.Fprintf(o.writer, "\x1b[%dm", style.intensity+color)
		if e != nil {
			o.root.InternalError(WrapError("Could write entry", e))
			return
		}

		if style.bold {
			_, e := fmt.Fprintf(o.writer, "\x1b[1m")
			if e != nil {
				o.root.InternalError(WrapError("Could write entry", e))
				return
			}

		}
	}

	n, e = fmt.Fprintf(o.writer, "%v ", event.Message)
	if e != nil {
		o.root.InternalError(WrapError("Could write entry", e))
		return
	}
	writsofar += n

	// Space out and then write the data fields
	for writsofar < o.IDOffset {
		n, _ = fmt.Fprintf(o.writer, " ")
		writsofar += n
	}

	if o.Color {
		if color != black {
			_, e := fmt.Fprintf(o.writer, "\x1b[0;%dm", lowintensity+color)
			if e != nil {
				o.root.InternalError(WrapError("Could write entry", e))
				return
			}

		} else {
			_, e := fmt.Fprintf(o.writer, "\x1b[0;%dm", highintensity+color)
			if e != nil {
				o.root.InternalError(WrapError("Could write entry", e))
				return
			}

		}
	}

	n, e = fmt.Fprintf(o.writer, "%16v %2v:%-2v",
		event.Event, event.TaskID, event.ParentID)
	if e != nil {
		o.root.InternalError(WrapError("Could write entry", e))
		return
	}
	writsofar += n

	for writsofar < o.DataOffset {
		n, e = fmt.Fprintf(o.writer, " ")
		if e != nil {
			o.root.InternalError(WrapError("Could write entry", e))
			return
		}
		writsofar += n
	}

	if len(event.Data) > 0 {
		event.Data.WriteTo(o.writer)
	}

	if o.Color {
		_, e := fmt.Fprintf(o.writer, "\x1b[0m")
		if e != nil {
			o.root.InternalError(WrapError("Could write entry", e))
			return
		}
	}
	_, e = fmt.Fprintf(o.writer, "\n")
	if e != nil {
		o.root.InternalError(WrapError("Could write entry", e))
		return
	}

	// end Event
}
