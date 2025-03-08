package main

import (
	"net/http"
	"github.com/Owen-Choh/SC4052-Cloud-Computing-Assignment-2/chatbot-backend/chatbot"
)

func SetUpMuxWithRoutes() *http.ServeMux {
	router := http.NewServeMux()

	router.HandleFunc("GET /api/", chatbot.HelloWorld)
	router.HandleFunc("POST /api/login/", chatbot.Login)

	return router
}