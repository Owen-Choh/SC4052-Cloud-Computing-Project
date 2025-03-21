package conversation

import (
	"database/sql"

	"github.com/Owen-Choh/SC4052-Cloud-Computing-Assignment-2/chatbot-backend/chatbot/types"
	"github.com/Owen-Choh/SC4052-Cloud-Computing-Assignment-2/chatbot-backend/utils"
)

type ConversationStore struct {
	db *sql.DB
}

func NewConversationStore(db *sql.DB) *ConversationStore {
	return &ConversationStore{db: db}
}

func (s *ConversationStore) GetConversationsByID(conversationid string) ([]types.Conversation, error) {
	rows, err := s.db.Query("SELECT * FROM conversations WHERE conversationid=? ORDER BY chatid", conversationid)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	conversations := []types.Conversation{}
	for rows.Next() {
		conversation, err := scanRowsIntoConversation(rows)
		if err != nil {
			return nil, err
		}
		conversations = append(conversations, *conversation)
	}

	return conversations, nil
}

func (s *ConversationStore) CreateConversation(conversationPayload types.CreateConversationPayload) (int, error) {
	currentTime, _ := utils.GetCurrentTime()
	// temp_filepath := "tempfilepath.pdf"

	res, dberr := s.db.Exec(
		"INSERT INTO conversations (conversationid, chatbotid, username, chatbotname, role, chat, createddate) VALUES (?, ?, ?, ?, ?, ?, ?)",
		conversationPayload.Conversationid,
		conversationPayload.Chatbotid,
		conversationPayload.Username,
		conversationPayload.Chatbotname,
		conversationPayload.Role,
		conversationPayload.Chat,
		currentTime,
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

func scanRowsIntoConversation(rows *sql.Rows) (*types.Conversation, error) {
	conversation := new(types.Conversation)

	err := rows.Scan(
		&conversation.Chatid,
		&conversation.Conversationid,
		&conversation.Chatbotid,
		&conversation.Username,
		&conversation.Chatbotname,
		&conversation.Role,
		&conversation.Chat,
		&conversation.Createddate,
	)
	if err != nil {
		return nil, err
	}
	return conversation, nil
}
