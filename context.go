package logberry

import (
	"os"
	"path"
	"path/filepath"
	"strings"
	"os/user"
)


type Context struct {
	Root *Root
	Parent *Context
	Class ContextClass
	Label string
}


//----------------------------------------------------------------------
//----------------------------------------------------------------------
func NewContext(parent *Context, class ContextClass, label string) *Context {

	c := &Context {
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
