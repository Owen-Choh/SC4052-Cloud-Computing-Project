package chatbotservice

import (
	"net/http"

	"github.com/Owen-Choh/SC4052-Cloud-Computing-Assignment-2/chatbot-backend/chatbot/types"
)

type Handler struct {
	store *types.ChatbotStoreInterface
}

func NewHandler(store *types.ChatbotStoreInterface) *Handler {
	return &Handler{store: store}
}

func (h *Handler) RegisterRoutes(router *http.ServeMux) {
	router.HandleFunc("GET /chatbot/{chatbotName}", h.GetChatbot)
	router.HandleFunc("POST /chatbot/", h.CreateChatbot)
}

func (h *Handler) GetChatbot(w http.ResponseWriter, r *http.Request) {

}

func (h *Handler) CreateChatbot(w http.ResponseWriter, r *http.Request) {
	
}