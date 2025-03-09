package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/Owen-Choh/SC4052-Cloud-Computing-Assignment-2/chatbot-backend/utils/validate"
	"github.com/Owen-Choh/SC4052-Cloud-Computing-Assignment-2/chatbot-backend/utils/middleware"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	dberr := validate.CheckAndInitDB()
	if dberr != nil {
		return
	}
	
	mainRouter := http.NewServeMux()
	mainStack := middleware.CreateStack(
		middleware.Logging, 
		middleware.CORS,
	)
	nonAuthRouter := SetUpNonAuthRouter()
  
	mainRouter.Handle("/api/", http.StripPrefix("/api", mainStack(nonAuthRouter)))
	// set server and start
	server := http.Server{
		Addr:    ":8080",
		Handler: mainRouter,
	}
	fmt.Println("Starting server on :8080...")
	err := server.ListenAndServe()
	if err != nil {
		fmt.Printf("Error starting server: %s\n", err)
	}
}
