package logberry


//----------------------------------------------------------------------
//----------------------------------------------------------------------
type ThreadSafeOutput struct {
	root *Root
	driver OutputDriver
	channel chan logentry
}


type logentry interface {
	Log(driver OutputDriver)
}

type componententry struct {
	component *Component
	class ContextEventClass
	msg string
	data *D
}

func (x *componententry) Log(driver OutputDriver) {
	driver.ComponentEvent(x.component, x.class, x.msg, x.data)
}

type taskentry struct {
	task *Task
	event ContextEventClass
}

func (x *taskentry) Log(driver OutputDriver) {
	driver.TaskEvent(x.task, x.event)
}


//----------------------------------------------------------------------
//----------------------------------------------------------------------
func NewThreadSafeOutput(driver OutputDriver, buffer int) *ThreadSafeOutput {

	ts := &ThreadSafeOutput {
		driver: driver,
		channel: make(chan logentry, buffer),
	}

	go ts.process()

	return ts
}

func (x *ThreadSafeOutput) process() {

	for {
		e := <- x.channel
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

	x.channel <- &componententry {
		component: component,
		class: class,
		msg: msg,
		data: data,
	}

}

func (x *ThreadSafeOutput) TaskEvent(task *Task,
	event ContextEventClass) {

	x.channel <- &taskentry {
		task: task,
		event: event,
	}

}
