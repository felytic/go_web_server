package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
)

const failurePercentage = 15

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

type handlerFuncType func(http.ResponseWriter, *http.Request)

func payResponseHandler(name string, url string) handlerFuncType {
	payResponse := UrlResponse{
		Url:  url,
		Name: name,
	}

	respJSON, err := json.Marshal(payResponse)

	return func(writer http.ResponseWriter, request *http.Request) {
		if err != nil {
			log.Print("[ERROR]", err)
			fmt.Fprintln(writer, err)
		}

		// throw 500 error with 15% probability
		if rand.Intn(100) < failurePercentage {
			log.Print("[", GetIP(request), "] 500 Internal Server Error")
			writer.WriteHeader(http.StatusInternalServerError)
		} else {
			log.Print("[", GetIP(request), "] 200 OK")
			writer.Header().Set("Content-Type", "application/json")
			writer.Write(respJSON)
		}
	}

}

func runPaymentServer(port string, name string, url string) {
	mux := http.NewServeMux()
	handler := payResponseHandler(name, url)
	mux.HandleFunc("/", handler)

	addr := ":" + port
	server := http.Server{
		Addr:    addr,
		Handler: mux,
	}
	log.Print("Started ", name, "server at port ", port)
	log.Fatal(server.ListenAndServe())

}

func main() {
	go runPaymentServer("8090", "Apple Pay", "https://apple.com")
	go runPaymentServer("8091", "Google Pay", "https://google.com")
	go runPaymentServer("8092", "Pay Pal", "https://paypal.com")
	runPaymentServer("8093", "Stripe", "https://stripe.com")
}
