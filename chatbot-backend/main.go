package main

import (
	"fmt"
	"net/http"
)

func main() {
	router := SetUpRoutes()
	// set server and start
	server:= http.Server{
		Addr: ":8080",
		Handler: router,
	}
	fmt.Println("Starting server on :8080...")
	err := server.ListenAndServe()
	if err != nil {
		fmt.Printf("Error starting server: %s\n", err)
	}
}