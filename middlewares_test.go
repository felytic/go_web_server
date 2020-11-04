package main

import (
	"./middlewares"
	"fmt"
	"net/http"
	"testing"
)

func dummyHandler(writer http.ResponseWriter, request *http.Request) {
	fmt.Fprintln(writer, "Plan not found")
}

type MyResponseWriter struct {
	header http.Header
}

func (w *MyResponseWriter) Header() http.Header {
	return w.header
}

func (w *MyResponseWriter) Write([]byte) (int, error) {
	return 0, nil
}

func (w *MyResponseWriter) WriteHeader(statusCode int) {
}

func TestWrap(t *testing.T) {
	newHandler := middlewares.Wrap(dummyHandler)
	writer := &MyResponseWriter{http.Header{}}
	request := new(http.Request)
	request.Host = "http://go.com"
	newHandler(writer, request)

	if writer.Header().Get("X-Server-Name") != request.Host {
		t.Errorf("Wrap() does not add X-Server-Name header")
	}
	if writer.Header().Get("X-Response-Time") == "" {
		t.Errorf("Wrap() does not add X-Server-Name header")
	}
}
