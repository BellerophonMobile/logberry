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
	logberry.SetOutputDriver(m)

  log := logberry.NewComponent("testcmpnt", &logberry.D{"Rolling": "basic"})

	log.Info("Done")
}
