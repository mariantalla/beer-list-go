package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

type Env struct {
	GATag         string
	InstanceIndex string
	Beers         map[string]string
}

func main() {

	db, err := sql.Open("mysql",
		"3AG4rpTntTHg4z96:S5paRHY1Ez5EaI6f@mysql-broker.local.pcfdev.io:3306/cf_ce575afb_bec6_4a80_83cb_de4c1d60c543")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	http.HandleFunc("/", Index)
	http.ListenAndServe(fmt.Sprintf(":%s", os.Getenv("PORT")), nil)
}

func Index(w http.ResponseWriter, req *http.Request) {
	beers := map[string]string{
		"awesome-place": "awesome-beer",
	}
	env := Env{
		GATag:         os.Getenv("GA_TAG"),
		InstanceIndex: os.Getenv("INSTANCE_INDEX"),
		Beers:         beers,
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
