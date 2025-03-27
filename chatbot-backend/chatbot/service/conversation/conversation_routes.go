package conversation

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

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
	conversationStore types.ConversationStoreInterface
	apiFileStore      types.APIFileStoreInterface
	genaiCtx          context.Context
	genaiClient       *genai.Client // Shared Gemini API client
}

func NewHandler(chatbotStore types.ChatbotStoreInterface, conversationStore types.ConversationStoreInterface, apifileStore types.APIFileStoreInterface, apiKey string) (*Handler, error) {
	// Initialize the Gemini client
	// modelName := "gemini-2.0-flash-thinking-exp-01-21"
	// modelName := "gemini-2.0-pro-exp-02-05"
	ctx, client := setupAiCtxAndClient(apiKey)

	return &Handler{
		chatbotStore:      chatbotStore,
		conversationStore: conversationStore,
		apiFileStore:      apifileStore,
		genaiCtx:          ctx,
		genaiClient:       client,
	}, nil
}

func (h *Handler) RegisterRoutes(router *http.ServeMux) {
	router.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello from conversations")
	})

	router.HandleFunc("GET /start/{username}/{chatbotName}", h.StartConversation)
	router.HandleFunc("POST /chat/{username}/{chatbotName}", h.ChatWithChatbot)
	router.HandleFunc("POST /chat/test/{username}/{chatbotName}", h.ChatWithChatbotTest)

	router.HandleFunc("GET /test", func(w http.ResponseWriter, r *http.Request) {
		t, err := utils.GetCurrentTime()
		log.Println("Current time:", t, err)
		tz, tzerr := utils.GetTimezone()
		log.Println("Timezone:", tz, tzerr)
	})
}

func (h *Handler) StartConversation(w http.ResponseWriter, r *http.Request) {
	// Get username and chatbot name from request
	username := r.PathValue("username")
	chatbotName := r.PathValue("chatbotName")
	if username == "" || chatbotName == "" {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid parameters"))
		return
	}

	chatbot, err := h.chatbotStore.GetChatbotByName(username, chatbotName)
	if err != nil {
		log.Println("Error getting chatbot to start conversation:", err)
		utils.WriteError(w, http.StatusNotFound, err)
		return
	}

	if !chatbot.IsShared {
		utils.WriteError(w, http.StatusForbidden, fmt.Errorf("chatbot is not shared"))
		return
	}

	// Generate a new conversation ID to track this conversation in db
	conversationID := utils.GenerateUUID().String()
	utils.WriteJSON(w, http.StatusOK, map[string]string{
		"conversationid": conversationID,
		"description":    chatbot.Description,
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
	if !chatbot.IsShared {
		utils.WriteError(w, http.StatusForbidden, fmt.Errorf("chatbot is not shared"))
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

	// Initialize the Gemini model
	modelName := "gemini-2.0-pro-exp-02-05"
	genaiModel := h.genaiClient.GenerativeModel(modelName)

	genaiModel.SetTemperature(0.9)
	genaiModel.SetTopK(40)
	genaiModel.SetTopP(0.95)
	genaiModel.SetMaxOutputTokens(8192)
	genaiModel.ResponseMIMEType = "text/plain"

	// files provided during configuration of chatbot
	systemFileURIs := []string{}
	if chatbot.Filepath != "" {
		systemFileURIs = []string{
			h.checkAndUploadToGemini(chatbot.Filepath, chatbot.Chatbotid),
		}
	}
	genaiModel.SystemInstruction = &genai.Content{
		Parts: getSystemInstructionParts(*chatbot),
	}

	log.Printf("start chatid: %v", chatRequest.Conversationid)
	session := genaiModel.StartChat()
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
	resp, err := session.SendMessage(h.genaiCtx, genai.Text(chatRequest.Message))
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) {
			log.Printf("%s\n", apiErr.Body)
		}
		log.Println("WARNING: api call is not working")
		return
	}

	currentTime, err := utils.GetCurrentTime()
	// save to database and collate response to send back to user
	h.conversationStore.CreateConversation(types.NewConversation{
		Conversationid: conversationID,
		Chatbotid:      chatbot.Chatbotid,
		Username:       chatbot.Username,
		Chatbotname:    chatbot.Chatbotname,
		Role:           "user",
		Chat:           chatRequest.Message,
		Createddate:    currentTime,
	})
	responseString := ""
	for _, part := range resp.Candidates[0].Content.Parts {
		go func(part genai.Part) {
			chat := string(part.(genai.Text))
			_, err := h.conversationStore.CreateConversation(types.NewConversation{
				Conversationid: conversationID,
				Chatbotid:      chatbot.Chatbotid,
				Username:       chatbot.Username,
				Chatbotname:    chatbot.Chatbotname,
				Role:           "model",
				Chat:           chat,
				Createddate:    currentTime,
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
	if chatbot.Description != "" {
		parts = append(parts, genai.Text("This is a description of what you are: "+chatbot.Description))
	}
	if chatbot.Behaviour != "" {
		parts = append(parts, genai.Text("This is how you should behave: "+chatbot.Behaviour))
	}
	if chatbot.Usercontext != "" {
		parts = append(parts, genai.Text("This is some context you should remember: "+chatbot.Usercontext))
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

func setupAiCtxAndClient(apiKey string) (context.Context, *genai.Client) {
	ctx := context.Background()

	client, err := genai.NewClient(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		log.Fatalf("Error creating client: %v", err)
	}

	// model := client.GenerativeModel(modelName)
	return ctx, client
}

func (h *Handler) checkAndUploadToGemini(path string, chatbotid int) string {
	apiFile, err := h.apiFileStore.GetAPIFileByFilepath(path)
	// if file not found in db, upload and store in db
	if err != nil {
		log.Printf("Error getting file from db: %v", err)
		fileURI := uploadToGemini(h.genaiCtx, h.genaiClient, path)

		// store the uri in db to reuse next time
		go func() {
			currentTime, _ := utils.GetCurrentTime()
			apiFile := types.NewAPIFile{
				Chatbotid:   chatbotid,
				Createddate: currentTime,
				Filepath:    path,
				Fileuri:     fileURI,
			}
			_, err := h.apiFileStore.CreateAPIFile(apiFile)
			if err != nil {
				log.Printf("Error storing file to db: %v", err)
			}
		}()
		return fileURI
	}

	storedTime, timeerr := time.Parse(config.Envs.Time_layout, apiFile.Createddate)
	// if file exist but file is too old
	if time.Since(storedTime) > time.Duration(config.Envs.API_FILE_EXPIRATION_HOUR)*time.Hour {
		log.Printf("File is too old, reuploading. previous stored time %s parsed time %s error %v", apiFile.Createddate, storedTime, timeerr)
		fileURI := uploadToGemini(h.genaiCtx, h.genaiClient, path)

		// store the uri in db to reuse next time
		go func() {
			currentTime, _ := utils.GetCurrentTime()
			apiFile := types.UpdateAPIFile{
				Fileid:      apiFile.Fileid,
				Chatbotid:   chatbotid,
				Createddate: currentTime,
				Filepath:    path,
				Fileuri:     fileURI,
			}
			err := h.apiFileStore.UpdateAPIFile(apiFile)
			if err != nil {
				log.Printf("Error updating file in db: %v", err)
			}
		}()
		return fileURI
	}

	log.Printf("File is still valid, using %s created at %s", apiFile.Fileuri, apiFile.Createddate)
	return apiFile.Fileuri
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
