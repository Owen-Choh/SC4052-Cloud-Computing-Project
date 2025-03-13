package chatbotservice

import (
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/Owen-Choh/SC4052-Cloud-Computing-Assignment-2/chatbot-backend/chatbot/auth"
	"github.com/Owen-Choh/SC4052-Cloud-Computing-Assignment-2/chatbot-backend/chatbot/types"
	"github.com/Owen-Choh/SC4052-Cloud-Computing-Assignment-2/chatbot-backend/utils"
)

var ErrChatbotNotFound = errors.New("chatbot not found")

type Handler struct {
	store types.ChatbotStoreInterface
	userstore types.UserStoreInterface
}

func NewHandler(store types.ChatbotStoreInterface, userstore types.UserStoreInterface) *Handler {
	return &Handler{
		store: store,
		userstore: userstore,
	}
}

func (h *Handler) RegisterRoutes(router *http.ServeMux) {
	router.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello from chatbot")})
	router.HandleFunc("GET /list", auth.WithJWTAuth(h.GetUserChatbot, h.userstore))
	router.HandleFunc("GET /{username}/{chatbotName}", h.GetChatbot)
	router.HandleFunc("POST /newchatbot", auth.WithJWTAuth(h.CreateChatbot, h.userstore))
}

func (h *Handler) GetUserChatbot(w http.ResponseWriter, r *http.Request) {
	username := auth.GetUsernameFromContext(r.Context())
	if username == "" {
		log.Println("username missing in request context")
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid request"))
		return
	}
	
	log.Println("authenticated user: "+username)
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
	username := auth.GetUsernameFromContext(r.Context())
	if username == "" {
		log.Println("username missing in request context")
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid request"))
		return
	}

	var payload types.CreateChatbotPayload
	if err := utils.ParseJSON(r, &payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	// newChatbot := types.NewChatbot{
	// 	Username: username,
	// 	Chatbotname: payload.Chatbotname,
	// 	Usercontext: payload.Usercontext,
	// 	File: "",
	// }
	
	// TODO add the chatbot to db
	utils.WriteJSON(w, http.StatusNotImplemented, fmt.Errorf("api is not ready"))

}