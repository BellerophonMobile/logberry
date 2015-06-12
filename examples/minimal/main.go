package main

import (
	"github.com/BellerophonMobile/logberry"
)

func somecomputation(data interface{}) (int, error) {
	return 7, nil
}

func efunc() error {
	return logberry.Main.Failure("Unrecoverable error",
		logberry.D{"Note":"Arbitrary error!"})
}

func main() {

	// Uncomment for JSON output
	// logberry.Std.SetOutputDriver(logberry.NewJSONOutput(os.Stdout))
	
	logberry.Main.BuildMetadata(buildmeta)

	logberry.Main.Info("Demo is functional")


	// Create some structured application data
	var data = struct {
		DataLabel string
		DataInt    int
	}{"alpha", 9}

	// Do some activity on that data, which may fail, within the component
	task := logberry.Main.SubTask("Compute numbers", &data)
	res, err := somecomputation(data)
	if err != nil {
		task.Error(err)
		return
	}
	task.Success(logberry.D{"Result": res})

	// An error has occurred out of nowhere!
	if e := efunc(); e != nil {
		return
	}

	logberry.Main.Info("Done")

}
