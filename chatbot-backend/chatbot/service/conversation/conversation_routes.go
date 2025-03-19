package conversation

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/Owen-Choh/SC4052-Cloud-Computing-Assignment-2/chatbot-backend/chatbot/auth"
	"github.com/Owen-Choh/SC4052-Cloud-Computing-Assignment-2/chatbot-backend/chatbot/config"
	"github.com/Owen-Choh/SC4052-Cloud-Computing-Assignment-2/chatbot-backend/chatbot/types"
	"github.com/Owen-Choh/SC4052-Cloud-Computing-Assignment-2/chatbot-backend/utils"
	"github.com/go-playground/validator/v10"

	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/googleapi"
	"google.golang.org/api/option"
)

var ErrChatbotNotFound = errors.New("chatbot not found")

type Handler struct {
	chatbotStore      types.ChatbotStoreInterface
	conversationStore *ConversationStore
	genaiClient *genai.Client // Shared Gemini API client
	genaiModel  *genai.GenerativeModel
}

func NewHandler(chatbotStore types.ChatbotStoreInterface, conversationStore *ConversationStore, apiKey string) (*Handler, error) {
	ctx := context.Background()

	// Initialize the Gemini client
	client, err := genai.NewClient(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		return nil, fmt.Errorf("error creating Gemini client: %w", err)
	}

	model := client.GenerativeModel("gemini-2.0-flash-thinking-exp-01-21")

	return &Handler{
		chatbotStore:      chatbotStore,
		conversationStore: conversationStore,
		genaiClient: client,
		genaiModel:  model,
	}, nil
}

func (h *Handler) RegisterRoutes(router *http.ServeMux) {
	router.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello from conversations")
	})

	router.HandleFunc("GET /start", h.StartConversation)
	router.HandleFunc("POST /chat/{username}/{chatbotName}", h.ChatWithChatbot)
	router.HandleFunc("POST /chat/test/{username}/{chatbotName}", h.ChatWithChatbotTest)
}

func (h *Handler) StartConversation(w http.ResponseWriter, r *http.Request) {
	conversationID := utils.GenerateUUID().String()
	utils.WriteJSON(w, http.StatusOK, map[string]string{
		"conversationid": conversationID,
	})
}

