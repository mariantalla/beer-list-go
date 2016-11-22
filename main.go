package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"

	"github.com/cloudfoundry-community/go-cfenv"
	_ "github.com/go-sql-driver/mysql"
)

type Env struct {
	GATag         string
	InstanceIndex string
	Beers         map[string]string
}

var (
	env *cfenv.App
	db  *sql.DB
)

func main() {
	env, _ = cfenv.Current()
	mysqlService, err := env.Services.WithName("db")
	check(err)

	credentials := mysqlService.Credentials

	dns := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s",
		credentials["username"].(string),
		credentials["password"].(string),
		credentials["hostname"].(string),
		int(credentials["port"].(float64)),
		credentials["name"].(string),
	)
	db, err = sql.Open("mysql", dns)
	check(err)
	defer db.Close()

	_, err = db.Exec("create table if not exists beers (id integer PRIMARY KEY AUTO_INCREMENT, region varchar(255) NOT NULL, brand varchar(255) NOT NULL)")
	check(err)

	http.Handle("/views/", http.StripPrefix("/views/", http.FileServer(http.Dir("views"))))
	http.HandleFunc("/", Index)
	http.HandleFunc("/create", Create)
	http.ListenAndServe(fmt.Sprintf(":%s", os.Getenv("PORT")), nil)
}

func Index(w http.ResponseWriter, req *http.Request) {
	beers := make(map[string]string)

	rows, err := db.Query("select region, brand from beers")
	check(err)

	for rows.Next() {
		var region string
		var brand string
		check(rows.Scan(&region, &brand))
		beers[region] = brand
	}

	customEnv := Env{
		GATag:         os.Getenv("GA_TAG"),
		InstanceIndex: os.Getenv("INSTANCE_INDEX"),
		Beers:         beers,
	}
	render(w, "views/index.html", customEnv)
}

func Create(w http.ResponseWriter, req *http.Request) {
	check(req.ParseForm())
	brand := req.PostForm.Get("brand")
	region := req.PostForm.Get("region")

	_, err := db.Exec(fmt.Sprintf("insert into beers (region, brand) values ('%s', '%s')", region, brand))
	check(err)

	http.Redirect(w, req, "/", http.StatusFound)
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

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
