package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
)

func main() {
	http.HandleFunc("/", Index) // http.Handle("/", http.FileServer(http.Dir("./views")))
	http.ListenAndServe(fmt.Sprintf(":%s", os.Getenv("PORT")), nil)
}

func Index(w http.ResponseWriter, req *http.Request) {
	render(w, "views/index.html")
}

func render(w http.ResponseWriter, tmpl string) {
	tmpl = fmt.Sprintf("%s", tmpl)
	t, err := template.ParseFiles(tmpl)
	if err != nil {
		log.Print("template parsing error: ", err)
	}
	err = t.Execute(w, "")
	if err != nil {
		log.Print("template executing error: ", err)
	}
}
