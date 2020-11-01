package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strconv"

	_ "github.com/mattn/go-sqlite3"
)

const port int = 8081

var db *sql.DB

type Plan struct {
	Id    int
	Name  string
	Price float32
}

func connectDB(path string) {
	conn, err := sql.Open("sqlite3", path)
	if err != nil {
		log.Fatal(err)
	}

	db = conn
	log.Print("Connected to database")
}

func disconnectDB() {
	db.Close()
	log.Print("Disconnected from database")
}

func getPlans() []Plan {
	query := `SELECT id, name, price FROM plan`

	rows, err := db.Query(query)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var plans []Plan

	for rows.Next() {
		plan := Plan{}
		err = rows.Scan(&plan.Id, &plan.Name, &plan.Price)
		if err != nil {
			log.Fatal(err)
		}

		plans = append(plans, plan)
	}

	return plans
}

func processRequest(writer http.ResponseWriter, request *http.Request) {
	plans := getPlans()
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
