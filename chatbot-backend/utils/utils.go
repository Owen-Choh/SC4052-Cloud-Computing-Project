package utils

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
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

func GetTimezone() (time.Time, error) {
	loc, err := time.LoadLocation("Asia/Singapore")
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
	formattedTime := sgtTime.Format("2006-01-02 15:04:05")
	log.Println("Current Time in SGT:", formattedTime)
	return formattedTime, nil
}
