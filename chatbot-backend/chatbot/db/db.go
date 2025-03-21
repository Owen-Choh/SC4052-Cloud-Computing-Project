package db

import (
	"database/sql"
	"errors"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

func GetDBConnection() (*sql.DB, error){
	return sql.Open("sqlite3", "./database_files/chatbot.db")
}

func InitDB() (bool, error) {
	// Open a connection to the database
	db, err := GetDBConnection()
	if err != nil {
		log.Println(err)
		return false, err
	}
	defer db.Close()

	if checktablesexist(db) {
		return true, errors.New("all tables already exist")
	}

	// Create tables
	_, err = db.Exec(`
	CREATE TABLE IF NOT EXISTS users (
		userid INTEGER PRIMARY KEY AUTOINCREMENT,
		username TEXT UNIQUE NOT NULL,
		password TEXT NOT NULL,
		createddate TEXT NOT NULL,
		lastlogin TEXT NOT NULL
	);`)
	if err != nil {
		log.Printf("Error initalising users table: %s\n", err)
		return false, err
	}

	_, err = db.Exec(`
	CREATE TABLE IF NOT EXISTS chatbots (
		chatbotid INTEGER PRIMARY KEY AUTOINCREMENT,
		username TEXT NOT NULL,
		chatbotname TEXT NOT NULL,
		description TEXT NOT NULL Default '',
		behaviour TEXT NOT NULL,
		usercontext TEXT NOT NULL,
		createddate TEXT NOT NULL,
		updateddate TEXT NOT NULL,
		lastused TEXT NOT NULL,
		isShared BOOLEAN NOT NULL DEFAULT FALSE,
		filepath TEXT NOT NULL,
		FOREIGN KEY(username) REFERENCES users(username),
		UNIQUE(username, chatbotname)
	);`)
	if err != nil {
		log.Printf("Error initalising chatbots table: %s\n", err)
		return false, err
	}

	_, err = db.Exec(`
	CREATE TABLE IF NOT EXISTS conversations (
		chatid INTEGER PRIMARY KEY AUTOINCREMENT,
		conversationid TEXT NOT NULL,
		chatbotid INTEGER NOT NULL,
		username TEXT NOT NULL,
		chatbotname TEXT NOT NULL,
		role TEXT NOT NULL,
		chat TEXT NOT NULL,
		createddate TEXT NOT NULL,
		FOREIGN KEY(chatbotid) REFERENCES chatbots(chatbotid),
		FOREIGN KEY(username) REFERENCES users(username),
		FOREIGN KEY(chatbotname) REFERENCES chatbots(chatbotname)
	);`)
	if err != nil {
		log.Printf("Error initalising conversations table: %s\n", err)
		return false, err
	}

	_, err = db.Exec(`
	CREATE TABLE IF NOT EXISTS apifiles (
		fileid INTEGER PRIMARY KEY AUTOINCREMENT,
		chatbotid INTEGER NOT NULL,
		createddate TEXT NOT NULL,
		filepath TEXT NOT NULL,
		fileuri TEXT NOT NULL,
		FOREIGN KEY(chatbotid) REFERENCES chatbots(chatbotid)
	);`)
	if err != nil {
		log.Printf("Error initalising conversations table: %s\n", err)
		return false, err
	}
	
	return true, err
}

func checktablesexist(db *sql.DB) bool {
	var answer bool = true

	rows, err := db.Query("SELECT name FROM sqlite_master WHERE type='table' AND name IN ('users', 'chatbots', 'conversations', 'apifiles');")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	// Read results into a map
	tableExists := map[string]bool{
		"users":    false,
		"chatbots": false,
		"conversations": false,
		"apifiles": false,
	}

	for rows.Next() {
		var tableName string
		if err := rows.Scan(&tableName); err != nil {
			log.Fatal(err)
		}
		tableExists[tableName] = true
	}

	// Print results
	for table, exists := range tableExists {
		if !exists {
			log.Printf("Table '%s' does not exist.\n", table)
			answer = false
		}
	}

	return answer
}
