package chatbotservice

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/Owen-Choh/SC4052-Cloud-Computing-Assignment-2/chatbot-backend/chatbot/auth"
	"github.com/Owen-Choh/SC4052-Cloud-Computing-Assignment-2/chatbot-backend/chatbot/config"
	"github.com/Owen-Choh/SC4052-Cloud-Computing-Assignment-2/chatbot-backend/chatbot/types"
	"github.com/Owen-Choh/SC4052-Cloud-Computing-Assignment-2/chatbot-backend/utils"
	"github.com/Owen-Choh/SC4052-Cloud-Computing-Assignment-2/chatbot-backend/utils/validate"
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

	for index := range chatbots {
		chatbots[index].Filepath = filepath.Base(chatbots[index].Filepath)
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

	chatbot.Filepath = filepath.Base(chatbot.Filepath)
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
	chatbotname := strings.TrimSpace(r.FormValue("chatbotname"))
	description := strings.TrimSpace(r.FormValue("description"))
	behaviour := strings.TrimSpace(r.FormValue("behaviour"))
	usercontext := strings.TrimSpace(r.FormValue("usercontext"))
	isShared := r.FormValue("isShared") == "true"

	// Handle file upload, get the paths first for validation
	file, header, err := r.FormFile("file")
	var fullDirPath string
	var filepath string
	if err == nil {
		defer file.Close()

		// Read the first 512 bytes to detect content type
		buffer := make([]byte, 512)
		_, err = file.Read(buffer)
		if err != nil {
			http.Error(w, "Failed to read file", http.StatusInternalServerError)
			return
		}
		fileType := http.DetectContentType(buffer)

		// Reset file reader position (important for saving later)
		file.Seek(0, 0)

		// Allowed file types
		allowedTypes := map[string]bool{
			"application/pdf": true, // PDF
			"image/jpeg":      true, // JPEG
		}
		// Validate file type
		if !allowedTypes[fileType] {
			utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid file type. Only PDF, JPG or JPEG are allowed."))
			return
		}

		fullDirPath = config.Envs.FILES_PATH + username + "/" + chatbotname
		filepath = fullDirPath + "/" + header.Filename
	} else {
		filepath = ""
	}

	var fileUpdatedDate string
	if filepath != "" {
		fileUpdatedDate, _ = utils.GetCurrentTime()
	} else {
		fileUpdatedDate = ""
	}

	// Create chatbot struct to validate fields first
	newChatbot := types.NewChatbot{
		Username:        username,
		Chatbotname:     chatbotname,
		Description:     description,
		Behaviour:       behaviour,
		IsShared:        isShared,
		Usercontext:     usercontext,
		File:            filepath,
		FileUpdatedDate: fileUpdatedDate,
	}
	if err := utils.Validate.Struct(newChatbot); err != nil {
		validate_error := err.(validator.ValidationErrors)
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid payload %v", validate_error))
		return
	}
	// check chatbot name, cannot have some special characters
	if chatbotname == "" || !validate.ValidChatbotNameRegex.MatchString(chatbotname) {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid chatbot name"))
		return
	}
	// check file name, cannot have some special characters
	if filepath != "" && !validate.ValidFileNameRegex.MatchString(header.Filename) {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid file name"))
		return
	}

	// Handle file upload
	if fullDirPath != "" && filepath != "" {
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
	}

	createdTime, _ := utils.GetCurrentTime()
	botID, err := h.chatbotStore.CreateChatbot(newChatbot)
	if err != nil {
		log.Println("Error creating chatbot:", err)
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusCreated, map[string]interface{}{
		"chatbotid":   botID,
		"createddate": createdTime,
		"updateddate": createdTime,
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
	chatbotname := strings.TrimSpace(r.FormValue("chatbotname"))
	description := strings.TrimSpace(r.FormValue("description"))
	behaviour := strings.TrimSpace(r.FormValue("behaviour"))
	usercontext := strings.TrimSpace(r.FormValue("usercontext"))
	isShared := r.FormValue("isShared") == "true"
	removeFile := r.FormValue("removeFile") == "true"

	// Handle file upload, get the paths first for validation
	file, header, err := r.FormFile("file")
	var fullDirPath string
	var newFilepath string
	if err == nil {
		defer file.Close()

		// Read the first 512 bytes to detect content type
		buffer := make([]byte, 512)
		_, err = file.Read(buffer)
		if err != nil {
			http.Error(w, "Failed to read file", http.StatusInternalServerError)
			return
		}
		fileType := http.DetectContentType(buffer)

		// Reset file reader position (important for saving later)
		file.Seek(0, 0)

		// Allowed file types
		allowedTypes := map[string]bool{
			"application/pdf": true, // PDF
			"image/jpeg":      true, // JPEG
		}
		// Validate file type
		if !allowedTypes[fileType] {
			utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid file type. Only PDF, JPG or JPEG are allowed."))
			return
		}

		fullDirPath = config.Envs.FILES_PATH + username + "/" + chatbotname
		newFilepath = fullDirPath + "/" + header.Filename
		log.Println("Saving file to:", newFilepath)
	} else {
		newFilepath = ""
	}

	var fileUpdatedDate string
	if newFilepath != "" {
		fileUpdatedDate, _ = utils.GetCurrentTime()
	} else {
		fileUpdatedDate = oldChatbot.FileUpdatedDate
	}

	// Create chatbot struct
	updateChatbot := types.UpdateChatbot{
		Chatbotid:       chatbotIDInt,
		Username:        username,
		Chatbotname:     chatbotname,
		Description:     description,
		Behaviour:       behaviour,
		IsShared:        isShared,
		Usercontext:     usercontext,
		File:            newFilepath,
		FileUpdatedDate: fileUpdatedDate,
	}
	if err := utils.Validate.Struct(updateChatbot); err != nil {
		validate_error := err.(validator.ValidationErrors)
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid payload %v", validate_error))
		return
	}
	// check chatbot name, cannot have some special characters
	if chatbotname == "" || !validate.ValidChatbotNameRegex.MatchString(chatbotname) {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid chatbot name"))
		return
	}
	// check file name, cannot have some special characters
	if newFilepath != "" && !validate.ValidFileNameRegex.MatchString(header.Filename) {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid file name"))
		return
	}

	// Handle file removal
	oldfilepath := oldChatbot.Filepath
	if removeFile || newFilepath != "" {
		if oldfilepath != "" {
			log.Println("Attempting to remove file:", oldfilepath)
			err = os.Remove(oldfilepath)
			if err != nil {
				if err != os.ErrNotExist {
					// Log the error but can continue
					log.Println("Error removing file but continuing:", oldfilepath, err)
					if removeFile && newFilepath == "" {
						// if user requested to remove file, then just skip to the end of the handler function behaviour
						log.Println("Error removing file ending early since user requested to remove file only")

						updateTime, _ := utils.GetCurrentTime()
						err = h.chatbotStore.UpdateChatbot(updateChatbot)
						if err != nil {
							log.Println("Error updating chatbot:", err)
							utils.WriteError(w, http.StatusInternalServerError, err)
							return
						}
						utils.WriteJSON(w, http.StatusOK, map[string]interface{}{
							"message":     "Chatbot updated successfully",
							"updateddate": updateTime,
						})
						return
					}
				} else {
					log.Println("Error removing file:", err)

					// if user only requested to remove file, then do not proceed
					if removeFile && newFilepath == "" {
						utils.WriteError(w, http.StatusInternalServerError, fmt.Errorf("failed to remove previously uploaded file"))
						return
					}
				}
			}
			log.Println("Removed old file:", oldfilepath)
		}
	}

	// Handle file upload
	if err == nil {
		// fullDirPath := config.Envs.FILES_PATH + username + "/" + chatbotname
		log.Println("Full directory path:", fullDirPath)
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
		// newFilepath = fullDirPath + "/" + header.Filename
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

	// set back to the old path so that db dont get updated with empty string
	if updateChatbot.File == "" && oldfilepath != "" {
		updateChatbot.File = oldfilepath
	}
	updateTime, _ := utils.GetCurrentTime()
	err = h.chatbotStore.UpdateChatbot(updateChatbot)
	if err != nil {
		log.Println("Error updating chatbot:", err)
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, map[string]interface{}{
		"message":     "Chatbot updated successfully",
		"updateddate": updateTime,
	})
}

func (h *Handler) DeleteChatbot(w http.ResponseWriter, r *http.Request) {
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

	oldfilepath := chatbot.Filepath
	if oldfilepath != "" {
		err = os.RemoveAll(config.Envs.FILES_PATH + chatbot.Username + "/" + chatbot.Chatbotname)
		if err != nil {
			log.Println("Error removing directory:", err)
			utils.WriteError(w, http.StatusInternalServerError, fmt.Errorf("failed to remove directory of chatbot"))
			return
		}
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
