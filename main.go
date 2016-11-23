package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"

	cfenv "github.com/cloudfoundry-community/go-cfenv"
	_ "github.com/go-sql-driver/mysql"
)

type Env struct {
	GATag         string
	InstanceIndex string
	Beers         map[string]string
}

var (
	env  *cfenv.App
	db   *sql.DB
	file *os.File
	err  error
)

func main() {
	file, err = os.Create("troublesome-file")
	check(err)
	defer file.Close()

	db = connectToDB()
	defer db.Close()
	renderCssAndImages()
	http.HandleFunc("/", Index)
	http.HandleFunc("/create", Create)
	http.ListenAndServe(fmt.Sprintf(":%s", os.Getenv("PORT")), nil)
}

func connectToDB() *sql.DB {
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

	_, err = db.Exec("create table if not exists beers (id integer PRIMARY KEY AUTO_INCREMENT, region varchar(255) NOT NULL, brand varchar(255) NOT NULL)")
	check(err)
	return db
}

func Create(w http.ResponseWriter, req *http.Request) {
	file.WriteString("Someone added a new beer!\n")
	check(req.ParseForm())
	brand := req.PostForm.Get("brand")
	region := req.PostForm.Get("region")
	fmt.Printf("creating beer: %s/%s\n", region, brand)

	_, err := db.Exec(fmt.Sprintf("insert into beers (region, brand) values ('%s', '%s')", region, brand))
	check(err)

	http.Redirect(w, req, "/", http.StatusFound)
}

func Index(w http.ResponseWriter, req *http.Request) {
	file.WriteString("Someone requested beers!\n")

	beers := make(map[string]string)

	rows, err := db.Query("select region, brand from beers")
	check(err)

	for rows.Next() {

		var region string
		var brand string
		check(rows.Scan(&region, &brand))
		fmt.Printf("found beer: %s/%s\n", region, brand)
		beers[region] = brand
	}

	env := Env{
		GATag:         os.Getenv("GA_TAG"),
		InstanceIndex: os.Getenv("INSTANCE_INDEX"),
		Beers:         beers,
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
