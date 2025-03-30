package user

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Owen-Choh/SC4052-Cloud-Computing-Assignment-2/chatbot-backend/chatbot/types"
)

func TestUserServiceRegisterHandler(t *testing.T) {
	userStore := &mockUserStore{}
	handler := NewHandler(userStore)

	t.Run("should fail if user payload invalid", func (t *testing.T) {
		payload := types.RegisterUserPayload{
			Username: "test-user",
			Password: "test-password",
		}

		marshalled, _ := json.Marshal(payload)

		request, err := http.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(marshalled))
		if err != nil {
			t.Fatal(err)
		}

			// Execute test request
			responseRecorder := httptest.NewRecorder()
			router := http.NewServeMux()
			router.HandleFunc("/register", handler.handleRegister)
			router.ServeHTTP(responseRecorder, request)

			// Check response
			if responseRecorder.Code != test.expected {
				t.Errorf("expected status code %d, got %d %s", test.expected, responseRecorder.Code, responseRecorder.Body.String())
			}
		})
	}
}

type mockUserStore struct{}

func (m *mockUserStore) GetUserByName(username string) (*types.User, error) {
	return nil, fmt.Errorf("user %s does not exist", username)
}

func (m *mockUserStore) GetUserByID(id int) (*types.User, error) {
	return nil, fmt.Errorf("user %d does not exist", id)
}

func (m *mockUserStore) CreateUser(types.RegisterUserPayload) error {
	return nil
}

func (m *mockUserStore) UpdateUserLastlogin(int) error {
	return nil
}
