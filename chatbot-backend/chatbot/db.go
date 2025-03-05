package chatbot

import (
	"fmt"
	"net/http"
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
)

func CheckDB(w http.ResponseWriter, r *http.Request) {
	// Open a connection to the database
	db, err := sql.Open("sqlite3", "./database/chatbot.db")
	if err != nil {
		fmt.Println(err)
		fmt.Fprint(w,err.Error())
	}
	defer db.Close()

	// Create a table (if not exists)
	_, err = db.Exec(`
	CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		username TEXT UNIQUE NOT NULL,
		password TEXT NOT NULL
	);`)
	if err != nil {
		fmt.Println(err)
		fmt.Fprint(w,err.Error())
	}

	fmt.Fprintf(w,"Database connected and table ensured.")
}