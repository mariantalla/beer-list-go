package main

import (
	"fmt"
	"net/http"
	"os"
)

func main() {
	http.Handle("/", http.FileServer(http.Dir("./views")))
	http.ListenAndServe(fmt.Sprintf(":%s", os.Getenv("PORT")), nil)
}
