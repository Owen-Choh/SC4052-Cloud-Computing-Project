package chatbotservice

import (
	"database/sql"

	"github.com/Owen-Choh/SC4052-Cloud-Computing-Assignment-2/chatbot-backend/chatbot/types"
	"github.com/Owen-Choh/SC4052-Cloud-Computing-Assignment-2/chatbot-backend/utils"
	"github.com/go-playground/locales/id"
)

type ChatbotStore struct {
	db *sql.DB
}

func NewStore(db *sql.DB) *ChatbotStore {
	return &ChatbotStore{db:db}
}

func (s *ChatbotStore) CreateChatbot(userPayload types.CreateChatbotPayload) (int, error) {
	currentTime, _ := utils.GetCurrentTime()
	temp_filepath := "tempfilepath.pdf"
	
	res, dberr := s.db.Exec(
		"INSERT INTO chatbots (userid, chatbotname, usercontext, createddate, updateddate, lastused, filepath) VALUES (?, ?, ?, ?, ?, ?, ?)",
		userPayload.Userid,
		userPayload.Chatbotname,
		userPayload.Usercontext,
		currentTime,
		currentTime,
		currentTime,
		temp_filepath,
	)
	if dberr != nil {
		return 0, dberr
	}

	id,err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}