package types

// UserStoreInterface defines the methods for user store
type UserStoreInterface interface {
	GetUserByID(id int) (*User, error)
	GetUserByName(username string) (*User, error)
	CreateUser(user RegisterUserPayload) error
}

// ChatbotStoreInterface defines the methods for chatbot store
type ChatbotStoreInterface interface {
	GetChatbotsByID(chatbotID int) (*Chatbot, error)
	GetChatbotsByUsername(username string) ([]Chatbot, error)
	GetChatbotByName(username string, chatbotName string) (*Chatbot, error)
	CreateChatbot(userPayload NewChatbot) (int, error)
	UpdateChatbot(chatbotPayload UpdateChatbot) error
	DeleteChatbot(chatbotID int) error
}

// ConversationStoreInterface defines the methods for conversation store
type ConversationStoreInterface interface {
	GetConversationsByID(conversationID string) ([]Conversation, error)
	GetConversationsByUserID(userID int) ([]Conversation, error)
	CreateConversation(conversationPayload NewConversation) (int, error)
	UpdateConversation(conversationPayload UpdateConversation) error
	DeleteConversation(conversationID int) error
}

// APIFileStoreInterface defines the methods for API file store
type APIFileStoreInterface interface {
	GetAPIFileByID(apiFileID int) (*APIFile, error)
	GetAPIFilesByUserID(userID int) ([]APIFile, error)
	GetAPIFileByFilepath(filepath string) (*APIFile, error)
	CreateAPIFile(apiFilePayload NewAPIFile) (int, error)
	UpdateAPIFile(apiFilePayload UpdateAPIFile) error
	DeleteAPIFile(apiFileID int) error
}
