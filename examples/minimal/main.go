package main

import (
	"errors"
	"github.com/BellerophonMobile/logberry"
)

func main() {

	logberry.AddOutput(logberry.NewStdOutput())

  logberry.Build(buildmeta)
  logberry.CommandLine()

  logberry.Info("The bananas have gotten loose!",
		&logberry.D{"Status": "insane"})

  logberry.Error("Could not launch", errors.New("Failure in hyperdrive!"))

	logberry.Warning("Power drain");

  logberry.Info("Continuing on", &logberry.D{"Code": "Red"})

}
