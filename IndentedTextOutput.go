// +build !appengine

package logberry

import (
	"fmt"
	"io"
	"os"
	"sort"
	"strings"
	"syscall"
	"time"

	"github.com/BellerophonMobile/logberry/terminal"
)

// IndentedTextOutput is an OutputDriver that writes out log events in
// a structured but more or less human readable form. It has the
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
type IndentedTextOutput struct {
	root   *Root
	writer io.Writer

	Program string

	Color bool

	IDOffset   int
	DataOffset int
	Indent     string
}

// NewIndentedStdOutput creates a new IndentedTextOutput attached to stdout.
func NewIndentedStdOutput(program string) *IndentedTextOutput {
	t := NewIndentedTextOutput(os.Stdout, program)
	t.Color = terminal.IsTerminal(syscall.Stdout)
	return t
}

// NewIndentedErrOutput creates a new IndentedTextOutput attached to stderr.
func NewIndentedErrOutput(program string) *IndentedTextOutput {
	t := NewIndentedTextOutput(os.Stderr, program)
	t.Color = terminal.IsTerminal(syscall.Stderr)
	return t
}

// NewIndentedTextOutput creates a new IndentedTextOutput attached to the given writer.
func NewIndentedTextOutput(w io.Writer, program string) *IndentedTextOutput {
	return &IndentedTextOutput{
		writer:     w,
		Program:    program,
		IDOffset:   80,
		DataOffset: 100,
		Indent:     "  ",
	}
}

// Attach notifies the OutputDriver of its Root.  It should only be
// called by a Root.
func (o *IndentedTextOutput) Attach(root *Root) {
	o.root = root
}

// Detach notifies the OutputDriver that it has been removed from its
// Root.  It should only be called by a root.
func (o *IndentedTextOutput) Detach() {
	o.root = nil
}

// Event outputs a generated log entry, as called by a Root or a
// chaining OutputDriver.
func (o *IndentedTextOutput) Event(event *Event) {

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

	// if len(event.Data) > 0 {
	// 	event.Data.WriteTo(o.writer)
	// }
	if len(event.Data) > 0 {
		fmt.Fprintf(o.writer, "\n")
		o.writeEventDataMap(event.Data, 1, true)
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

func (o *IndentedTextOutput) writeEventDataMap(datamap EventDataMap, indent int, printIndent bool) {
	if len(datamap) == 0 {
		return
	}

	keys := make([]string, len(datamap))
	i := 0
	for k := range datamap {
		keys[i] = k
		i++
	}
	sort.Strings(keys)

	for i, k := range keys {
		v := datamap[k]

		if printIndent {
			for n := 0; n < indent; n++ {
				fmt.Fprintf(o.writer, o.Indent)
			}
		}

		if strings.ContainsAny(k, "\"= {}[]") {
			fmt.Fprintf(o.writer, "%q", k)
		} else {
			fmt.Fprintf(o.writer, "%v", k)
		}

		switch v := v.(type) {
		case EventDataMap:
			if len(v) > 1 {
				fmt.Fprintf(o.writer, ":\n")
				o.writeEventDataMap(v, indent+1, true)
			} else if len(v) == 1 {
				fmt.Fprintf(o.writer, ".")
				o.writeEventDataMap(v, indent, false)
			}

		case EventDataSlice:
			fmt.Fprintf(o.writer, " = [")
			for i, item := range v {
				switch item := item.(type) {
				case EventDataMap:
					fmt.Fprintf(o.writer, "{\n")

					o.writeEventDataMap(item, indent+1, true)

					fmt.Fprintf(o.writer, "\n")
					for n := 0; n < indent; n++ {
						fmt.Fprintf(o.writer, o.Indent)
					}

					if i == len(v)-1 {
						fmt.Fprintf(o.writer, "}")
					} else {
						fmt.Fprintf(o.writer, "}, ")
					}

				default:
					item.WriteTo(o.writer)
					if i < len(v)-1 {
						fmt.Fprintf(o.writer, ", ")
					}
				}
			}

			fmt.Fprintf(o.writer, "]")

		default:
			fmt.Fprintf(o.writer, " = ")
			v.WriteTo(o.writer)
		}

		if i < len(keys)-1 {
			fmt.Fprintf(o.writer, "\n")
		}
	}
}
