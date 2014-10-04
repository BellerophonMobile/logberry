package main

import (
	"errors"
	"github.com/BellerophonMobile/logberry"
)

func main() {

  cmp := logberry.Main.Component("testcmpnt", &logberry.D{"Mode": "basic"})

  cmp.Info("Component is processing", &logberry.D{"Task": "757"})

	cmp.Info("Generic message")

	cmp.Info("Generic data", 7, 42, 39)

  cmp.Error("Aborting processing", errors.New("CPU meltdown"))

	cmp.Cleanup()

}
