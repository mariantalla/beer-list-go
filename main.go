package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
)

type Env struct {
	GATag string
	InstanceIndex string
}

func main() {
	http.HandleFunc("/", Index)
	http.ListenAndServe(fmt.Sprintf(":%s", os.Getenv("PORT")), nil)
}

func Index(w http.ResponseWriter, req *http.Request) {
	env := Env{
		GATag: os.Getenv("GA_TAG"),
		InstanceIndex: os.Getenv("INSTANCE_INDEX"),
	}
	render(w, "views/index.html", env)
}

func render(w http.ResponseWriter, tmpl string, env Env) {
	tmpl = fmt.Sprintf("%s", tmpl)
	t, err := template.ParseFiles(tmpl)
	if err != nil {
		log.Print("template parsing error: ", err)
	}
	err = t.Execute(w, env)
	if err != nil {
		log.Print("template executing error: ", err)
	}
}
