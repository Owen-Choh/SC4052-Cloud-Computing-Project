package types

import ()

type ChatbotStoreInterface interface {
	GetChatbotByName(botname string) (*Chatbot, error)
	GetChatbotByID(id int) (*Chatbot, error)
	CreateChatbot(CreateChatbotPayload) error
}

type CreateChatbotPayload struct {
	Userid      int    `json:"userid"`
	Chatbotname string `json:"chatbotname"`
	Usercontext string `json:"usercontext"`
	File        string `json:"file"`
}
type Chatbot struct {
	Chatbotid   int    `json:"chatbotid"`
	Userid      int    `json:"userid"`
	Chatbotname string `json:"chatbotname"`
	Usercontext string `json:"usercontext"`
	Createddate string `json:"createddate"`
	Updateddate string `json:"updateddate"`
	Lastused    string `json:"lastused"`
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
