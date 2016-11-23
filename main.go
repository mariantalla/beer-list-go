package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
)

type Env struct {
	GATag         string
	InstanceIndex string
}

func main() {
	renderCssAndImages()
	http.HandleFunc("/", Index)
	http.ListenAndServe(fmt.Sprintf(":%s", os.Getenv("PORT")), nil)
}

func Index(w http.ResponseWriter, req *http.Request) {
	env := Env{
		GATag:         os.Getenv("GA_TAG"),
		InstanceIndex: os.Getenv("INSTANCE_INDEX"),
	}
	render(w, "views/index.html", env)
}

func render(w http.ResponseWriter, htmlTemplate string, env Env) {
	htmlTemplate = fmt.Sprintf("%s", htmlTemplate)
	template, err := template.ParseFiles(htmlTemplate)
	check(err)

	err = template.Execute(w, env)
	check(err)
}

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func renderCssAndImages() {
	http.Handle("/views/", http.StripPrefix("/views/", http.FileServer(http.Dir("views"))))
}
