package main

import (
	"net/http"

	"google.golang.org/appengine"
)

func main() {
	http.HandleFunc("/", indexHandler)
	appengine.Main()
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("OK"))
}
