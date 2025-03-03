package chatbot

import (
	"fmt"
	"net/http"
)

func HelloWorld(w http.ResponseWriter, r *http.Request) {
	fmt.Println("HelloWorld received request at /api/")
	// write to the response which returns to client
	fmt.Fprintf(w, "Hello world!")
}