package main

import (
	"errors"
	"github.com/BellerophonMobile/logberry"
)

func main() {
	logberry.AddOutput(logberry.NewStdOutput())

  logberry.Build(buildmeta)
  logberry.CommandLine()

  log := logberry.NewComponent("testcmpnt")

  log.Info("The bananas have gotten loose!", &logberry.Data{"Fruit": "bananas"})
  log.Error("Could not launch", errors.New("Failure in hyperdrive!"))
}
