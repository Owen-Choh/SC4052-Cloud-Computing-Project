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
	router.HandleFunc("GET /chatbot/{chatbotName}", h.GetChatbot)
	router.HandleFunc("POST /chatbot/", h.CreateChatbot)
}

func (h *Handler) GetChatbot(w http.ResponseWriter, r *http.Request) {
	var payload types.GetChatbotPayload
	if err := utils.ParseJSON(r, payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	if err := utils.Validate.Struct(payload); err != nil {
		validate_error := err.(validator.ValidationErrors)
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid payload %v", validate_error))
		return
	}

	chatbot, err := h.store.GetChatbotByName(payload.Username, payload.Chatbotname)
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