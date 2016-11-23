package main

import "net/http"

func main() {
	http.Handle("/", http.FileServer(http.Dir("./views")))
	http.ListenAndServe(":4000", nil)
}
