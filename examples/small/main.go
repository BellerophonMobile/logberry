package main

//go:generate go run ../../env/util/build-metadata.go -workspace=../../ -out=build

import (
	"github.com/BellerophonMobile/logberry"
	"github.com/BellerophonMobile/logberry/env"	
	//	"os"
)

func main() {

	// Uncomment this and "os" import for JSON output
	// logberry.Std.SetOutputDriver(logberry.NewJSONOutput(os.Stdout))

	// Report build information; a script generates buildmetadata
	env.LogBuildMetadata(logberry.Main, buildmetadata)

	// Report that the program is initialized & running
	logberry.Main.Ready()

	// Create some structured application data and log it
	var data = struct {
		DataLabel string
		DataInt   int
	}{"alpha", 9}

	logberry.Main.Info("Reporting some data", data)

	// Create a program component---a long-running, multi-use entity.
	computerlog := logberry.Main.Component("computer")

	// Execute a task within that component, which may fail
	task := computerlog.Task("Compute numbers", &data)
	res, err := somecomputation()
	if err != nil {
		task.Error(err)
		return
	}
	task.Success(logberry.D{"Result": res})

	// Generate an application specific event reporting some other data
	var req = struct {
		User string
	}{"tjkopena"}

	computerlog.Event("request", "Received request", req)

	// Run a function under the component
	if e := arbitraryfunc(computerlog); e != nil {
		// Handle the error here
	}

	// The component ends
	computerlog.Finalized()

	// The program shuts down
	logberry.Main.Stopped()

	// Wait for all log messages to be output
	logberry.Std.Stop()

}

func somecomputation() (int, error) {
	return 7, nil // errors.New("Could not compute")
	// return 7, errors.New("Could not compute")
}

func arbitraryfunc(component *logberry.Task) error {

	// Start a long-running task, using Begin() to log start & begin timer
	task := component.Task("Arbitrary computation")

	// Report some intermediate progress
	task.Info("Intermediate progress", logberry.D{"Best": 9})

	// An error has occurred out of nowhere!  Log & return an error
	// noting that this task has failed, data associated with the error,
	// wrapping the underlying cause, and noting this source location
	return task.Failure("Random unrecoverable error",
		logberry.D{"Bounds": "x-axis"})

}
