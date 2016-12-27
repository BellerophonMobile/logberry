package main

import (
	"errors"
	"github.com/BellerophonMobile/logberry"
)

func somework(parent *logberry.Task) error {
	task := parent.Task("Some work")
	e := errors.New("inscrutable library error")
	return task.WrapError("Mix-up in job order", e)
}

func main() {
	defer logberry.Std.Stop()

	logberry.Main.Ready()

	task := logberry.Main.Task("Compute something")
	task.AddData(logberry.D{"Priority": 17})
	task.Failure("Value out of bounds")
	
	task = logberry.Main.Task("Bigger computation")
	err := errors.New("Bad input")
	task.Error(err)

	task = logberry.Main.Task("Another computation")
	err = somework(task)
	if err != nil {
		task.Error(err)
	} else {
		task.Success()
	}	

	task = logberry.Main.Task("Final computation")
	err = logberry.NewError("Unspeakable horror", 49)
	task.Error(err)

	logberry.Main.Stopped()

}
