package main

import (
	"github.com/BellerophonMobile/logberry"
	"os"
)

func main() {

	m := logberry.NewFanOutput()
	m.AddOutputDriver(logberry.NewStdOutput())
	m.AddOutputDriver(logberry.NewJSONOutput(os.Stdout))
	logberry.Std.SetOutputDriver(m)

	log := logberry.Main.Component("testcmpnt", &logberry.D{"Rolling": "Thunder"})

	log.Info("Done")
}
