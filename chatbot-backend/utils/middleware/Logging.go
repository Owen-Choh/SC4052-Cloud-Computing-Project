package middleware

import (
	"log"
	"net/http"
	"time"
)

type wrappedResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (w *wrappedResponseWriter) WriteHeader(statusCode int) {
	w.statusCode = statusCode
	w.ResponseWriter.WriteHeader(statusCode)
}

func Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		log.Println(r.Header.Get("Origin"))
		wrappedResponseWriter := &wrappedResponseWriter{
			ResponseWriter: w, 
			statusCode: http.StatusOK,
		}

		next.ServeHTTP(wrappedResponseWriter, r)
		log.Printf("Received %s at %s Replied with %d %s\n", r.Method, r.URL.Path, wrappedResponseWriter.statusCode, time.Since(start))
	})
}