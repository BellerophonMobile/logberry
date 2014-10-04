package main

import (
//	"errors"
	"github.com/BellerophonMobile/logberry"
)

func main() {

	logberry.Main.Build(buildmeta)

	logberry.Main.CommandLine()

	logberry.Main.Environment()

	logberry.Main.Process()

//  logberry.Info("Bananas are loose!", &logberry.D{"Status": "insane"})

//  logberry.Info("Continuing on", &logberry.D{"Code": "Red", "Power": 2})

//	logberry.Warning("Power drain");

//  logberry.Error("Could not launch", errors.New("Failure in hyperdrive!"))

}

/*
 * Verbosity controllable down a hierarchy branch?
 * Unique IDs for Contexts
 * Actions
 */
