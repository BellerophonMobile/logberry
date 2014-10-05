package main

import (
	"github.com/BellerophonMobile/logberry"
	"io/ioutil"
	"net/http"
)


func geticon() error {

	url := "https://raw.githubusercontent.com/BellerophonMobile/logberry/master/docs/logberry.png"

	get := logberry.Main.LongResourceTask("Download strawberry icon", url)
	res, err := http.Get(url)
	if err != nil {
		return get.Error(err)
	} else if res.StatusCode != http.StatusOK {
		return get.Failure(http.StatusText(res.StatusCode))
	}
	get.Success()

	return nil

}


func main() {

	myfilename := "/home/nouser/doesnotexist"
	read := logberry.Main.ResourceTask("Read app data", myfilename)
	if _,err := ioutil.ReadFile(myfilename); err != nil {
		read.Error(err)
	} else {
		read.Success()
	}

	if e := geticon(); e != nil {
		logberry.Main.Error("Could not get icon", e)
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
