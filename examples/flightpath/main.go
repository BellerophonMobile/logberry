package main

import (
	"github.com/BellerophonMobile/logberry"
)


type destination struct {
	Planet string
	Priority int
}

var enginelog *logberry.Component

func route(dest destination) ([]byte,error) {
	return []byte{0xDE, 0xAD, 0xBE, 0xEF}, logberry.NewError("Unstable route")
}

func startengines(dest destination) error {

	enginelog = logberry.Main.Component("engines")

	enginelog.Info("Spooling engines", &logberry.D{"Block": "C", "Power": 2})

	task := enginelog.CalculationTask("Make flightpath", dest)
	_, err := route(dest)
	if err != nil {
		return task.Error(err)
	}
	task.Success()

	enginelog.Info("Ignition")

	return nil
}

func monitorengines() {
	enginelog.Warning("Low power levels", "80%");
}

func stopengines() error {

	coildevice := "/firefly/engine/coil/3"
	task := enginelog.ResourceTask("Writing coil parameters", coildevice)
	return task.Failure("Unknown coil")

}


func main() {
	var target = destination{"Hera", 7}

	logberry.Main.Build(buildmeta)
	logberry.Main.Info("Flight computer boot")

	err := startengines(target)
	if err != nil {
		logberry.Main.Recovered("No flightpath, flying straight!", err)
	}

	monitorengines()

	err = stopengines()
	if err != nil {
		logberry.Main.Fatal("Crash landing!", err)
	}

}
