package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/Owen-Choh/SC4052-Cloud-Computing-Assignment-2/chatbot-backend/chatbot/config"
	"github.com/Owen-Choh/SC4052-Cloud-Computing-Assignment-2/chatbot-backend/chatbot/service/chatbotservice"
	"github.com/Owen-Choh/SC4052-Cloud-Computing-Assignment-2/chatbot-backend/chatbot/service/user"
	"github.com/Owen-Choh/SC4052-Cloud-Computing-Assignment-2/chatbot-backend/utils/middleware"
	"github.com/Owen-Choh/SC4052-Cloud-Computing-Assignment-2/chatbot-backend/utils/validate"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	

	dbConnection, dberr := validate.CheckAndInitDB()
	if dberr != nil {
		return
	}

	mainRouter := http.NewServeMux()
	mainStack := middleware.CreateStack(
		middleware.Logging,
		middleware.CORS,
	)
	mainRouter.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// write to the response which returns to client
		fmt.Fprintf(w, "Hello world!")
	})

	userSubRouter := http.NewServeMux()
	userStore := user.NewStore(dbConnection)
	userHandler := user.NewHandler(userStore)
	userHandler.RegisterRoutes(userSubRouter)
	
	mainRouter.Handle("/api/user/", http.StripPrefix("/api/user", mainStack(userSubRouter)))

	chatbotSubRouter := http.NewServeMux()
	chatbotStore := chatbotservice.NewStore(dbConnection)
	conversationStore := chatbotservice.NewConversationStore(dbConnection)
	chatbotHandler := chatbotservice.NewHandler(chatbotStore, userStore, conversationStore)
	chatbotHandler.RegisterRoutes(chatbotSubRouter)
	
	mainRouter.Handle("/api/chatbot/", http.StripPrefix("/api/chatbot", mainStack(chatbotSubRouter)))

	
	// set server and start
	server := http.Server{
		Addr:    ":" + config.Envs.Port,
		Handler: mainRouter,
	}
	log.Printf("Starting server on :%s...\n", config.Envs.Port)
	err := server.ListenAndServe()
	if err != nil {
		log.Printf("Error starting server: %s\n", err)
	}
}
