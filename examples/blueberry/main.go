package main

import (
	"bytes"
	"github.com/BellerophonMobile/logberry"
	"io/ioutil"
	"net/http"
)

type swap struct {
	old string
	new string
}

var swaps = [...]swap{
	swap{"fill:#ed1c24;", "fill:#1c7bed;"},
	swap{"fill:#bc151b;", "fill:#153dbc;"},
	swap{"fill:#f6836c;", "fill:#6cc9f6;"},
	swap{"fill:#870e12;", "fill:#150e87;"},
}

func geticon(cxt logberry.Context, url string) ([]byte, error) {

	task := cxt.LongResourceTask("Download icon", url)
	res, err := http.Get(url)
	if err != nil {
		// Can fail via basic IO error
		return nil, task.Error(err)
	} else if res.StatusCode != http.StatusOK {
		// Can successfully get a reponse, but still fail to get icon
		return nil, task.Failure(http.StatusText(res.StatusCode))
	}
	defer res.Body.Close()

	bytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, task.Error(err)
	}

	return bytes, task.Success()

	// end geticon
}

func main() {

	url := "https://raw.githubusercontent.com/BellerophonMobile/logberry/master/docs/logberry.svg"

	icon, err := geticon(logberry.Main, url)
	if err != nil {
		logberry.Main.Fatal("Could not download icon", err)
	}

	// Super inefficient.  Bare with it for this example...
	for i := range swaps {
		icon = bytes.Replace(icon, []byte(swaps[i].old), []byte(swaps[i].new), -1)
	}

	err = ioutil.WriteFile("test.svg", icon, 0644)
	if err != nil {
		logberry.Main.Fatal("Could not write new icon", err)
	}

	logberry.Main.Info("All done!")

	// end main
}
