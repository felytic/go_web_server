package main

import (
	"context"
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"time"

	"./payments"
	"./plans"

	"github.com/joho/godotenv"
	_ "github.com/mattn/go-sqlite3"
)

var dbHost string
var port string
var shutdownTimeout time.Duration

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

func isAllResponsesOk(responses []payments.PaymentResponse) bool {
	for _, response := range responses {
		if response.Err != nil {
			return false
		}
		if response.Res.StatusCode != 200 {
			return false
		}
	}

	return true
}

func subscribeHandler(writer http.ResponseWriter, request *http.Request) {
	planId, err := strconv.Atoi(request.URL.Query().Get("id"))

	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	_, err = plans.GetPlanById(db, planId)

	if err != nil {
		writer.WriteHeader(http.StatusNotFound)
		fmt.Fprintln(writer, "Plan not found")
		return
	}

	responses := payments.GetPaymentResponses(planId)
	if isAllResponsesOk(responses) {
		fmt.Fprintln(writer, responses)
	} else {
		page_template := template.Must(template.ParseFiles("templates/error.html"))
		page_template.Execute(writer, nil)
	}

	log.Print("[REQUEST] ", request.URL)
}

func initEnv() {
	err := godotenv.Load(".env")

	port = os.Getenv("APP_PORT")
	dbHost = os.Getenv("DB_HOST")
	duration, _ := strconv.Atoi(os.Getenv("APP_SHUTDOWN_TIMEOUT"))
	shutdownTimeout = time.Duration(duration) * time.Second

	if port == "" || dbHost == "" || err != nil {
		log.Fatal("Error loading environment variables")
	}
}

func createServer(addr string) *http.Server {
	connectDB(dbHost)
	mux := http.NewServeMux()
	mux.HandleFunc("/", indexHandler)
	mux.HandleFunc("/subscribe", subscribeHandler)

	return &http.Server{Addr: addr, Handler: mux}
}

func runServer(server *http.Server) {
	// Run server in gorutine
	go func() {
		if err := server.ListenAndServe(); err != nil {
			log.Fatal(err)
		}
	}()
	log.Print("Started server at port ", port)
}

func stopServer(server *http.Server) {
	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatal(err)
	}
	log.Print("Server stopped succesfully")

	disconnectDB()
}

func main() {
	// Read environmet variables
	initEnv()

	addr := ":" + port

	server := createServer(addr)

	runServer(server)

	// Wait for interrupt signal
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)
	<-interrupt

	stopServer(server)
}
