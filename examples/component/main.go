package main

import (
	"errors"
	"github.com/BellerophonMobile/logberry"
)

func main() {

  log := logberry.NewComponent("testcmpnt", &logberry.D{"Rolling": "basic"})

  log.Info("Component is processing", &logberry.D{"ID": "757"})

	log.Info("Generically report a message")

	log.Info("Generically report some data", 7, 42, 39)

	log.Service("Contacted some resource", &logberry.D{ "Host": "localhost:9" })

  log.Error("Aborting processing", errors.New("Abnormal event!"))

}
