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
	router.HandleFunc("GET /logout", h.logout)
	router.HandleFunc("GET /auth/check", auth.WithJWTAuth(h.checkAuth, h.store))
	router.HandleFunc("GET /logout", h.logout)
	router.HandleFunc("GET /auth/check", auth.WithJWTAuth(h.checkAuth, h.store))

	// admin routes
}

func (h *Handler) logout(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    "",            // Empty value
		Path:     "/",           // Match the original path
		HttpOnly: true,
		Secure:   true,          // Keep this for HTTPS
		MaxAge:   -1,            // Tells browser to delete cookie
		Expires:  time.Unix(0, 0), // Optional extra
	})
	utils.WriteJSON(w, http.StatusOK, nil)
}

func (h *Handler) checkAuth(w http.ResponseWriter, r *http.Request) {
	// auth should be handled by middleware, if it reaches here, it means auth is successful
	userid := auth.GetUserIDFromContext(r.Context())
	username := auth.GetUsernameFromContext(r.Context())
	if userid == -1 || username == "" {
		utils.WriteError(w, http.StatusInternalServerError, fmt.Errorf("failed to get user info from request context"))	
		return
	}

	utils.WriteJSON(w, http.StatusOK, map[string]interface{}{
		"user": map[string]interface{}{
			"userid":   strconv.Itoa(userid),
			"username": username,
		},
		"expiresAt": time.Now().Add(auth.GetExpirationDuration()).Format(time.RFC3339),
	})

	log.Printf("checked cookie for user %s\n", username)
}

func (h *Handler) logout(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    "",  // Empty value
		Path:     "/", // Match the original path
		HttpOnly: true,
		Secure:   true,            // Keep this for HTTPS
		MaxAge:   -1,              // Tells browser to delete cookie
		Expires:  time.Unix(0, 0), // Optional extra
	})
	utils.WriteJSON(w, http.StatusOK, nil)
}

func (h *Handler) checkAuth(w http.ResponseWriter, r *http.Request) {
	// auth should be handled by middleware, if it reaches here, it means auth is successful
	userid := auth.GetUserIDFromContext(r.Context())
	username := auth.GetUsernameFromContext(r.Context())
	if userid == -1 || username == "" {
		utils.WriteError(w, http.StatusInternalServerError, fmt.Errorf("failed to get user info from request context"))
		return
	}

	utils.WriteJSON(w, http.StatusOK, map[string]interface{}{
		"user": map[string]interface{}{
			"userid":   strconv.Itoa(userid),
			"username": username,
		},
		"expiresAt": time.Now().Add(auth.GetExpirationDuration()).Format(time.RFC3339),
	})

	log.Printf("checked cookie for user %s\n", username)
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
		log.Printf("error validating login payload %s: %s\n", payload.Username, err)
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid payload %v", validate_error))
		return
	}

	u, err := h.store.GetUserByName(payload.Username)
	if err != nil {
		log.Printf("error querying by username: %s\n", err)
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("not found, invalid username or password"))
		return
	}

	if !auth.ComparePassword(u.Password, []byte(payload.Password)) {
		log.Printf("someone tried to login with wrong password\n")
		utils.WriteError(w, http.StatusUnauthorized, fmt.Errorf("not found, invalid username or password"))
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
    Expires:  time.Now().Add(auth.GetExpirationDuration()),
	})

	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    token,
		HttpOnly: true,
		Secure:   true, // Ensure it's only sent over HTTPS
		Path:     "/",
		Expires:  time.Now().Add(auth.GetExpirationDuration()),
	})

	utils.WriteJSON(w, http.StatusOK, map[string]interface{}{
		"user": map[string]interface{}{
			"userid":   strconv.Itoa(u.Userid),
			"username": u.Username,
		},
		"expiresAt": time.Now().Add(auth.GetExpirationDuration()).Format(time.RFC3339),
	})

	log.Printf("user %s logged in\n", u.Username)

	log.Printf("user %s logged in\n", u.Username)
}

func (h *Handler) handleRegister(w http.ResponseWriter, r *http.Request) {
	// Parse the form for both application/x-www-form-urlencoded and multipart/form-data
	log.Println(r.Header.Get("Content-Type"))
	if err := r.ParseMultipartForm(1000); err != nil {
		log.Println("Error parsing register form:", err)
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	payload := types.RegisterUserPayload{
		Username: r.FormValue("username"),
		Password: r.FormValue("password"),
	}

	if err := utils.Validate.Struct(payload); err != nil {
		validate_error := err.(validator.ValidationErrors)
		log.Printf("error validating register payload %s: %s\n", payload.Username, err)
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid payload %v", validate_error))
		return
	}

	// check if user exists
	_, err := h.store.GetUserByName(payload.Username)
	if err == nil {
		log.Printf("user %s already exists\n", payload.Username)
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
