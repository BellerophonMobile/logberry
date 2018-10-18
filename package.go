/*
Package logberry implements a structured logging framework.  It is
focused on generating logs, rather than managing them, and tries to be
lightweight while capturing more semantics and structure than is
typical, in readable and easily parsed forms.

There are four central concepts/objects/interfaces:

 D              - Data to be published with an event.
 Task           - A component, function, or logic path generating events.
 OutputDriver   - Serializer for publishing events.
 Root           - An interface between Tasks and OutputDrivers.

Also important are two less fundamental but included concepts/objects:

 Error          - A generic structured error report.
 BuildMetadata  - A simple representation of the build environment.

Higher level documentation is available from the repository and README:
  https://github.com/BellerophonMobile/logberry

*/
package logberry

import (
	"os"
	"path"
)

// Std is the default Root created at startup.
var Std *Root

// Main is the default Task created at startup, roughly intended to
// represent main program execution.
var Main *Task

func init() {

	//-- Construct the standard default root
	Std = NewRoot(24)
	Std.AddOutputDriver(NewStdOutput(path.Base(os.Args[0])))
	Std.AddErrorListener(new(StdErrorListener))

	//-- Construct the standard default task manually so no event
	Main = &Task{
		uid:       newtaskuid(),
		component: "main",
		activity:  "Component main",
		root:      Std,
	}

}
