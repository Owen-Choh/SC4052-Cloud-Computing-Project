package main

import (
	"net/http"
	"github.com/Owen-Choh/SC4052-Cloud-Computing-Assignment-2/chatbot-backend/chatbot/user"
	"github.com/Owen-Choh/SC4052-Cloud-Computing-Assignment-2/chatbot-backend/chatbot/auth"
)

func SetUpNonAuthRouter() *http.ServeMux {
	router := http.NewServeMux()

	router.HandleFunc("GET /", user.HelloWorld)
	router.HandleFunc("POST /login", auth.Login)

	return router
}