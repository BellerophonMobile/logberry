package main

import (
	"errors"
	"github.com/BellerophonMobile/logberry"
)

func main() {

  obj := logberry.Main.Instance("Object7")

  obj.Info("Component is processing", &logberry.D{"Task": "757"})

	obj.Info("Generic message")

	obj.Info("Generic data", 7, 42, 39)

  obj.Error("Aborting processing", errors.New("CPU meltdown"))

	obj.Finalize()

}
