package auth

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Owen-Choh/SC4052-Cloud-Computing-Assignment-2/chatbot-backend/chatbot/config"
	"github.com/Owen-Choh/SC4052-Cloud-Computing-Assignment-2/chatbot-backend/chatbot/types"
	"github.com/golang-jwt/jwt/v5"
)

// secret: jwt-secret
//
//	"userid":    2,
//	"username":  "testuser",
//	"expiredAt": "2125-03-20T12:18:41+08:00"
const invalidUserToken string = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyaWQiOjIsInVzZXJuYW1lIjoidGVzdHVzZXIiLCJleHBpcmVkQXQiOiIyMTI1LTAzLTIwVDEyOjE4OjQxKzA4OjAwIn0.BBAGT2RXbt1OZo67Mq6iEsUMl4ScEZzi0c2FyycDb7U"

// "userid":    1,
// "username":  "test-user",
// "expiredAt": "2025-03-20T12:18:41+08:00"
const expiredToken string = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyaWQiOjEsInVzZXJuYW1lIjoidGVzdC11c2VyIiwiZXhwaXJlZEF0IjoiMjAyNS0wMy0yMFQxMjoxODo0MSswODowMCJ9.4NGPWwdbRBgM3AE5Uvv70YpWKFLZct89P5u9LXYOTYw"

func TestCreateJWT(t *testing.T) {
	secret := []byte("jwt-secret")

	token, err := CreateJWT(secret, 1, "test user")
	if err != nil {
		t.Errorf("error creating jwt: %v", err)
	}

	if token == "" {
		t.Error("expected token to be not empty")
	}
}

func TestGetTokenFromRequest(t *testing.T) {
	tests := []struct {
		name         string
		req          *http.Request
		cookieString string
		expected     string
	}{
		{
			name: "valid cookie",
			req: &http.Request{
				Header: http.Header{"Authorization": []string{"Bearer test-token"}},
			},
			cookieString: "test-token",
			expected:     "test-token",
		},
		{
			name:     "no input",
			req:      nil,
			expected: "",
		},
		{
			name: "no cookie",
			req: &http.Request{
				Header: http.Header{"Authorization": []string{"Bearer test-token"}},
			},
			cookieString: "nil",
			expected:     "",
		},
		{
			name: "empty cookie",
			req: &http.Request{
				Header: http.Header{"Authorization": []string{"Bearer test-token"}},
			},
			cookieString: "",
			expected:     "",
		},
		{
			name: "trimmed token",
			req: &http.Request{
				Header: http.Header{"Authorization": []string{"Bearer test-token"}},
			},
			cookieString: "   test-token   ",
			expected:     "test-token",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if test.req != nil && test.cookieString != "nil" {
				cookie := &http.Cookie{
					Name:  CookieName,
					Value: test.cookieString,
				}
				test.req.AddCookie(cookie)
			}

			token := GetTokenFromRequest(test.req)
			if token != test.expected {
				t.Errorf("expected token %s, got %s", test.expected, token)
			}
		})
	}
}

func TestValidateToken(t *testing.T) {

	t.Run("valid token", func(t *testing.T) {
		secret := []byte(config.Envs.JWTSecret)
		token, err := CreateJWT(secret, 1, "test user")
		if err != nil {
			t.Errorf("error creating jwt: %v", err)
		}
		validatedToken, err := validateToken(token)
		if err != nil {
			t.Fatalf("expected error for token: %v", err)
		}
		if validatedToken == nil {
			t.Fatal("expected non-nil token")
		}

		if !validatedToken.Valid {
			t.Error("expected valid token")
		}

		if validatedToken.Claims.(jwt.MapClaims)["username"] != "test user" {
			t.Errorf("expected username to be test user, got %v", validatedToken.Claims.(jwt.MapClaims)["username"])
		}
		if validatedToken.Claims.(jwt.MapClaims)["userid"] != "1" {
			t.Errorf("expected userid to be string 1, got %T %v", validatedToken.Claims.(jwt.MapClaims)["userid"], validatedToken.Claims.(jwt.MapClaims)["userid"])
		}
	})

	t.Run("different secret", func(t *testing.T) {
		secret := []byte("jwt-secret")
		token, err := CreateJWT(secret, 1, "testuser")
		if err != nil {
			t.Errorf("error creating jwt: %v", err)
		}
		validatedToken, err := validateToken(token)
		if err == nil {
			t.Error("expected error for token created with different secret")
		}
		if validatedToken != nil {
			t.Errorf("expected nil token for token created with different secret, got %v", validatedToken)
		}
	})

	t.Run("invalid token", func(t *testing.T) {
		// invalid token
		_, err := validateToken("invalid-token")
		if err == nil {
			t.Error("expected error for invalid token")
		}
	})

	t.Run("expired token", func(t *testing.T) {
		_, err := validateToken(expiredToken)
		if err == nil {
			t.Error("expected error for expired token")
		}
	})
}

