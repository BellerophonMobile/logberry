package main

import (
	"github.com/BellerophonMobile/logberry"
)

type test struct {
	StringField string
	IntField int
}


func main() {

  log := logberry.NewComponent("testcmpnt", &logberry.D{"Alpha": "Beta"})

	value := &test{ StringField: "Banana", IntField: 24 }


	task := log.Task("Some non-trivial activity", value)

	task.Success()


	task = log.ResourceTask("Read some resource", "/dev/null")

	task.Failure("Can't read /dev/null")

}
