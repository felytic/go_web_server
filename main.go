package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"

	"./plans"

	_ "github.com/mattn/go-sqlite3"
)

const port int = 8081
const dbPath = "./db.sqlite3"

var db *sql.DB

func connectDB(path string) {
	conn, err := sql.Open("sqlite3", path)
	plans.HandleErr(err)

	db = conn
	log.Print("Connected to database")
}

func disconnectDB() {
	db.Close()
	log.Print("Disconnected from database")
}

func indexHandler(writer http.ResponseWriter, request *http.Request) {
	plans := plans.GetPlans(db)
	page_template := template.Must(template.ParseFiles("templates/plans.html"))
	page_template.Execute(writer, plans)
	log.Print("[REQUEST] ", request.URL)
}

func subscribeHandler(writer http.ResponseWriter, request *http.Request) {
	planId, err := strconv.Atoi(request.URL.Query().Get("id"))
	plans.HandleErr(err)
	plan := plans.GetPlanById(db, planId)
	fmt.Fprintln(writer, plan)
	log.Print("[REQUEST] ", request.URL)
}

func errorHandler(writer http.ResponseWriter, request *http.Request) {
	page_template := template.Must(template.ParseFiles("templates/error.html"))
	page_template.Execute(writer, nil)
	log.Print("[REQUEST] ", request.URL)
}

func main() {
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/subscribe", subscribeHandler)
	http.HandleFunc("/error", errorHandler)

	connectDB(dbPath)
	defer disconnectDB()

	log.Print("Started server at port ", port)
	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(port), nil))
}
