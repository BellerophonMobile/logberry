package env

import (
	"os"
	"os/user"
	"path"
	"path/filepath"
	"strings"
	"github.com/BellerophonMobile/logberry"
)

// BuildMetadataEvent generates a configuration log event for the
// given Task reporting the build configuration, as captured by the
// passed object.  A utility script to generate such metadata
// automatically is in the util/ folder of the Logberry repository.
func LogBuildMetadata(main *logberry.Task, build *BuildMetadata) {
	main.Event(logberry.CONFIGURATION, "Build metadata", logberry.DAggregate([]interface{}{build}))
}

// BuildSignatureEvent generates a configuration log event for the
// given Task reporting build configuration, as captured by the given
// string.  A utility script to generate such metadata automatically
// is in the util/ folder of the Logberry repository.  It can be
// useful to use this string rather than a BuildMetadata object so
// that it can be passed in through the standard go tools command
// line, i.e., via linker flags.
func LogBuildSignature(main *logberry.Task, build string) {
	main.Event(logberry.CONFIGURATION, "Build signature", logberry.D{"Signature": build})
}

// ConfigurationEvent generates a configuration log event for the
// given Task reporting parameters or other initialization data.  The
// variadic data parameter is aggregated as a D and reporting with the
// event, as is the data permanently associated with the Task.  The
// given data is not associated to the Task permanently.
func LogConfiguration(main *logberry.Task, data ...interface{}) {
	main.Event(logberry.CONFIGURATION, "Configuration", data...)
}

// CommandLineEvent generates a configuration log event for the given
// Task reporting the command line used to execute the currently
// executing process.
func LogCommandLine(main *logberry.Task) error {

	hostname, err := os.Hostname()
	if err != nil {
		return logberry.WrapError("Could not retrieve hostname", err)
	}

	u, err := user.Current()
	if err != nil {
		return logberry.WrapError("Could not retrieve user info", err)
	}

	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		return logberry.WrapError("Could not retrieve program path", err)
	}

	prog := path.Base(os.Args[0])

	d := logberry.D{
		"Host":    hostname,
		"User":    u.Username,
		"Path":    dir,
		"Program": prog,
		"Args":    os.Args[1:],
	}

	main.Event(logberry.CONFIGURATION, "Command line", d)

	return nil

}

// EnvironmentEvent generates a configuration log event for the given
// Task reporting the current operating system host environment
// variables of the currently executing process.
func LogEnvironment(main *logberry.Task) {

	d := logberry.D{}
	for _, e := range os.Environ() {
		pair := strings.Split(e, "=")
		d[pair[0]] = pair[1]
	}

	main.Event(logberry.CONFIGURATION, "Environment", d)

}

// ProcessEvent generates a configuration log event for the given Task
// reporting identifiers for the currently executing process.
func LogProcess(main *logberry.Task) error {

	hostname, err := os.Hostname()
	if err != nil {
		return logberry.WrapError("Could not retrieve hostname", err)
	}

	wd, err := os.Getwd()
	if err != nil {
		return logberry.WrapError("Could not retrieve working dir", err)
	}

	u, err := user.Current()
	if err != nil {
		return logberry.WrapError("Could not retrieve user info", err)
	}

	d := logberry.D{
		"Host": hostname,
		"WD":   wd,
		"UID":  u.Uid,
		"User": u.Username,
		"PID":  os.Getpid(),
	}

	main.Event(logberry.CONFIGURATION, "Process", d)

	return nil
	
}
