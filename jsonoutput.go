package logberry

import (
	"fmt"
	"io"
	"time"
	"encoding/json"
)


type JSONOutput struct {
	root *Root
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


//----------------------------------------------------------------------
func (x *JSONOutput) Attach(root *Root) {
	x.root = root
}

func (x *JSONOutput) Detach() {
	x.root = nil
}


//----------------------------------------------------------------------
func (x *JSONOutput) timestamp() string {

	if (x.DifferentialTime) {
		return time.Since(x.Start).String()
	}

	return time.Now().Format(time.RFC3339)

	// end timestamp
}


func (x *JSONOutput) internalerror(err error) {
 
	fmt.Fprintf(x.writer, "{\"EntryType\":\"error\",\"Error\": %q}\n",
		err.Error())

	x.root.InternalError(WrapError(err, "Could not output log entry"))

}


//----------------------------------------------------------------------
//----------------------------------------------------------------------
func (x *JSONOutput) contextevent(entrytype string,
	context Context,
	event ContextEventClass,
	msg string,
	data *D) {

	var entry = make(map[string]interface{})
	entry["EntryType"] = entrytype
	entry["Event"] = ContextEventClassText[event]
	entry["UID"] = context.GetUID()
	entry["Tag"] = x.root.Tag
	entry["Label"] = context.GetLabel()
	entry["Msg"] = msg
	entry["Data"] = data
	entry["Time"] = x.timestamp()

	var bytes []byte
	var err error
	bytes, err = json.Marshal(entry)
	if err != nil {
		x.internalerror(WrapError(err, "Could not marshal log entry"))
		return
	}

	x.writer.Write(bytes)
	x.writer.Write([]byte("\n"))

}

func (x *JSONOutput) ComponentEvent(component *Component,
  event ContextEventClass,
  msg string,
  data *D) {

	if event < 0 || event >= contexteventclasssentinel {
		x.internalerror(NewError("ContextEventClass out of range for component event",
			component.GetUID(), event))
		return
	}

	x.contextevent("component", component, event, msg, data)

	// end ComponentEvent
}

func (x *JSONOutput) TaskEvent(task *Task,
  event ContextEventClass) {

	var msg string = task.Activity

	switch event {
	case START:
		msg += " start"

	case FINISH:
		if (task.Timed) {
			msg += " finished"
		}

	case SUCCESS:
		msg += " success"

	case ERROR:
		msg += " failed"

	default:
		x.internalerror(NewError("ContextEventClass out of range for task event",
			task.GetUID(), event))
		return

	}

	x.contextevent("task", task, event, msg, task.Data)

	// end TaskEvent
}

func (x *JSONOutput) TaskProgress(task *Task,
  event ContextEventClass,
  msg string,
  data *D) {

	if event < 0 || event >= contexteventclasssentinel {
		x.internalerror(NewError("ContextEventClass out of range for task progress",
			task.GetUID(), event))
		return
	}

	x.contextevent("task", task, event, msg, data)

	// end TaskProgress
}
