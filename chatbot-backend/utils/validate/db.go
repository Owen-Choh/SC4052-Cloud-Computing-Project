package validate

import (
	"log"
	"github.com/Owen-Choh/SC4052-Cloud-Computing-Assignment-2/chatbot-backend/chatbot/db"
)

func CheckAndInitDB() error {
	isInitialised, dberr := db.InitDB()
	if isInitialised {
		if dberr != nil {
			log.Printf("Did not reinitialise db as %s", dberr)
		} else {
			log.Println("Database initialised")
		}
	} else if dberr != nil {
		log.Println("Abort server start up due to error initalising database")
		return dberr
	}
	return nil
}