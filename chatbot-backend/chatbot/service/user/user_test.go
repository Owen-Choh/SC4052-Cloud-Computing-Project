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

	tests := []struct {
		name        string
		payload     types.RegisterUserPayload
		payloadType string
		expected    int
	}{
		{
			name: "json payload",
			payload: types.RegisterUserPayload{
				Username: "testuser",
				Password: "test-password",
			},
			payloadType: "application/json",
			expected:    http.StatusBadRequest, // Assuming JSON is not accepted
		},
		{
			name: "multipart form payload",
			payload: types.RegisterUserPayload{
				Username: "testuser",
				Password: "test-password",
			},
			payloadType: "multipart/form-data",
			expected:    http.StatusCreated, // Assuming successful registration
		},
		{
			name: "username with special characters",
			payload: types.RegisterUserPayload{
				Username: "test\\user",
				Password: "test-password",
			},
			payloadType: "multipart/form-data",
			expected:    http.StatusBadRequest, // Assuming successful registration
		},
		{
			name: "username with spaces",
			payload: types.RegisterUserPayload{
				Username: "test user",
				Password: "test-password",
			},
			payloadType: "multipart/form-data",
			expected:    http.StatusBadRequest, // Assuming successful registration
		},
		{
			name: "password too short",
			payload: types.RegisterUserPayload{
				Username: "testuser",
				Password: "test",
			},
			payloadType: "multipart/form-data",
			expected:    http.StatusBadRequest, // Assuming successful registration
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var requestBody bytes.Buffer
			var contentType string

			if test.payloadType == "application/json" {
				jsonData, _ := json.Marshal(test.payload)
				requestBody.Write(jsonData)
				contentType = "application/json"
			} else if test.payloadType == "multipart/form-data" {
				writer := multipart.NewWriter(&requestBody)
				_ = writer.WriteField("username", test.payload.Username)
				_ = writer.WriteField("password", test.payload.Password)
				writer.Close() // Must close before using data
				contentType = writer.FormDataContentType()
			} else {
				log.Fatal("Invalid payload type")
			}

			// Create request with the correct Content-Type
			request, err := http.NewRequest(http.MethodPost, "/register", &requestBody)
			if err != nil {
				t.Fatal(err)
			}
			request.Header.Set("Content-Type", contentType)

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
