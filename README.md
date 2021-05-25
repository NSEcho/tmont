# tmont
Simple HTML template reloader for Go projects

# Usage

Initialize `TMonitor` struct by calling the method `New` passing the slice of filenames you would like to monitor and then call `Watch` method.

To get the template back, simply call the `Get` method providing the name of the filename.

# Example

```golang
package main

import (
	"fmt"
	"net/http"

	"github.com/lateralusd/tmont"
)

var files = []string{"/Users/myHomeDir/templaterino"}
var t *tmont.TMonitor

func main() {
	t = tmont.New(files...)
	t.Watch()

	http.HandleFunc("/", handlerino)
	http.ListenAndServe(":8000", nil)
}

func handlerino(w http.ResponseWriter, r *http.Request) {
	tpl := t.Get("/Users/myHomeDir/templaterino")
	fmt.Println("tpl is", tpl)
	tpl.Execute(w, nil)
}
```

# Note

Project is experimental, made for fun but will be used in my own projects
