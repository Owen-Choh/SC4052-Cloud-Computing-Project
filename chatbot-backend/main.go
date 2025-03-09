package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/Owen-Choh/SC4052-Cloud-Computing-Assignment-2/chatbot-backend/chatbot/db"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	isInitialised, dberr := db.InitDB()
	if isInitialised {
		if dberr != nil {
			log.Printf("Database initialised but %s", dberr)
		} else {
			log.Println("Database initialised")
		}
	} else if dberr != nil {
		log.Println("Abort server start up due to error initalising database")
		return
	}
	
	router := SetUpMuxWithRoutes()

	routerWithCORS := CORS(router)
	// set server and start
	server := http.Server{
		Addr:    ":8080",
		Handler: routerWithCORS,
	}
	fmt.Println("Starting server on :8080...")
	err := server.ListenAndServe()
	if err != nil {
		fmt.Printf("Error starting server: %s\n", err)
	}
}

func CORS(next http.Handler) http.HandlerFunc {
  return func(w http.ResponseWriter, r *http.Request) {
    w.Header().Add("Access-Control-Allow-Origin", "*")
    w.Header().Add("Access-Control-Allow-Credentials", "true")
    w.Header().Add("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
    w.Header().Add("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")

    if r.Method == "OPTIONS" {
        http.Error(w, "No Content", http.StatusNoContent)
        return
    }

    next.ServeHTTP(w, r)
  }
}