package chatbotservice

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/Owen-Choh/SC4052-Cloud-Computing-Assignment-2/chatbot-backend/chatbot/auth"
	"github.com/Owen-Choh/SC4052-Cloud-Computing-Assignment-2/chatbot-backend/chatbot/types"
	"github.com/Owen-Choh/SC4052-Cloud-Computing-Assignment-2/chatbot-backend/utils"
	"github.com/go-playground/validator/v10"
)

var ErrChatbotNotFound = errors.New("chatbot not found")

type Handler struct {
	chatbotStore types.ChatbotStoreInterface
	userStore    types.UserStoreInterface
}

func NewHandler(chatbotStore types.ChatbotStoreInterface, userstore types.UserStoreInterface) *Handler {
	return &Handler{
		chatbotStore: chatbotStore,
		userStore:    userstore,
	}
}

func (h *Handler) RegisterRoutes(router *http.ServeMux) {
	router.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello from chatbot")
	})
	router.HandleFunc("GET /list", auth.WithJWTAuth(h.GetUserChatbot, h.userStore))
	router.HandleFunc("GET /details/{username}/{chatbotName}", h.GetChatbot)
	router.HandleFunc("POST /", auth.WithJWTAuth(h.CreateChatbot, h.userStore))
	router.HandleFunc("PUT /{chatbotid}", auth.WithJWTAuth(h.UpdateChatbot, h.userStore))
	router.HandleFunc("DELETE /{chatbotid}", auth.WithJWTAuth(h.DeleteChatbot, h.userStore))
}

func (h *Handler) GetUserChatbot(w http.ResponseWriter, r *http.Request) {
	username := auth.GetUsernameFromContext(r.Context())
	if username == "" {
		log.Println("username missing in request context")
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid request"))
		return
	}

	log.Println("authenticated user: " + username)
	chatbots, err := h.chatbotStore.GetChatbotsByUsername(username)
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

	chatbot, err := h.chatbotStore.GetChatbotByName(username, chatbotName)
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
		out, err := os.OpenFile(filepath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0666)
		// out, err := os.Create(filepath)
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

	botID, err := h.chatbotStore.CreateChatbot(newChatbot)
	if err != nil {
		log.Println("Error creating chatbot:", err)
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusCreated, map[string]interface{}{
		"chatbotid": botID,
	})
}

func (h *Handler) UpdateChatbot(w http.ResponseWriter, r *http.Request) {
	username := auth.GetUsernameFromContext(r.Context())
	if username == "" {
		log.Println("username missing in request context set by jwt")
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid request"))
		return
	}

	chatbotID := r.PathValue("chatbotid")
	if chatbotID == "" {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid chatbot ID"))
		return
	}
	chatbotIDInt, converr := strconv.Atoi(chatbotID)
	if converr != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid chatbot ID"))
		return
	}

	err := r.ParseMultipartForm(10 << 20) // 10MB limit
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("failed to parse form"))
		return
	}

	oldChatbot, err := h.chatbotStore.GetChatbotsByID(chatbotIDInt)
	if err != nil {
		log.Printf("Error getting chatbot %d: %v\n", chatbotIDInt, err)
		utils.WriteError(w, http.StatusForbidden, fmt.Errorf("unauthorized"))
		return
	}
	if username != oldChatbot.Username {
		utils.WriteError(w, http.StatusForbidden, fmt.Errorf("unauthorized"))
		return
	}

	// Extract chatbot fields from form
	chatbotname := r.FormValue("chatbotname")
	description := r.FormValue("description")
	behaviour := r.FormValue("behaviour")
	usercontext := r.FormValue("usercontext")
	isShared := r.FormValue("isShared") == "true"
	removeFile := r.FormValue("removeFile") == "true"

	// Handle file removal
	oldfilepath := oldChatbot.Filepath
	if removeFile {
		if oldfilepath != "" {
			log.Println("Attempting to remove file:", oldfilepath)
			// Check if file is locked
			f, err := os.OpenFile(oldfilepath, os.O_RDWR, 0666)
			if err != nil {
				log.Println("File seems locked or inaccessible:", err)
			} else {
				f.Close()
			}

			err = os.Remove(oldfilepath)
			if err != nil {
				log.Println("Error removing file:", err)
				utils.WriteError(w, http.StatusInternalServerError, fmt.Errorf("failed to remove file"))
				return
			}
		}
	}

	// Handle file upload
	file, header, err := r.FormFile("file")
	var newFilepath string
	if err == nil {
		defer file.Close()

		fullDirPath := "database_files/uploads/" + username + "/" + chatbotname
		err := os.MkdirAll(fullDirPath, os.ModePerm) // Create the directory if it doesn’t exist
		if err != nil {
			log.Println("Error creating directory:", err)
			utils.WriteError(w, http.StatusInternalServerError, fmt.Errorf("failed to save new file"))
			return
		}

		if header.Filename == "" {
			log.Println("No filename detected")
			utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("failed to save new file"))
			return
		}
		// Save the uploaded file
		newFilepath = fullDirPath + "/" + header.Filename
		log.Println("Saving file to:", newFilepath)
		out, err := os.Create(newFilepath)
		if err != nil {
			log.Println("Error saving file:", err)
			utils.WriteError(w, http.StatusInternalServerError, fmt.Errorf("failed to save new file"))
			return
		}
		defer out.Close()
		io.Copy(out, file)
	} else {
		newFilepath = "" // No file uploaded
	}

	if newFilepath == "" {
		newFilepath = oldfilepath
	}

	// Create chatbot struct
	updateChatbot := types.UpdateChatbot{
		Chatbotid:   chatbotIDInt,
		Username:    username,
		Chatbotname: chatbotname,
		Description: description,
		Behaviour:   behaviour,
		IsShared:    isShared,
		Usercontext: usercontext,
		File:        newFilepath,
	}
	if err := utils.Validate.Struct(updateChatbot); err != nil {
		validate_error := err.(validator.ValidationErrors)
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid payload %v", validate_error))
		return
	}

	err = h.chatbotStore.UpdateChatbot(updateChatbot)
	if err != nil {
		log.Println("Error updating chatbot:", err)
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, map[string]interface{}{
		"message": "Chatbot updated successfully",
	})
}

func (h *Handler) DeleteChatbot(w http.ResponseWriter, r *http.Request){
	username := auth.GetUsernameFromContext(r.Context())
	if username == "" {
		log.Println("username missing in request context set by jwt")
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid request"))
		return
	}

	chatbotID := r.PathValue("chatbotid")
	if chatbotID == "" {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid chatbot ID"))
		return
	}
	chatbotIDInt, converr := strconv.Atoi(chatbotID)
	if converr != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid chatbot ID"))
		return
	}

	chatbot, err := h.chatbotStore.GetChatbotsByID(chatbotIDInt)
	if err != nil {
		log.Printf("Error getting chatbot %d for deletion: %v\n", chatbotIDInt, err)
		utils.WriteError(w, http.StatusForbidden, fmt.Errorf("unauthorized"))
		return
	}
	if username != chatbot.Username {
		utils.WriteError(w, http.StatusForbidden, fmt.Errorf("unauthorized"))
		return
	}

	err = h.chatbotStore.DeleteChatbot(chatbotIDInt)
	if err != nil {
		log.Println("Error deleting chatbot:", err)
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, map[string]interface{}{
		"message": "Chatbot deleted successfully",
	})
}