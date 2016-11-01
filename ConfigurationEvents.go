package logberry

import (
	"os"
	"os/user"
	"path"
	"path/filepath"
	"strings"
)

var contexttask = Task{
	uid: 0,
	component: "[process]",
}

// BuildMetadataEvent generates a configuration log event for the
// given Task reporting the build configuration, as captured by the
// passed object.  A utility script to generate such metadata
// automatically is in the util/ folder of the Logberry repository.
func BuildMetadataEvent(main *Task, build *BuildMetadata) {
	main.root.event(main, CONFIGURATION, "Build metadata", DAggregate([]interface{}{build}))
}

// BuildSignatureEvent generates a configuration log event for the
// given Task reporting build configuration, as captured by the given
// string.  A utility script to generate such metadata automatically
// is in the util/ folder of the Logberry repository.  It can be
// useful to use this string rather than a BuildMetadata object so
// that it can be passed in through the standard go tools command
// line, i.e., via linker flags.
func BuildSignatureEvent(main *Task, build string) {
	main.root.event(main, CONFIGURATION, "Build signature", D{"Signature": build})
}

// ConfigurationEvent generates a configuration log event for the
// given Task reporting parameters or other initialization data.  The
// variadic data parameter is aggregated as a D and reporting with the
// event, as is the data permanently associated with the Task.  The
// given data is not associated to the Task permanently.
func ConfigurationEvent(main *Task, data ...interface{}) {
	d := DAggregate(append(data, main.data))
	main.root.event(main, CONFIGURATION, "Configuration", d)
}

// CommandLineEvent generates a configuration log event for the given
// Task reporting the command line used to execute the currently
// executing process.
func CommandLineEvent(main *Task) {

	hostname, err := os.Hostname()
	if err != nil {
		main.root.internalerror(WrapError("Could not retrieve hostname", err))
		return
	}

	u, err := user.Current()
	if err != nil {
		main.root.internalerror(WrapError("Could not retrieve user info", err))
		return
	}

	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		main.root.internalerror(WrapError("Could not retrieve program path", err))
		return
	}

	prog := path.Base(os.Args[0])

	d := D{
		"Host":    hostname,
		"User":    u.Username,
		"Path":    dir,
		"Program": prog,
		"Args":    os.Args[1:],
	}
	d.CopyFrom(main.data)

	main.root.event(main, CONFIGURATION, "Command line", d)

}

// EnvironmentEvent generates a configuration log event for the given
// Task reporting the current operating system host environment
// variables of the currently executing process.
func EnvironmentEvent(main *Task) {

	d := D{}
	for _, e := range os.Environ() {
		pair := strings.Split(e, "=")
		d[pair[0]] = pair[1]
	}
	d.CopyFrom(main.data)

	main.root.event(main, CONFIGURATION, "Environment", d)

}

// ProcessEvent generates a configuration log event for the given Task
// reporting identifiers for the currently executing process.
func ProcessEvent(main *Task) {

	hostname, err := os.Hostname()
	if err != nil {
		main.root.internalerror(WrapError("Could not retrieve hostname", err))
		return
	}

	wd, err := os.Getwd()
	if err != nil {
		main.root.internalerror(WrapError("Could not retrieve working dir", err))
		return
	}

	u, err := user.Current()
	if err != nil {
		main.root.internalerror(WrapError("Could not retrieve user info", err))
		return
	}

	d := D{
		"Host": hostname,
		"WD":   wd,
		"UID":  u.Uid,
		"User": u.Username,
		"PID":  os.Getpid(),
	}
	d.CopyFrom(main.data)

	main.root.event(main, CONFIGURATION, "Process", d)

}
