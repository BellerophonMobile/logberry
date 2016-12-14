package main

import (
	"time"
	"net/http"
	"fmt"
	"log"
	"github.com/BellerophonMobile/logberry"
	"github.com/BellerophonMobile/logberry/sseoutput"	
)

func main() {

	opts := sseoutput.Options{
		HistoryLimit: 2,
	}
	
	sse,err := sseoutput.New(&opts)
	if err != nil {
		log.Fatal(err)
	}
	
	logberry.Std.SetOutputDriver(sse)

		go func() {
			c := 0
			for {
				logberry.Main.Info("Important log data!", logberry.D{"Count": c})
				fmt.Printf("Generated %v\n", c)
				time.Sleep(2 * time.Second)
				c++
	}
	}()
	
	http.HandleFunc("/events", sse.Handler())
	http.HandleFunc("/view", Viewer)
	log.Fatal(http.ListenAndServe(":8080", nil))
	
}

func Viewer(w http.ResponseWriter, r *http.Request) {
	log.Println("Viewer")
	
	fmt.Fprint(w, `
<!DOCTYPE html>
<html>
<body>
Events:<br/>

	<script type="text/javascript">
	    var source = new EventSource('/events');
      source.onopen = function(e) {
          console.log("OnOpen:" + e)
      }
      source.onerror = function(e) {
          console.log("OnError: " + e)
      }
	    source.onmessage = function(e) {
          console.log("OnMessage:" + e.data)
	        document.body.innerHTML += e.data + '<br>';
	    }
      source.addEventListener("urgentupdate", function(e) {
          console.log("Update:" + e)
	        document.body.innerHTML += e.data + '<br>';
      });
	</script>
</body>
</html>
`)
	
}
