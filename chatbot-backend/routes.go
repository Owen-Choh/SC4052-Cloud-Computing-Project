package main

import (
	"net/http"
	"github.com/Owen-Choh/SC4052-Cloud-Computing-Assignment-2/chatbot-backend/chatbot"
)

func SetUpRoutes() *http.ServeMux {
	router := http.NewServeMux()

	router.HandleFunc("GET /api/", chatbot.HelloWorld)

	return router
}