package main

import (
	"github.com/BellerophonMobile/logberry"
)

func main() {

	logberry.Main.Info("Demo is functional")

	logberry.Main.Info("Computed data", logberry.D{"X": 23, "Label": "priority"})

	logberry.Main.Failure("Arbritrary failure")

}
