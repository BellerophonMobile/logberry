package main

import (
	"github.com/BellerophonMobile/logberry"
	"os"
)

func main() {

	// This is not necessary, logberry's core supports multiple outputs.
	// But this permits you to do things like feed the core output to a
	// ThreadSafeOutput, which in turns feeds to a multiplexer and then
	// through that to several outputs.
	m := logberry.NewMultiplexerOutput()
	m.AddOutputDriver(logberry.NewStdOutput())
	m.AddOutputDriver(logberry.NewJSONOutput(os.Stdout))
	logberry.Std.SetOutputDriver(m)

	log := logberry.Main.Component("testcmpnt", &logberry.D{"Rolling": "Thunder"})

	log.Info("Done")
}
