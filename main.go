package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/cloudfoundry-community/go-cfenv"
	_ "github.com/go-sql-driver/mysql"
)

type Env struct {
	GATag         string
	InstanceIndex string
	Beers         map[string]string
}

var env *cfenv.App

func main() {
	http.HandleFunc("/", Index)
	http.ListenAndServe(fmt.Sprintf(":%s", os.Getenv("PORT")), nil)
}

func Index(w http.ResponseWriter, req *http.Request) {
	beers := map[string]string{
		"awesome-place": "awesome-beer",
	}

	// beers := make(map[string]string)
	//
	env, _ = cfenv.Current()
	mysqlService, err := env.Services.WithName("db")
	check(err)

	credentials := mysqlService.Credentials

	fmt.Println(strings.Split(credentials["uri"].(string), "//")[1])
	db, err := sql.Open("mysql", "CLQnXGd1Yoe72dME:JPVgiVLipd5FtFcq@tcp(mysql-broker.local.pcfdev.io:3306)/cf_6a53c48b_291b_405f_86e8_cbb21352f660")
	check(err)
	defer db.Close()

	_, err = db.Exec("create table if not exists beers (id integer PRIMARY KEY AUTO_INCREMENT, region varchar(255) NOT NULL, brand varchar(255) NOT NULL)")
	check(err)
	// _, err = db.Exec("insert into beers (region, brand) values ('db-region', 'db-beer')")
	// check(err)

	rows, err := db.Query("select region, brand from beers")
	check(err)

	fmt.Println("got some results supposedly")

	for rows.Next() {
		fmt.Println("oh dear")
		var region sql.NullString
		var brand sql.NullString
		check(rows.Scan(&region, &brand))
		fmt.Printf("new row: %s %s", region.String, brand.String)
		beers[region.String] = brand.String
	}

	customEnv := Env{
		GATag:         os.Getenv("GA_TAG"),
		InstanceIndex: os.Getenv("INSTANCE_INDEX"),
		Beers:         beers,
	}
	render(w, "views/index.html", customEnv)
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
