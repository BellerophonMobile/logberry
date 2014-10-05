package main

import (
	"github.com/BellerophonMobile/logberry"
)

var target = struct {
	Destination string
	Priority int
}{
	"Tau Ceti",
	1,
}

func main() {

	logberry.Main.Build(buildmeta)

	logberry.Main.CommandLine()

	logberry.Main.Info("Spooling engines", &logberry.D{"Block": "C", "Power": 2})
	logberry.Main.Info("Calculating flightpath", target)
	logberry.Main.Info("Ignition")

	logberry.Main.Warning("Power drain");

	e := logberry.Main.Failure("Failure in hyperdrive!")

	logberry.Main.Fatal("Aborting flight", e)

}
