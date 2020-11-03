package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
)

const port = "8090"

type UrlResponse struct {
	Url  string `json:"url"`
	Name string `json:"name"`
}

func GetIP(r *http.Request) string {
	forwarded := r.Header.Get("X-FORWARDED-FOR")
	if forwarded != "" {
		return forwarded
	}
	return r.RemoteAddr
}

func appleHandler(writer http.ResponseWriter, request *http.Request) {
	appleResponse := UrlResponse{
		Url:  "https://apple.com",
		Name: "Apple Pay",
	}

	respJSON, err := json.Marshal(appleResponse)

	if err != nil {
		log.Print("[ERROR]", err)
		fmt.Fprintln(writer, err)
	}

	// throw 500 error with 30% probability
	if rand.Intn(10) < 3 {
		log.Print("[", GetIP(request), "] 500 Internal Server Error")
		writer.WriteHeader(http.StatusInternalServerError)
	} else {
		log.Print("[", GetIP(request), "] 200 OK")
		writer.Header().Set("Content-Type", "application/json")
		writer.Write(respJSON)
	}
}

func main() {
	http.HandleFunc("/", appleHandler)
	log.Print("Started server at port ", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
