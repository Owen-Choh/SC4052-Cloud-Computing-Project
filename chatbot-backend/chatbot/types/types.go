package types

import ()

type ChatbotStoreInterface interface {
	GetChatbotByName(botname string) (*Chatbot, error)
	GetChatbotByID(id int) (*Chatbot, error)
	CreateChatbot(CreateChatbotPayload) error
}

type CreateChatbotPayload struct {
	userid      int    `json:"userid"`
	chatbotname string `json:"chatbotname"`
	usercontext string `json:"usercontext"`
	file        string `json:"file"`
}
type Chatbot struct {
	chatbotid   int    `json:"chatbotid"`
	userid      int    `json:"userid"`
	chatbotname string `json:"chatbotname"`
	usercontext string `json:"usercontext"`
	createddate string `json:"createddate"`
	updateddate string `json:"updateddate"`
	lastused    string `json:"lastused"`
	filepath    string `json:"filepath"`
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
