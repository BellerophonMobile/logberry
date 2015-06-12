package main

import (
	"github.com/BellerophonMobile/logberry"
)

func main() {
	
	logberry.Main.Info("Demo is functional")

	logberry.Main.Event("computation", "Computed data",
		logberry.D{"X": -234, "Label": "priority"})

	logberry.Main.Failure("Arbritrary failure")

}
