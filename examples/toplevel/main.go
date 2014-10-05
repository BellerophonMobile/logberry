package main

import (
	"github.com/BellerophonMobile/logberry"
)


func startEngines() error {

	logberry.Main.Info("Spooling engines",
		&logberry.D{"Block": "C", "Power": 2})

	logberry.Main.Warning("Power drain");

	return logberry.Main.Failure("Failure in hyperdrive!")

}


func main() {

	logberry.Main.Build(buildmeta)

	logberry.Main.CommandLine()

	flightdata := struct {
		Destination string
		Cargo string
	}{
		"Tau Ceti",
		"Hopes and dreams",
	}

	shipdata := struct {
		Name string
		Tier int
	}{
		"Hyperion",
		9,
	}

	logberry.Main.Info("Starting flight", flightdata, shipdata)

	e := startEngines()
	if e != nil {
		logberry.Main.Fatal("Flight aborted", e)
	}

}
