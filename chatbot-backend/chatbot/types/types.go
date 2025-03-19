package types

import ()

type ChatbotStoreInterface interface {
	GetChatbotsByUsername(username string) ([]Chatbot, error)
	GetChatbotByName(username string, chatbotName string) (*Chatbot, error)
	CreateChatbot(userPayload NewChatbot) (int, error)
}

type NewChatbot struct {
	Username    string `json:"Username" validate:"required"`
	Chatbotname string `json:"chatbotname" validate:"required,min=3"`
	Behaviour   string `json:"behaviour"`
	Usercontext string `json:"usercontext"`
	IsShared    bool   `json:"isShared"`
	File        string `json:"file"`
}
type CreateChatbotPayload struct {
	Chatbotname string `json:"chatbotname" validate:"required,min=3"`
	Behaviour   string `json:"behaviour"`
	Usercontext string `json:"usercontext"`
	IsShared    bool   `json:"isShared" validate:"required"`
	File        string `json:"file"`
}

type Chatbot struct {
	Chatbotid   int    `json:"chatbotid"`
	Username      string    `json:"username"`
	Chatbotname string `json:"chatbotname"`
	Behaviour   string `json:"behaviour"`
	Usercontext string `json:"usercontext"`
	Createddate string `json:"createddate"`
	Updateddate string `json:"updateddate"`
	Lastused    string `json:"lastused"`
	IsShared    bool   `json:"isShared"`
	Filepath    string `json:"filepath"`
}

type UserStoreInterface interface {
	GetUserByName(username string) (*User, error)
	GetUserByID(id int) (*User, error)
	CreateUser(RegisterUserPayload) error
}

type User struct {
	Userid      int    `json:"userid"`
	Username    string `json:"username"`
	Password    string `json:"password"`
	Createddate string `json:"createdDate"`
	Lastlogin   string `json:"lastLogin"`
}

type NewUser struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type RegisterUserPayload struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required,min=3"`
}

type LoginUserPayload struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// ChatRequest defines the request body for chatting with a chatbot.
type ChatRequest struct {
	Conversationid string `json:"conversationid" validate:"required"`
	Message string `json:"message" validate:"required"`
}

// ChatResponse defines the response body for chatting with a chatbot.
type ChatResponse struct {
	Response string `json:"response"`
}

type Conversation struct {
	Chatid         int    `json:"chatid"`
	Conversationid string `json:"conversationid"`
	Chatbotid      int    `json:"chatbotid"`
	Username       string `json:"username"`
	Chatbotname    string `json:"chatbotname"`
	Role           string `json:"role"`
	Chat           string `json:"chat"`
	Createddate    string `json:"createddate"`	
}

type CreateConversationPayload struct {
	Conversationid string `json:"conversationid"`
	Chatbotid      int    `json:"chatbotid"`
	Username       string `json:"username"`
	Chatbotname    string `json:"chatbotname"`
	Role           string `json:"role"`
	Chat           string `json:"chat"`
	Createddate    string `json:"createddate"`
}

type APIfile struct {
	Fileid int `json:"fileid"`
	Chatbotid int `json:"chatbotid"`
	Createddate string `json:"createddate"`
	Filepath string `json:"filepath"`
	Fileuri string `json:"fileuri"`
}