package auth

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/Harascho1/goWEB/config"
	"github.com/Harascho1/goWEB/types"
	"github.com/Harascho1/goWEB/utils"
	"github.com/golang-jwt/jwt/v5"
)

type contextKey string

const UserKey contextKey = "userID"

func CreateJWT(secret []byte, userId int) (string, error) {
	expiration := time.Second * time.Duration(config.Envs.JWTExpirationInSeconds)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userID":   strconv.Itoa(userId),
		"expireAt": time.Now().Add(expiration).Unix(),
	})

	tokenString, err := token.SignedString(secret)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func GetUserIdFromContext(ctx context.Context) int {
	userId, ok := ctx.Value(UserKey).(int)
	if !ok {
		return -1
	}
	return userId
}

func WithJWTAuth(handerFunc http.HandlerFunc, store types.UserStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tokenString := utils.GetTokenFromRequest(r)
		log.Println(tokenString)

		token, err := validateToken(tokenString)
		if err != nil {
			log.Printf("failed to validated token: %v", err)
			permisionDenied(w)
			return
		}

		if !token.Valid {
			log.Printf("invalid token")
			permisionDenied(w)
			return
		}

		claims := token.Claims.(jwt.MapClaims)
		userIdStr := claims["userID"].(string)

		userId, _ := strconv.Atoi(userIdStr)

		u, err := store.GetUserByID(userId)
		if err != nil {
			log.Printf("failed to get user by id: %v", err)
			permisionDenied(w)
			return
		}

		ctx := r.Context()
		ctx = context.WithValue(ctx, UserKey, u.ID)

		r = r.WithContext(ctx)

		handerFunc(w, r)
	}
}

func permisionDenied(w http.ResponseWriter) {
	utils.WriteError(w, http.StatusForbidden, fmt.Errorf("Permision denied"))
}

func validateToken(tokenString string) (*jwt.Token, error) {
	return jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return []byte(config.Envs.JWTSecret), nil
	})
}
