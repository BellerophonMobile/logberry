/*
Package logberry implements a structured logging framework.  It is
focused on generating logs, rather than managing them, and tries to be
lightweight while capturing more semantics and structure than is
typical, in readable and easily parsed forms.

There are four central concepts/objects/interfaces:

 D              - Data to be published with an event.
 Task           - A component, function, or logic that generates events.
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
	"io/ioutil"
)

// Std is the default Root created at startup.
var Std *Root

// Main is the default Task created at startup, roughly intended to
// represent main program execution.
var Main *Task

func init() {

	//-- Construct the standard default root
	Std = NewRoot(24)

	stdout := NewStdOutput(path.Base(os.Args[0]))

	filterfn := os.Getenv("LogberryFilter")
	if filterfn != "" {
		bytes,err := ioutil.ReadFile(filterfn)
		if err != nil {
			panic("Could not read filter file:" + err.Error())
		}
		
		filter,err := NewFilterOutput(bytes, "std", stdout)
		if err != nil {
			panic("Could not parse filter file:" + err.Error())
		}

		Std.AddOutputDriver(filter)		
	} else {
		Std.AddOutputDriver(stdout)
	}

	//-- Construct the standard default task
	Main = Std.Component("main")

	// end init
}
