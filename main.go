package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"strconv"

	"./plans"

	"github.com/joho/godotenv"
	_ "github.com/mattn/go-sqlite3"
)

var dbHost string
var port string
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

func initEnv() {
	err := godotenv.Load(".env")

	port = os.Getenv("APP_PORT")
	dbHost = os.Getenv("DB_HOST")

	if port == "" || dbHost == "" || err != nil {
		log.Fatal("Error loading environment variables")
	}
}

func main() {
	initEnv()

	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/subscribe", subscribeHandler)
	http.HandleFunc("/error", errorHandler)

	connectDB(dbHost)
	defer disconnectDB()

	log.Print("Started server at port ", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
