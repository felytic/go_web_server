package middlewares

import (
	"log"
	"net/http"
	"time"
)

type writerWithTimer struct {
	http.ResponseWriter
	start time.Time
}

func Wrap(handler http.HandlerFunc) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		// Substitude original writer with responseWriterWithTimer
		timedWriter := &writerWithTimer{writer, time.Now()}
		timedWriter.Header().Set("X-Server-Name", request.Host)
		handler.ServeHTTP(timedWriter, request)

		duration := writer.Header().Get("X-Response-Time")
		log.Printf("[REQUEST] %s FROM: %s DURATION: %s", request.URL, request.RemoteAddr, duration)
	}
}

// Adds X-Response-Time header to each response
func (w *writerWithTimer) Write(b []byte) (int, error) {
	duration := time.Now().Sub(w.start)
	w.Header().Set("X-Response-Time", duration.String())

	return w.ResponseWriter.Write(b)
}
