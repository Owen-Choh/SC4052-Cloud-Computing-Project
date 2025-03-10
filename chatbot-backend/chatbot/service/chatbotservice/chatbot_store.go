package chatbotservice

import "database/sql"

type ChatbotStore struct {
	db *sql.DB
}