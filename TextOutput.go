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
// or very wide terminals, generally putting each event on 1 or 2
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
	high_intensity int = 90
	low_intensity  int = 30
)

type terminalstyle struct {
	color     int
	bold      bool
	intensity int
}

var defaultstyle = terminalstyle{white,  false, high_intensity}  // default

var eventstyles = map[string]terminalstyle{
  BEGIN:         {black,  false, high_intensity},  // begin
  END:           {black,  false, high_intensity},  // end
  CONFIGURATION: {blue,   false, low_intensity},   // configuration
  READY:         {green,  true,  high_intensity},  // ready
  STOPPED:       {white,  false, high_intensity},  // stopped
  INFO:          {white,  false, high_intensity},  // info
  WARNING:       {yellow, false, high_intensity},  // warning
  ERROR:         {red,    true,  high_intensity},  // error
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
		IDOffset:   84,
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

	style,ok := eventstyles[event.Event]
	if !ok {
		style = defaultstyle
	}
//	if event.highlight {
//		style.bold = true
//	}

	var writsofar int  // Track characters writen so far, to space;
										 // terminal commands produce no text, so don't
										 // include in writsofar

	
	// Set the color
	var color int = style.color

	// Write the timestamp, program tag, and component

	if o.Color {
		var c = high_intensity
		if color != black {
			c = low_intensity
		}
		
		_, e := fmt.Fprintf(o.writer, "\x1b[%dm", c+color)
		if e != nil {
			o.root.internalerror(WrapError("Could write entry", e))
			return
		}
	}

	n,e := fmt.Fprintf(o.writer, "%v %v %v ",
		event.Timestamp.Format(time.RFC3339), o.Program, event.Component)
	if e != nil {
		o.root.internalerror(WrapError("Could write entry", e))
		return
	}
	writsofar += n

	
	// Write the message
	if o.Color {
		_,e := fmt.Fprintf(o.writer, "\x1b[%dm", style.intensity+color)
		if e != nil {
			o.root.internalerror(WrapError("Could write entry", e))
			return
		}
		
		if style.bold {
			_,e := fmt.Fprintf(o.writer, "\x1b[1m")
			if e != nil {
				o.root.internalerror(WrapError("Could write entry", e))
				return
			}

		}
	}

	n,e = fmt.Fprintf(o.writer, "%v ", event.Message)
	if e != nil {
		o.root.internalerror(WrapError("Could write entry", e))
		return
	}
	writsofar += n


	// Space out and then write the data fields
	for writsofar < o.IDOffset {
		n,_ = fmt.Fprintf(o.writer, " ")
		writsofar += n
	}

	if o.Color {
		if color != black {
			_, e := fmt.Fprintf(o.writer, "\x1b[0;%dm", low_intensity+color)
			if e != nil {
				o.root.internalerror(WrapError("Could write entry", e))
				return
			}

		} else {
			_,e := fmt.Fprintf(o.writer, "\x1b[0;%dm", high_intensity+color)
			if e != nil {
				o.root.internalerror(WrapError("Could write entry", e))
				return
			}

		}
	}

	n,e = fmt.Fprintf(o.writer, "%16v %2v:%-2v",
		event.Event, event.TaskID, event.ParentID)
	if e != nil {
		o.root.internalerror(WrapError("Could write entry", e))
		return
	}
	writsofar += n
	
	for writsofar < o.DataOffset {
		n,e = fmt.Fprintf(o.writer, " ")
		if e != nil {
			o.root.internalerror(WrapError("Could write entry", e))
			return
		}
		writsofar += n
	}

	event.Data.WriteTo(o.writer)

	if o.Color {
		_,e := fmt.Fprintf(o.writer, "\x1b[0m")
		if e != nil {
			o.root.internalerror(WrapError("Could write entry", e))
			return
		}
	}
	_,e = fmt.Fprintf(o.writer, "\n")
	if e != nil {
		o.root.internalerror(WrapError("Could write entry", e))
		return
	}

	// end Event
}
