package chatbot

import (
	"fmt"
	"net/http"
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
)

func InitDB(w http.ResponseWriter, r *http.Request) {
	// Open a connection to the database
	db, err := sql.Open("sqlite3", "./database/chatbot.db")
	if err != nil {
		fmt.Println(err)
	}
	defer db.Close()

	// Create tables
	_, err = db.Exec(`
	CREATE TABLE IF NOT EXISTS users (
		userid INTEGER PRIMARY KEY AUTOINCREMENT,
		username TEXT UNIQUE NOT NULL,
		password TEXT NOT NULL
	);`)
	if err != nil {
		fmt.Printf("Error initalising users table: %s\n", err)
	}

	_, err = db.Exec(`
	CREATE TABLE IF NOT EXISTS chatbots (
		chatbotid INTEGER PRIMARY KEY AUTOINCREMENT,
		userid INTEGER NOT NULL,
		chatbotname TEXT NOT NULL,
		usercontext TEXT NOT NULL,
		filepath TEXT NOT NULL,
		FOREIGN KEY(userid) REFERENCES users(userid),
		UNIQUE(userid, chatbotname)
	);`)
	if err != nil {
		fmt.Printf("Error initalising chatbots table: %s\n", err)

		fmt.Fprint(w, "Error occured while initializing database")
		return
	}

	fmt.Fprint(w, "Database initialized")
}