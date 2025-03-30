package chatbotservice

import (
	"database/sql"

	"github.com/Owen-Choh/SC4052-Cloud-Computing-Assignment-2/chatbot-backend/chatbot/types"
	"github.com/Owen-Choh/SC4052-Cloud-Computing-Assignment-2/chatbot-backend/utils"
)

type ChatbotStore struct {
	db *sql.DB
}

func NewStore(db *sql.DB) types.ChatbotStoreInterface {
	return &ChatbotStore{db: db}
}

func (s *ChatbotStore) GetChatbotsByID(chatbotID int) (*types.Chatbot, error) {
	rows, err := s.db.Query("SELECT * FROM chatbots WHERE chatbotid=?", chatbotID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	chatbots := new(types.Chatbot)
	rows.Next()
	chatbots, err = scanRowsIntoChatbot(rows)

	if err != nil {
		return nil, err
	}

	return chatbots, nil
}

func (s *ChatbotStore) GetChatbotsByUsername(username string) ([]types.Chatbot, error) {
	rows, err := s.db.Query("SELECT * FROM chatbots WHERE username=?", username)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	chatbots := []types.Chatbot{}
	for rows.Next() {
		bot, err := scanRowsIntoChatbot(rows)
		if err != nil {
			return nil, err
		}
		chatbots = append(chatbots, *bot)
	}

	return chatbots, nil
}

func (s *ChatbotStore) GetChatbotByName(username string, chatbotName string) (*types.Chatbot, error) {
	rows, err := s.db.Query("SELECT * FROM chatbots WHERE username=? AND chatbotName=?", username, chatbotName)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

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

func (s *ChatbotStore) CreateChatbot(userPayload types.NewChatbot) (int, error) {
	currentTime, _ := utils.GetCurrentTime()
	// temp_filepath := "tempfilepath.pdf"

	res, dberr := s.db.Exec(
		"INSERT INTO chatbots (username, chatbotname, description, behaviour, usercontext, createddate, updateddate, lastused, isShared, filepath) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)",
		userPayload.Username,
		userPayload.Chatbotname,
		userPayload.Description,
		userPayload.Behaviour,
		userPayload.Usercontext,
		currentTime,
		currentTime,
		currentTime,
		userPayload.IsShared,
		userPayload.File,
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

func (s *ChatbotStore) UpdateChatbot(chatbotPayload types.UpdateChatbot) error {
	currentTime, _ := utils.GetCurrentTime()

	_, err := s.db.Exec(
		"UPDATE chatbots SET chatbotname=?, description=?, behaviour=?, usercontext=?, updateddate=?, isShared=?, filepath=? WHERE chatbotid=? AND username=?",
		chatbotPayload.Chatbotname,
		chatbotPayload.Description,
		chatbotPayload.Behaviour,
		chatbotPayload.Usercontext,
		currentTime,
		chatbotPayload.IsShared,
		chatbotPayload.File,
		chatbotPayload.Chatbotid,
		chatbotPayload.Username,
	)
	return err
}

func (s *ChatbotStore) DeleteChatbot(chatbotID int) error {
	_, err := s.db.Exec("DELETE FROM chatbots WHERE chatbotid=?", chatbotID)
	return err
}

func scanRowsIntoChatbot(rows *sql.Rows) (*types.Chatbot, error) {
	chatbot := new(types.Chatbot)

	err := rows.Scan(
		&chatbot.Chatbotid,
		&chatbot.Username,
		&chatbot.Chatbotname,
		&chatbot.Description,
		&chatbot.Behaviour,
		&chatbot.Usercontext,
		&chatbot.Createddate,
		&chatbot.Updateddate,
		&chatbot.Lastused,
		&chatbot.IsShared,
		&chatbot.Filepath,
	)
	if err != nil {
		return nil, err
	}
	return chatbot, nil
}
