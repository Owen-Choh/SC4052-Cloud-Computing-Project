package main

import (
	"net/http"

	"github.com/Owen-Choh/SC4052-Cloud-Computing-Assignment-2/chatbot-backend/chatbot/auth"
	"github.com/Owen-Choh/SC4052-Cloud-Computing-Assignment-2/chatbot-backend/chatbot/chatbotservice"
)

func SetUpNonAuthRouter() *http.ServeMux {
	router := http.NewServeMux()

	router.HandleFunc("POST /login", auth.Login)

	return router
}

func SetUpAuthRouter() *http.ServeMux {
	router := http.NewServeMux()

	router.HandleFunc("GET /chatbot", chatbotservice.HelloBot)

	return router
}