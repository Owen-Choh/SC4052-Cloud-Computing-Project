package chatbotservice

import (
	"fmt"
	"net/http"
)

func HelloBot(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, I am a chatbot.")
}