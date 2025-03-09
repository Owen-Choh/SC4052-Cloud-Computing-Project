package utils

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/Owen-Choh/SC4052-Cloud-Computing-Assignment-2/chatbot-backend/chatbot/config"
)

func ParseJSON(r *http.Request, payload any) error {
	if r.Body == nil {
		return fmt.Errorf("missing request body")
	}

	return json.NewDecoder(r.Body).Decode(payload)
}

func WriteJSON(w http.ResponseWriter, statusCode int, payload any) error {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	return json.NewEncoder(w).Encode(payload)
}

func WriteError(w http.ResponseWriter, statusCode int, err error) {
	WriteJSON(w, statusCode, map[string]string{"error": err.Error()})
}

// returns time object in specified timezone. Errors and local time if error
func GetTimezone() (time.Time, error) {
	loc, err := time.LoadLocation(config.Envs.Timezone)
	if err != nil {
		log.Println("Error loading timezone:", err)
		return time.Now(), err
	}

	// Get current time in Singapore time
	sgtTime := time.Now().In(loc)
	return sgtTime, nil
}

func GetCurrentTime() (string, error) {
	sgtTime, err := GetTimezone()
	if err != nil {
		return "", err
	}

	// Format the time (YYYY-MM-DD HH:MM:SS)
	formattedTime := sgtTime.Format(config.Envs.Default_Time)
	log.Println("Current Time in SGT:", formattedTime)
	return formattedTime, nil
}
