package user

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/Owen-Choh/SC4052-Cloud-Computing-Assignment-2/chatbot-backend/chatbot/auth"
	"github.com/Owen-Choh/SC4052-Cloud-Computing-Assignment-2/chatbot-backend/chatbot/types"
	"github.com/Owen-Choh/SC4052-Cloud-Computing-Assignment-2/chatbot-backend/utils"
)

type Handler struct {
	store types.UserStoreInterface
}

func NewHandler(store types.UserStoreInterface) *Handler {
	return &Handler{store: store}
}

func (h *Handler) RegisterRoutes(router *http.ServeMux) {
	router.HandleFunc("/login", h.handleLogin)
	router.HandleFunc("/register", h.handleRegister)

	// admin routes
}

func (h *Handler) handleLogin(w http.ResponseWriter, r *http.Request) {
	var user types.LoginUserPayload
	if err := utils.ParseJSON(r, &user); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	u, err := h.store.GetUserByName(user.Username)
	if err != nil {
		log.Printf("error querying by username: %s\n", err)
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("not found, invalid email or password"))
		return
	}

	if !auth.ComparePassword(u.Password, []byte(user.Password)) {
		log.Printf("someone tried to login with wrong password\n")
		utils.WriteError(w, http.StatusUnauthorized, fmt.Errorf("not found, invalid email or password"))
		return
	}

	utils.WriteJSON(w, http.StatusOK, map[string]string{
		"userid":   strconv.Itoa(u.Userid),
		"username": u.Username,
		"token": "",
	})
}

func (h *Handler) handleRegister(w http.ResponseWriter, r *http.Request) {
	var user types.RegisterUserPayload
	if err := utils.ParseJSON(r, &user); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	// check if user exists
	_, err := h.store.GetUserByName(user.Username)
	if err == nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("user %s already exists", user.Username))
		return
	}

	hashedPassword, err := auth.HashPassword(user.Password)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	err = h.store.CreateUser(types.RegisterUserPayload{
		Username: user.Username,
		Password: hashedPassword,
	})
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusCreated, nil)
}
