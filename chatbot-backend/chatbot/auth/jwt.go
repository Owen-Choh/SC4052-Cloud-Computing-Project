package auth

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/Owen-Choh/SC4052-Cloud-Computing-Assignment-2/chatbot-backend/chatbot/config"
	"github.com/Owen-Choh/SC4052-Cloud-Computing-Assignment-2/chatbot-backend/chatbot/types"
	"github.com/Owen-Choh/SC4052-Cloud-Computing-Assignment-2/chatbot-backend/utils"
	"github.com/golang-jwt/jwt/v5"
)

type contextKey string

const CookieName = "token"
const UserIDKey contextKey = "userid"
const UsernameKey contextKey = "username"

func CreateJWT(secret []byte, userid int, username string) (string, error) {
	expiration := GetExpirationDuration()

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userid":    strconv.Itoa(userid),
		"username":  username,
		"expiredAt": time.Now().Add(expiration).Unix(),
	})

	tokenString, err := token.SignedString(secret)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func WithJWTAuth(handlerFunc http.HandlerFunc, store types.UserStoreInterface) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tokenString := GetTokenFromRequest(r)
		if tokenString == "" {
			utils.WriteError(w, http.StatusTeapot, fmt.Errorf("token missing, permission denied"))
			return
		}
		// tokenString = tokenString[7:] //remove the bearer prefix
		token, err := validateToken(tokenString)
		if err != nil {
			log.Printf("failed to validate token: %v", err)
			permissionDenied(w)
			return
		}

		if !token.Valid {
			log.Printf("invalid token")
			permissionDenied(w)
			return
		}

		claims := token.Claims.(jwt.MapClaims)
		str := claims["userid"].(string)
		userID, err := strconv.Atoi(str)
		if err != nil {
			log.Printf("failed to convert userID to int: %v", err)
			permissionDenied(w)
			return
		}

		u, err := store.GetUserByID(userID)
		if err != nil {
			log.Printf("failed to get user by id: %v", err)
			permissionDenied(w)
			return
		}

		if u == nil {
			log.Printf("user not found")
			permissionDenied(w)
			return
		}

		if u.Username != claims["username"].(string) || u.Userid != userID {
			log.Printf("jwt claims mismatched for userid %d, wrong userid %t, wrong username %t ", u.Userid, u.Userid != userID, u.Username != claims["username"].(string))
			permissionDenied(w)
			return
		}

		ctx := r.Context()
		ctx = context.WithValue(ctx, UserIDKey, u.Userid)
		ctx = context.WithValue(ctx, UsernameKey, u.Username)
		r = r.WithContext(ctx)

		handlerFunc(w, r)
	}
}

func GetTokenFromRequest(r *http.Request) string {
	// This code is getting the token from cookie
	token, err := r.Cookie("token")
	if err != nil {
		return ""
	}
	return token.Value
	// This code is getting the token from header
	// tokenAuth := r.Header.Get("Authorization")
	// if tokenAuth != "" {
	// 	return strings.TrimSpace(tokenAuth)
	// }

	// return ""
}

func validateToken(t string) (*jwt.Token, error) {
	token, err := jwt.Parse(t, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}

		return []byte(config.Envs.JWTSecret), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		expiredAt := int64(claims["expiredAt"].(float64))
		if time.Now().Unix() > expiredAt {
			return nil, fmt.Errorf("token has expired")
		}
	}

	return token, nil
}

func permissionDenied(w http.ResponseWriter) {
	utils.WriteError(w, http.StatusForbidden, fmt.Errorf("permission denied"))
}

func GetUserIDFromContext(ctx context.Context) int {
	userID, ok := ctx.Value(UserIDKey).(int)
	if !ok {
		return -1
	}
	return userID
}

func GetUsernameFromContext(ctx context.Context) string {
	username, ok := ctx.Value(UsernameKey).(string)
	if !ok {
		return ""
	}
	return username
}

func GetExpirationDuration() time.Duration {
	return time.Second * time.Duration(config.Envs.JWTExpirationInSeconds)
}
