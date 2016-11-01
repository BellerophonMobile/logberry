package main

import (
	"github.com/BellerophonMobile/logberry"
	"gopkg.in/natefinch/lumberjack.v2"
)

// To direct output, e.g., to a file, simply attach an OutputDriver
// targeting that destination.  The built-in text and JSON outputs
// both take io.Writer instances as targets.  This example adds JSON
// output to a rolling log managed by Lumberjack, in addition to the
// default console text output.
func main() {
	logger := &lumberjack.Logger{
		Filename: "/tmp/foo.log",
	}

	// Use SetOutputDriver() to have only the managed JSON output.
	logberry.Std.AddOutputDriver(logberry.NewJSONOutput(logger))

	// Log where we're writing this log...
	logberry.Main.Info("Wrote output", &logberry.D{"Log": logger.Filename})

	logberry.Std.Stop()

}