func (h *Handler) ChatWithChatbotTest(w http.ResponseWriter, r *http.Request) {
	log.Println("ChatWithChatbotTest reply with test response")
	conversation, err := h.conversationStore.GetConversationsByID("2f9328h-fonvh0-2249")
	if err != nil {
		log.Println("Error getting test conversation:", err)
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	utils.WriteJSON(w, http.StatusOK, types.ChatResponse{Response: conversation[1].Chat})
}

func (h *Handler) ChatWithChatbot(w http.ResponseWriter, r *http.Request) {
	username := r.PathValue("username")
	chatbotName := r.PathValue("chatbotName")

	if username == "" || chatbotName == "" {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid parameters"))
		return
	}

	// Get chatbot for context
	chatbot, err := h.chatbotStore.GetChatbotByName(username, chatbotName)
	if err != nil {
		if errors.Is(err, ErrChatbotNotFound) {
			utils.WriteError(w, http.StatusNotFound, err)
		} else {
			utils.WriteError(w, http.StatusInternalServerError, err)
		}
		return
	}

	// Read input message from request
	var chatRequest types.ChatRequest
	if err := utils.ParseJSON(r, &chatRequest); err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid request body"))
		return
	}
	if err := utils.Validate.Struct(chatRequest); err != nil {
		validate_error := err.(validator.ValidationErrors)
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid payload %v", validate_error))
		return
	}

	conversationID := chatRequest.Conversationid
	conversations, err := h.conversationStore.GetConversationsByID(conversationID)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	// Initialize Gemini API client
	log.Printf("Initializing Gemini API client for %s\n", conversationID)
	apiKey := config.Envs.GEMINI_API_KEY
	if apiKey == "" {
		log.Fatalln("Environment variable GEMINI_API_KEY not set")
		utils.WriteError(w, http.StatusInternalServerError, fmt.Errorf("internal server error"))
		return
	}
	modelName := "gemini-2.0-flash-thinking-exp-01-21"
	ctx, client, model := setupAiModel(apiKey, modelName)
	defer client.Close()

	model.SetTemperature(0.9)
	model.SetTopK(40)
	model.SetTopP(0.95)
	model.SetMaxOutputTokens(8192)
	model.ResponseMIMEType = "text/plain"

	systemFileURIs := []string{}
	if chatbot.Filepath != "" {
		systemFileURIs = []string{uploadToGemini(ctx, client, chatbot.Filepath)}
	}
	model.SystemInstruction = &genai.Content{
		Parts: getSystemInstructionParts(*chatbot),
	}

	log.Printf("start chatid: %v", chatRequest.Conversationid)
	session := model.StartChat()
	// append the file to history as system instruction only allow text
	session.History = []*genai.Content{
		{
			Role: "user",
			Parts: []genai.Part{
				genai.Text("Here are some files you can use:"),
				genai.FileData{URI: systemFileURIs[0]},
			},
		},
	}
	// append the actual conversation from db
	conversationHistory := getContentFromConversions(conversations)
	session.History = append(session.History, conversationHistory...)

	log.Printf("sending msg for conversationid: %s\n", conversationID)
	resp, err := session.SendMessage(ctx, genai.Text(chatRequest.Message))
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) {
			log.Printf("%s\n", apiErr.Body)
		}
		log.Println("WARNING: api call is not working")
		return
	}

	// save to database and collate response to send back to user
	h.conversationStore.CreateConversation(types.CreateConversationPayload{
		Conversationid: conversationID,
		Chatbotid:      chatbot.Chatbotid,
		Username:       chatbot.Username,
		Chatbotname:    chatbot.Chatbotname,
		Role:           "user",
		Chat:           chatRequest.Message,
	})
	responseString := ""
	for _, part := range resp.Candidates[0].Content.Parts {
		go func(part genai.Part) {
			chat := string(part.(genai.Text))
			_, err := h.conversationStore.CreateConversation(types.CreateConversationPayload{
				Conversationid: conversationID,
				Chatbotid:      chatbot.Chatbotid,
				Username:       chatbot.Username,
				Chatbotname:    chatbot.Chatbotname,
				Role:           "model",
				Chat:           chat,
			})

			if err != nil {
				log.Printf("Error saving conversation: %v", err)
			}
		}(part)

		responseString += string(part.(genai.Text))
	}

	log.Printf("responding to conversation: %s\n", conversationID)
	utils.WriteJSON(w, http.StatusOK, types.ChatResponse{Response: responseString})
}

func getSystemInstructionParts(chatbot types.Chatbot) []genai.Part {
	parts := []genai.Part{} // Initialize empty slice
	if chatbot.Behaviour != "" {
		parts = append(parts, genai.Text("This is how you should behave: "+chatbot.Behaviour))
	}
	if chatbot.Usercontext != "" {
		parts = append(parts, genai.Text("This is the context you should remember: "+chatbot.Usercontext))
	}
	return parts
}

func getContentFromConversions(conversations []types.Conversation) []*genai.Content {
	content := []*genai.Content{}
	for _, conversation := range conversations {
		content = append(content, &genai.Content{
			Role: conversation.Role,
			Parts: []genai.Part{
				genai.Text(conversation.Chat),
			},
		})
	}
	return content
}

func setupAiModel(apiKey string, modelName string) (context.Context, *genai.Client, *genai.GenerativeModel) {
	ctx := context.Background()

	client, err := genai.NewClient(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		log.Fatalf("Error creating client: %v", err)
	}

	model := client.GenerativeModel(modelName)
	return ctx, client, model
}

func uploadToGemini(ctx context.Context, client *genai.Client, path string) string {
	file, err := os.Open(path)
	if err != nil {
		log.Fatalf("Error opening file: %v", err)
	}
	defer file.Close()

	fileData, err := client.UploadFile(ctx, "", file, nil)
	if err != nil {
		log.Fatalf("Error uploading file: %v", err)
	}

	log.Printf("Uploaded file %s as: %s", fileData.DisplayName, fileData.URI)
	return fileData.URI
}
