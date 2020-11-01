package main

import (
	"./plans"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strconv"

	_ "github.com/mattn/go-sqlite3"
)

const port int = 8081

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

func processRequest(writer http.ResponseWriter, request *http.Request) {
	plans := plans.GetPlans(db)
	fmt.Fprintln(writer, plans)
	log.Print("[REQUEST] ", request.URL)
}

func main() {
	http.HandleFunc("/", processRequest)

	connectDB("./db.sqlite3")
	defer disconnectDB()

	log.Print("Started server at port ", port)
	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(port), nil))
}
