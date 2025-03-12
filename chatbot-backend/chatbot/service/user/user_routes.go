package user

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/Owen-Choh/SC4052-Cloud-Computing-Assignment-2/chatbot-backend/chatbot/auth"
	"github.com/Owen-Choh/SC4052-Cloud-Computing-Assignment-2/chatbot-backend/chatbot/config"
	"github.com/Owen-Choh/SC4052-Cloud-Computing-Assignment-2/chatbot-backend/chatbot/types"
	"github.com/Owen-Choh/SC4052-Cloud-Computing-Assignment-2/chatbot-backend/utils"
	"github.com/go-playground/validator/v10"
)

type Handler struct {
	store types.UserStoreInterface
}

func NewHandler(store types.UserStoreInterface) *Handler {
	return &Handler{store: store}
}

func (h *Handler) RegisterRoutes(router *http.ServeMux) {
	router.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello from user")
	})
	router.HandleFunc("POST /login", h.handleLogin)
	router.HandleFunc("POST /register", h.handleRegister)

	// admin routes
}

func (h *Handler) handleLogin(w http.ResponseWriter, r *http.Request) {
	// Parse the form for both application/x-www-form-urlencoded and multipart/form-data
	if err := r.ParseMultipartForm(1000); err != nil {
		log.Println("Error parsing login form:", err)
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	payload := types.LoginUserPayload{
		Username: r.FormValue("username"),
		Password: r.FormValue("password"),
	}

	if err := utils.Validate.Struct(payload); err != nil {
		validate_error := err.(validator.ValidationErrors)
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid payload %v", validate_error))
		return
	}

	u, err := h.store.GetUserByName(payload.Username)
	if err != nil {
		log.Printf("error querying by username: %s\n", err)
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("not found, invalid email or password"))
		return
	}

	if !auth.ComparePassword(u.Password, []byte(payload.Password)) {
		log.Printf("someone tried to login with wrong password\n")
		utils.WriteError(w, http.StatusUnauthorized, fmt.Errorf("not found, invalid email or password"))
		return
	}

	secret := []byte(config.Envs.JWTSecret)
	token, err := auth.CreateJWT(secret, u.Userid, u.Username)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "token",
    Value:    token,
    HttpOnly: true,
    Secure:   true, // Ensure it's only sent over HTTPS
    Path:     "/",
    Expires:  time.Now().Add(auth.GetExpirationDuration()), // 1 day expiry
	})

	utils.WriteJSON(w, http.StatusOK, map[string]string{
		"userid":   strconv.Itoa(u.Userid),
		"username": u.Username,
	})
}

func (h *Handler) handleRegister(w http.ResponseWriter, r *http.Request) {
	// Parse the form for both application/x-www-form-urlencoded and multipart/form-data
	if err := r.ParseMultipartForm(1000); err != nil {
		log.Println("Error parsing login form:", err)
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	payload := types.LoginUserPayload{
		Username: r.FormValue("username"),
		Password: r.FormValue("password"),
	}
	
	if err := utils.Validate.Struct(payload); err != nil {
		validate_error := err.(validator.ValidationErrors)
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid payload %v", validate_error))
		return
	}

	// check if user exists
	_, err := h.store.GetUserByName(payload.Username)
	if err == nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("user %s already exists", payload.Username))
		return
	}

	hashedPassword, err := auth.HashPassword(payload.Password)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	err = h.store.CreateUser(types.RegisterUserPayload{
		Username: payload.Username,
		Password: hashedPassword,
	})
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusCreated, nil)
}
