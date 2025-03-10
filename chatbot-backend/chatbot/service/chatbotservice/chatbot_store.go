package chatbotservice

import (
	"database/sql"

	"github.com/Owen-Choh/SC4052-Cloud-Computing-Assignment-2/chatbot-backend/chatbot/types"
	"github.com/Owen-Choh/SC4052-Cloud-Computing-Assignment-2/chatbot-backend/utils"
)

type ChatbotStore struct {
	db *sql.DB
}

func NewStore(db *sql.DB) *ChatbotStore {
	return &ChatbotStore{db: db}
}

func (s *ChatbotStore) GetChatbotByName(username string, chatbotName string) (*types.Chatbot, error) {
	rows, err := s.db.Query("SELECT * FROM chatbots WHERE username=? AND chatbotName=?", username, chatbotName)
	if err != nil {
		return nil, err
	}

	chatbot := new(types.Chatbot)
	for rows.Next() {
		chatbot, err = scanRowsIntoChatbot(rows)
		if err != nil {
			return nil, err
		}
	}

	if chatbot.Chatbotid == 0 {
		return nil, ErrChatbotNotFound
	}

	return chatbot, nil
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

	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

func scanRowsIntoChatbot(rows *sql.Rows) (*types.Chatbot, error) {
	chatbot := new(types.Chatbot)

	err := rows.Scan(
		&chatbot.Chatbotid,
		&chatbot.Userid,
		&chatbot.Chatbotname,
		&chatbot.Usercontext,
		&chatbot.Createddate,
		&chatbot.Updateddate,
		&chatbot.Lastused,
		&chatbot.Filepath,
	)
	if err != nil {
		return nil, err
	}
	return chatbot, nil
}