func TestValidateTokenMiddleware(t *testing.T) {
	secret := []byte(config.Envs.JWTSecret)
	validToken, err := CreateJWT(secret, 1, "test-user")
	if err != nil {
		t.Fatalf("error creating validToken jwt: %v", err)
	}
	mismatchUsernameToken, err := CreateJWT(secret, 1, "test user")
	if err != nil {
		t.Fatalf("error creating mismatchUsernameToken jwt: %v", err)
	}

	tests := []struct {
		name       string
		authHeader bool
		token      string
		expected   int
	}{
		{
			name:       "valid token",
			authHeader: true,
			token:      validToken,
			expected:   http.StatusOK,
		},
		{
			name:       "no header",
			authHeader: false,
			token:      "",
			expected:   http.StatusForbidden,
		},
		{
			name:       "no token",
			authHeader: true,
			token:      "",
			expected:   http.StatusForbidden,
		},
		{
			name:       "invalid user",
			authHeader: true,
			token:      invalidUserToken,
			expected:   http.StatusForbidden,
		},
		{
			name:       "invalid token",
			authHeader: true,
			token:      expiredToken,
			expected:   http.StatusForbidden,
		},
		{
			name:       "wrong username token",
			authHeader: true,
			token:      mismatchUsernameToken,
			expected:   http.StatusForbidden,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			request, err := http.NewRequest(http.MethodGet, "/test", nil)
			if err != nil {
				t.Fatal(err)
			}
			if test.authHeader {
				request.AddCookie(&http.Cookie{
					Name:  CookieName,
					Value: test.token,
				})
			}

			var capturedCtx context.Context
			responseRecorder := httptest.NewRecorder()
			router := http.NewServeMux()
			router.HandleFunc("/test", WithJWTAuth(func(w http.ResponseWriter, r *http.Request) {
				capturedCtx = r.Context()
				w.WriteHeader(http.StatusOK)
			}, &mockUserStore{}))
			router.ServeHTTP(responseRecorder, request)

			if responseRecorder.Code != test.expected {
				t.Errorf("expected status code %d, got %d", test.expected, responseRecorder.Code)
			}

			// check if user details is set in context
			if responseRecorder.Code == http.StatusOK {
				if capturedCtx == nil {
					t.Fatal("expected context to be captured")
				}
				userID := capturedCtx.Value(UserIDKey)
				if userID == nil {
					t.Error("expected user id to be set in context")
				} else {
					if userID != 1 {
						t.Errorf("expected user id to be 1, got %v", userID)
					}
				}
				username := capturedCtx.Value(UsernameKey)
				if username == nil {
					t.Error("expected username to be set in context")
				} else {
					if username != "test-user" {
						t.Errorf("expected username to be test-user, got %v", username)
					}
				}
			}
		})
	}
}

type mockUserStore struct{}

func (m *mockUserStore) GetUserByName(username string) (*types.User, error) {
	return nil, nil
}

func (m *mockUserStore) GetUserByID(id int) (*types.User, error) {
	return &types.User{
		Userid:      1,
		Username:    "test-user",
		Password:    "does not matter",
		Createddate: "2021-03-20T12:00:00+08:00",
		Lastlogin:   "2021-03-20T12:00:00+08:00",
	}, nil
}
func (m *mockUserStore) CreateUser(types.RegisterUserPayload) error {
	return nil
}
