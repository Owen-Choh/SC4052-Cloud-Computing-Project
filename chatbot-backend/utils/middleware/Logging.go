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

// Implement http.Flusher if the underlying ResponseWriter is a Flusher
func (w *wrappedResponseWriter) Flush() {
	if flusher, ok := w.ResponseWriter.(http.Flusher); ok {
		flusher.Flush() // Delegate the Flush call to the underlying Flusher
	}
	// If the underlying ResponseWriter is not a Flusher
}

func Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		wrappedResponseWriter := &wrappedResponseWriter{
			ResponseWriter: w, 
			statusCode: http.StatusOK,
		}

		next.ServeHTTP(wrappedResponseWriter, r)
		log.Printf("Received %s at %s Replied with %d %s\n", r.Method, r.URL.Path, wrappedResponseWriter.statusCode, time.Since(start))
	})
}