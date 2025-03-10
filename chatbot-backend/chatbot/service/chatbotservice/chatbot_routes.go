package chatbotservice

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/Owen-Choh/SC4052-Cloud-Computing-Assignment-2/chatbot-backend/chatbot/types"
	"github.com/Owen-Choh/SC4052-Cloud-Computing-Assignment-2/chatbot-backend/utils"
	"github.com/go-playground/validator/v10"
)

var ErrChatbotNotFound = errors.New("chatbot not found")

type Handler struct {
	store types.ChatbotStoreInterface
}

func NewHandler(store types.ChatbotStoreInterface) *Handler {
	return &Handler{store: store}
}

func (h *Handler) RegisterRoutes(router *http.ServeMux) {
	router.HandleFunc("GET /chatbot/{username}/{chatbotName}", h.GetChatbot)
	router.HandleFunc("POST /chatbot/", h.CreateChatbot)
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
		if errors.Is(err, ErrChatbotNotFound){
			utils.WriteError(w, http.StatusNotFound, err)
		} else {
			utils.WriteError(w, http.StatusInternalServerError, err)
		}
		return
	}

	utils.WriteJSON(w, http.StatusOK, chatbot)
}

func (h *Handler) CreateChatbot(w http.ResponseWriter, r *http.Request) {
	
}