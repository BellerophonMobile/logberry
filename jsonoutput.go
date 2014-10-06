package logberry

import (
	"encoding/json"
	"fmt"
	"io"
	"time"
)

type JSONOutput struct {
	root   *Root
	writer io.Writer

	Start            time.Time
	DifferentialTime bool
}

//----------------------------------------------------------------------
func NewJSONOutput(w io.Writer) *JSONOutput {
	return &JSONOutput{
		writer:           w,
		Start:            time.Now(),
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

	if x.DifferentialTime {
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
	event string,
	msg string,
	data *D) {

	var entry = make(map[string]interface{})
	entry["EntryType"] = entrytype
	entry["Event"] = event
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
	event ComponentEventClass,
	msg string,
	data *D) {

	if InvalidComponentEventClass(event) {
		x.internalerror(NewError("ComponentEventClass out of range",
			component.GetUID(), event))
		return
	}

	x.contextevent("component", component, ComponentEventClassText[event], msg, data)

	// end ComponentEvent
}

func (x *JSONOutput) TaskEvent(task *Task,
	event TaskEventClass) {

	var msg string = task.Activity

	switch event {
	case TASK_BEGIN:
		msg += " start"

	case TASK_END:
		if task.Timed {
			msg += " success"
		}

	case TASK_ERROR:
		msg += " failure"

	default:
		x.internalerror(NewError("TaskEventClass out of range for TaskEvent()",
			task.GetUID(), event))
		return

	}

	x.contextevent("task", task, TaskEventClassText[event], msg, task.Data)

	// end TaskEvent
}

func (x *JSONOutput) TaskProgress(task *Task,
	event TaskEventClass,
	msg string,
	data *D) {

	if InvalidTaskEventClass(event) {
		x.internalerror(NewError("TaskEventClass out of range for TaskProgress()",
			task.GetUID(), event))
		return
	}

	x.contextevent("task", task, TaskEventClassText[event], msg, data)

	// end TaskProgress
}
