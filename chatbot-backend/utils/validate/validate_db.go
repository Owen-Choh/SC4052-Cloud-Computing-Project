package validate

import (
	"database/sql"
	"log"

	"github.com/Owen-Choh/SC4052-Cloud-Computing-Assignment-2/chatbot-backend/chatbot/db"
)

func CheckAndInitDB() (*sql.DB, error) {
	isInitialised, dberr := db.InitDB()
	if isInitialised {
		if dberr != nil {
			log.Printf("Did not reinitialise db as %s", dberr)
		} else {
			log.Println("Database initialised")
		}
	} else if dberr != nil {
		log.Println("Abort server start up due to error initalising database")
		return nil, dberr
	}

	return db.GetDBConnection()
}