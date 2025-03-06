package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/Owen-Choh/SC4052-Cloud-Computing-Assignment-2/chatbot-backend/chatbot"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	isInitialised, dberr := chatbot.InitDB()
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
	
	router := SetUpRoutes()
	// set server and start
	server := http.Server{
		Addr:    ":8080",
		Handler: router,
	}
	fmt.Println("Starting server on :8080...")
	err := server.ListenAndServe()
	if err != nil {
		fmt.Printf("Error starting server: %s\n", err)
	}
}
