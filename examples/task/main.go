package main

import (
	"github.com/BellerophonMobile/logberry"
	"io/ioutil"
)


func main() {

	myfilename := "/home/nouser/doesnotexist"
	read := logberry.Main.Task("Read app data", myfilename)
	if _,err := ioutil.ReadFile(myfilename); err != nil {
		read.Error(err)
	} else {
		read.Success()
	}

/*

	var value = struct{
		StringField string
		IntField int
	}{
		StringField: "Banana",
		IntField: 24,
	}

	background := logberry.Main.Task("A quick task", value)

	task := logberry.Main.Task("Some non-trivial activity", value)

	task.SetData("User", "joe")

	task.Success()


	task = logberry.Main.LongTask("A longer task")
	task.Success()

*/
}
