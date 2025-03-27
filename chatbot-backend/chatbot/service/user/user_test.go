package user

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Owen-Choh/SC4052-Cloud-Computing-Assignment-2/chatbot-backend/chatbot/types"
)

func TestUserServiceHandlers(t *testing.T) {
	userStore := &mockUserStore{}
	handler := NewHandler(userStore)

	t.Run("should fail if user payload invalid", func(t *testing.T) {
		payload := types.RegisterUserPayload{
			Username: "test-user",
			Password: "test-password",
		}

		marshalled, _ := json.Marshal(payload)

		request, err := http.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(marshalled))
		if err != nil {
			t.Fatal(err)
		}

		responseRecorder := httptest.NewRecorder()
		router := http.NewServeMux()
		router.HandleFunc("/register", handler.handleRegister)
		router.ServeHTTP(responseRecorder, request)

		if responseRecorder.Code != http.StatusBadRequest {
			t.Errorf("expected status code %d, got %d", http.StatusBadRequest, responseRecorder.Code)
		}
	})
}

type mockUserStore struct{}

func (m *mockUserStore) GetUserByName(username string) (*types.User, error) {
	return nil, nil
}

func (m *mockUserStore) GetUserByID(id int) (*types.User, error) {
	return nil, nil
}
func (m *mockUserStore) CreateUser(types.RegisterUserPayload) error {
	return nil
}
