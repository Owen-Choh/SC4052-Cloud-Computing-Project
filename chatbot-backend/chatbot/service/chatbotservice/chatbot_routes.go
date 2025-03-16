package chatbotservice

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
	store             types.ChatbotStoreInterface
	userstore         types.UserStoreInterface
	conversationStore *ConversationStore
	// genaiClient *genai.Client // Gemini API client
	// genaiModel  *genai.GenerativeModel
}

func NewHandler(store types.ChatbotStoreInterface, userstore types.UserStoreInterface, conversationStore *ConversationStore) *Handler {
	// ctx := context.Background()

	// Initialize the Gemini client
	// client, err := genai.NewClient(ctx, option.WithAPIKey(apiKey))
	// if err != nil {
	// 	return nil, fmt.Errorf("error creating Gemini client: %w", err)
	// }

	// model := client.GenerativeModel("gemini-2.0-flash-thinking-exp-01-21")

	return &Handler{
		store:             store,
		userstore:         userstore,
		conversationStore: conversationStore,
		// genaiClient: client,
		// genaiModel:  model,
	}
}

func (h *Handler) RegisterRoutes(router *http.ServeMux) {
	router.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello from chatbot")
	})
	router.HandleFunc("GET /list", auth.WithJWTAuth(h.GetUserChatbot, h.userstore))
	router.HandleFunc("GET /{username}/{chatbotName}", h.GetChatbot)
	router.HandleFunc("POST /newchatbot", auth.WithJWTAuth(h.CreateChatbot, h.userstore))

	router.HandleFunc("POST /chat/{username}/{chatbotName}", h.ChatWithChatbot)
}

// ... (GetUserChatbot and GetChatbot remain the same)

func (h *Handler) GetUserChatbot(w http.ResponseWriter, r *http.Request) {
	username := auth.GetUsernameFromContext(r.Context())
	if username == "" {
		log.Println("username missing in request context")
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid request"))
		return
	}

	log.Println("authenticated user: " + username)
	chatbots, err := h.store.GetChatbotsByUsername(username)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, chatbots)
}

func (h *Handler) GetChatbot(w http.ResponseWriter, r *http.Request) {
	username := r.PathValue("username")
	chatbotName := r.PathValue("chatbotName")

	if username == "" || chatbotName == "" {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid parameters"))
		return
	}

	chatbot, err := h.store.GetChatbotByName(username, chatbotName)
	if err != nil {
		if errors.Is(err, ErrChatbotNotFound) {
			utils.WriteError(w, http.StatusNotFound, err)
		} else {
			utils.WriteError(w, http.StatusInternalServerError, err)
		}
		return
	}

	utils.WriteJSON(w, http.StatusOK, chatbot)
}

func (h *Handler) CreateChatbot(w http.ResponseWriter, r *http.Request) {
	username := auth.GetUsernameFromContext(r.Context())
	if username == "" {
		log.Println("username missing in request context")
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid request"))
		return
	}

	err := r.ParseMultipartForm(10 << 20) // 10MB limit
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("failed to parse form"))
		return
	}

	// Extract chatbot fields from form
	chatbotname := r.FormValue("chatbotname")
	behaviour := r.FormValue("behaviour")
	usercontext := r.FormValue("usercontext")
	isShared := r.FormValue("isShared") == "true"

	// Handle file upload
	file, header, err := r.FormFile("file")
	var filepath string
	if err == nil {
		defer file.Close()

		fullDirPath := "database_files/uploads/" + username + "/" + chatbotname
		err := os.MkdirAll(fullDirPath, os.ModePerm) // Create the directory if it doesn’t exist
		if err != nil {
			log.Println("Error creating directory:", err)
			utils.WriteError(w, http.StatusInternalServerError, fmt.Errorf("failed to save file"))
			return
		}

		// Save the uploaded file
		filepath = fullDirPath + "/" + header.Filename
		log.Println("Saving file to:", filepath)
		out, err := os.Create(filepath)
		if err != nil {
			log.Println("Error saving file:", err)
			utils.WriteError(w, http.StatusInternalServerError, fmt.Errorf("failed to save file"))
			return
		}
		defer out.Close()
		io.Copy(out, file)
	} else {
		filepath = "" // No file uploaded
	}

	// Create chatbot struct
	newChatbot := types.NewChatbot{
		Username:    username,
		Chatbotname: chatbotname,
		Behaviour:   behaviour,
		IsShared:    isShared,
		Usercontext: usercontext,
		File:        filepath,
	}
	if err := utils.Validate.Struct(newChatbot); err != nil {
		validate_error := err.(validator.ValidationErrors)
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid payload %v", validate_error))
		return
	}

	botID, err := h.store.CreateChatbot(newChatbot)
	if err != nil {
		log.Println("Error creating chatbot:", err)
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusCreated, map[string]interface{}{
		"chatbotid": botID,
	})
}

