package logberry

//----------------------------------------------------------------------
//----------------------------------------------------------------------
type ThreadSafeOutput struct {
	root    *Root
	driver  OutputDriver
	channel chan logevent
}

type logevent interface {
	Log(driver OutputDriver)
}

type componentevent struct {
	component *Component
	class     ContextEventClass
	msg       string
	data      *D
}

func (x *componentevent) Log(driver OutputDriver) {
	driver.ComponentEvent(x.component, x.class, x.msg, x.data)
}

type taskevent struct {
	task  *Task
	event ContextEventClass
}

func (x *taskevent) Log(driver OutputDriver) {
	driver.TaskEvent(x.task, x.event)
}

type taskprogress struct {
	task  *Task
	event ContextEventClass
	msg   string
	data  *D
}

func (x *taskprogress) Log(driver OutputDriver) {
	driver.TaskProgress(x.task, x.event, x.msg, x.data)
}

//----------------------------------------------------------------------
//----------------------------------------------------------------------
func NewThreadSafeOutput(driver OutputDriver, buffer int) *ThreadSafeOutput {

	ts := &ThreadSafeOutput{
		driver:  driver,
		channel: make(chan logevent, buffer),
	}

	go ts.process()

	return ts
}

func (x *ThreadSafeOutput) process() {

	for {
		e := <-x.channel
		e.Log(x.driver)
	}

}

//----------------------------------------------------------------------
func (x *ThreadSafeOutput) Attach(root *Root) {
	x.root = root
	x.driver.Attach(root)
}

func (x *ThreadSafeOutput) Detach() {
	x.driver.Detach()
	x.root = nil
}

func (x *ThreadSafeOutput) ComponentEvent(component *Component,
	class ContextEventClass,
	msg string,
	data *D) {

	x.channel <- &componentevent{
		component: component,
		class:     class,
		msg:       msg,
		data:      data,
	}

}

func (x *ThreadSafeOutput) TaskEvent(task *Task,
	event ContextEventClass) {

	x.channel <- &taskevent{
		task:  task,
		event: event,
	}

}

func (x *ThreadSafeOutput) TaskProgress(task *Task,
	event ContextEventClass,
	msg string,
	data *D) {

	x.channel <- &taskprogress{
		task:  task,
		event: event,
		msg:   msg,
		data:  data,
	}

}
