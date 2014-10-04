package logberry

import (
	"os"
	"path"
	"path/filepath"
	"strings"
	"os/user"
	"sync/atomic"
)


type Context struct {
	ID uint64
	Root *Root
	Parent *Context
	Class ContextClass
	Label string
}

var numcontexts uint64


//----------------------------------------------------------------------
//----------------------------------------------------------------------
func NewContext(parent *Context, class ContextClass, label string) *Context {

	c := &Context {
		ID: atomic.AddUint64(&numcontexts, 1),
		Parent: parent,
		Class: class,
		Label: label,
	}

	if parent != nil {
		c.Root = parent.Root
	}

	return c

}

func (x *Context) Finalize() {
	x.Root.Report(x, FINALIZE, "Finalize", nil)
}


//----------------------------------------------------------------------
//----------------------------------------------------------------------
func (x *Context) Build(build *BuildMetadata) {
	x.Root.Report(x, CONFIGURATION, "Build", DBuild(build))
}

func (x *Context) CommandLine() {

	hostname,err := os.Hostname()
	if err != nil {
		x.Root.InternalError(WrapError(err, "Could not retrieve hostname"))
		return
	}

	u,err := user.Current()
	if err != nil {
		x.Root.InternalError(WrapError(err, "Could not retrieve user info"))
		return
	}

	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
  if err != nil {
		x.Root.InternalError(WrapError(err, "Could not retrieve program path"))
		return
  }

	prog := path.Base(os.Args[0])

	d := D{
		"Host": hostname,
		"User": u.Username,
		"Path": dir,
		"Program": prog,
		"Args": os.Args[1:],
	}

	x.Root.Report(x, CONFIGURATION, "Command line", &d)

}

func (x *Context) Environment() {

	d := D{}
	for _, e := range os.Environ() {
    pair := strings.Split(e, "=")
		d[pair[0]] = pair[1]
  }
	x.Root.Report(x, CONFIGURATION, "Environment", &d)

}

func (x *Context) Process() {

	hostname,err := os.Hostname()
	if err != nil {
		x.Root.InternalError(WrapError(err, "Could not retrieve hostname"))
		return
	}

	wd,err := os.Getwd()
  if err != nil {
		x.Root.InternalError(WrapError(err, "Could not retrieve working dir"))
		return
  }

	u,err := user.Current()
	if err != nil {
		x.Root.InternalError(WrapError(err, "Could not retrieve user info"))
		return
	}

	d := D{
		"Host": hostname,
		"WD": wd,
		"UID": u.Uid,
		"User": u.Username,
		"PID": os.Getpid(),
	}

	x.Root.Report(x, CONFIGURATION, "Process", &d)

}

//----------------------------------------------------------------------
//----------------------------------------------------------------------
func (x *Context) Info(msg string, data ...interface{}) {
	x.Root.Report(x, INFO, msg, DAggregate(data))
}

func (x *Context) Warning(msg string, data ...interface{}) {
	x.Root.Report(x, WARNING, msg, DAggregate(data))
}

func (x *Context) Error(msg string, err error, data ...interface{}) error {

	// Note that this can't/shouldn't just throw err into the data blob
	// because the standard errors package error doesn't expose
	// anything, even the message.  So you basically have to reduce to a
	// string via Error().

	x.Root.Report(x, ERROR, msg, DAggregate(data).Set("Error", err.Error()))
	return WrapError(err, msg)

}

// Failure is the same as Error but doesn't take an error object.
func (x *Context) Failure(msg string, data ...interface{}) error {
	x.Root.Report(x, ERROR, msg, DAggregate(data))
	return NewError(msg)
}

// Generally only the top level should invoke fatal, not components.
func (x *Context) Fatal(msg string, err error, data ...interface{}) {
	x.Root.Report(x, FATAL, msg, DAggregate(data).Set("Error", err.Error()))
	os.Exit(1)
}