func (h *Handler) ChatWithChatbot(w http.ResponseWriter, r *http.Request) {
	username := r.PathValue("username")
	chatbotName := r.PathValue("chatbotName")

	if username == "" || chatbotName == "" {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid parameters"))
		return
	}

	// Get chatbot for context
	chatbot, err := h.store.GetChatbotByName(username, chatbotName)
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
	if conversationID == "" {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid parameters"))
		return
	}

	// conversations, err := h.conversationStore.GetConversationsByID(conversationID)
	// if err != nil {
	// 	utils.WriteError(w, http.StatusInternalServerError, err)
	// 	return
	// }

	// Initialize Gemini API client
	log.Println("Initializing Gemini API client")
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

	log.Println("add system file")
	systemFileURIs := []string{}
	if chatbot.Filepath != "" {
		systemFileURIs = []string{uploadToGemini(ctx, client, chatbot.Filepath)}
	}
	log.Println("add system instruction")
	model.SystemInstruction = &genai.Content{
		Parts: getSystemInstructionParts(*chatbot),
	}

	log.Println("start chat")
	session := model.StartChat()
	//session.History = getContentFromConversions(conversations)
	session.History = []*genai.Content{
		{
			Role: "user",
			Parts: []genai.Part{
				genai.Text("Here are some files you can use:"),
				genai.FileData{URI: systemFileURIs[0]},
			},
		},
	}

	log.Println("send msg")
	resp, err := session.SendMessage(ctx, genai.Text(chatRequest.Message))
	if err != nil {
		// log.Fatalf("Error sending message: %v", err)
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) {
			log.Fatalf("%s", apiErr.Body)
		}
		log.Fatalln("shutting down server as api not working")
		return
	}
	log.Println("Got response")
	log.Printf("Response: %v\n", resp)
	// save to database and colate response to send back to user
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

	log.Printf("opened file %s", path)

	log.Printf("Uploading file %s", path)

	fileData, err := client.UploadFile(ctx, "", file, nil)
	// fileData, err := client.UploadFileFromPath(ctx, path, nil)
	if err != nil {
		log.Fatalf("Error uploading file: %v", err)
	}

	log.Printf("Uploaded file %s as: %s", fileData.DisplayName, fileData.URI)
	return fileData.URI
}

/*
func (h *Handler) ChatWithChatbot(w http.ResponseWriter, r *http.Request) {
	username := r.PathValue("username")
	chatbotName := r.PathValue("chatbotName")

	if username == "" || chatbotName == "" {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid parameters"))
		return
	}

	// Get chatbot for context
	chatbot, err := h.store.GetChatbotByName(username, chatbotName)
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

	// Construct prompt, including chatbot behavior, user context, and file content
	prompt := h.buildPrompt(chatbot, chatRequest.Message)

	// Send the prompt to Gemini and get the response
	response, err := h.sendPromptToGemini(r.Context(), prompt)
	if err != nil {
		log.Printf("Error from Gemini API: %v", err)
		utils.WriteError(w, http.StatusInternalServerError, fmt.Errorf("error getting response from AI model"))
		return
	}

	// Send the response back to the client
	utils.WriteJSON(w, http.StatusOK, ChatResponse{Response: response})
}

// buildPrompt constructs the full prompt to send to Gemini.
func (h *Handler) buildPrompt(chatbot types.Chatbot, userMessage string) string {
	promptBuilder := strings.Builder{}

	// Add chatbot behavior
	if chatbot.Behaviour != "" {
		promptBuilder.WriteString(fmt.Sprintf("You are a chatbot with the following behavior: %s\n", chatbot.Behaviour))
	}

	// Add user context
	if chatbot.Usercontext != "" {
		promptBuilder.WriteString(fmt.Sprintf("This is the user context: %s\n", chatbot.Usercontext))
	}

	// Add file content
	if chatbot.File != "" {
		fileContent, err := h.readFileContent(chatbot.File)
		if err != nil {
			log.Printf("Error reading file: %v", err)
			promptBuilder.WriteString("Error: Could not load the file associated with this chatbot. Please check chatbot config.\n")
		} else {
			promptBuilder.WriteString(fmt.Sprintf("Here is some more information that might be useful:\n%s\n", fileContent))
		}
	}

	// Add user's message
	promptBuilder.WriteString(fmt.Sprintf("User: %s\n", userMessage))
	promptBuilder.WriteString("Assistant: ")

	return promptBuilder.String()
}

// readFileContent reads the content of the uploaded file.
func (h *Handler) readFileContent(filepath string) (string, error) {
	content, err := os.ReadFile(filepath)
	if err != nil {
		return "", err
	}
	return string(content), nil
}

// sendPromptToGemini sends the prompt to the Gemini API and returns the response.
func (h *Handler) sendPromptToGemini(ctx context.Context, prompt string) (string, error) {
	resp, err := h.model.GenerateContent(ctx, generativeai.Text(prompt))
	if err != nil {
		return "", fmt.Errorf("error generating response: %w", err)
	}

	var result strings.Builder
	for _, part := range resp.Candidates[0].Content.Parts {
		result.WriteString(string(part.(generativeai.Text)))
	}
	return result.String(), nil
}

func (h *Handler) Close() error {
	return h.genaiClient.Close()
}
*/
