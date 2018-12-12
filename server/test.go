// Just a simple test, to get the Docker thing
// running properly ...

package main

import (
	"io"
	"log"
	"net/http"
)

const page string = "<h1>Hello World!</h1><p>How are you today?!?!?!??!</p>"

func servePage(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, page)
}

func main() {
	http.HandleFunc("/", servePage)
	log.Fatal(http.ListenAndServe(":80", nil))
}
