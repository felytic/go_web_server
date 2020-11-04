package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
)

type handlerFuncType func(http.ResponseWriter, *http.Request)

const failurePercentage = 15

func getRequestIP(r *http.Request) string {
	forwarded := r.Header.Get("X-FORWARDED-FOR")

	if forwarded != "" {
		return forwarded
	}
	return r.RemoteAddr
}

func payResponseHandler(name string, url string) handlerFuncType {
	return func(writer http.ResponseWriter, request *http.Request) {
		prefix := fmt.Sprintf("[%s] %s ", name, getRequestIP(request))

		// throw 500 error with 15% probability
		if rand.Intn(100) < failurePercentage {
			log.Print(prefix, "500 Internal Server Error")
			writer.WriteHeader(http.StatusInternalServerError)

		} else {
			log.Print(prefix, "200 OK")
			fmt.Fprintln(writer, url)
		}
	}
}

func runPaymentServer(port string, name string, url string) {
	addr := ":" + port
	mux := http.NewServeMux()
	handler := payResponseHandler(name, url)
	mux.HandleFunc("/", handler)

	server := http.Server{
		Addr:    addr,
		Handler: mux,
	}

	log.Print("Started ", name, " server at port ", port)
	log.Fatal(server.ListenAndServe())
}

func main() {
	go runPaymentServer("8090", "Apple Pay", "https://apple.com")
	go runPaymentServer("8091", "Google Pay", "https://google.com")
	go runPaymentServer("8092", "Pay Pal", "https://paypal.com")
	runPaymentServer("8093", "Stripe", "https://stripe.com")
}
