package auth

import (
	"encoding/json"
	"log"
	"net/http"
)

type User struct {
	UserID   int    `json:"userid"`
	Username string `json:"username"`
}

func Login(w http.ResponseWriter, r *http.Request) {
	log.Println("Login received request at /login")

	err := r.ParseMultipartForm(1000)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("form value: %+v\n", r.Form)
	log.Printf("form value: %+v\n", r.FormValue("username"))
	log.Printf("form value: %+v\n", r.FormValue("password"))
	// mock user for now
	if r.FormValue("username") != "testuser" || r.FormValue("password") != "123" {
		log.Println("Invalid username or password")
		http.Error(w, "Invalid username or password", http.StatusUnauthorized)
		return
	} else {

		// mock user object
		user := User{
			UserID:   0,
			Username: "testuser",
		}

		// Set response headers
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		// Encode user as JSON and send back to client
		if err := json.NewEncoder(w).Encode(user); err != nil {
			log.Println("Error encoding JSON response:", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
	}
}
