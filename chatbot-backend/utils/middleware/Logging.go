package middleware

import (
	"log"
	"net/http"
)

func Logging(next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Received %s %s\n", r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
	}
}