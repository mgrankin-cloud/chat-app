package jwt

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"log"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/mgrankin-cloud/messenger/internal/domain/models"
)

var secretKey = generateRandomKey(32)

func generateRandomKey(len int) string {
	key := make([]byte, len)
	_, err := rand.Read(key)
	if err != nil {
		log.Fatalf("Ошибка генерации секретного ключа: %v", err)
	}
	return base64.StdEncoding.EncodeToString(key)
}

func NewToken(user models.User, app models.App, duration time.Duration) (bool, string, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["uid"] = user.ID
	claims["email"] = user.Email
	claims["role"] = user.Role
	claims["exp"] = time.Now().Add(duration).Unix()
	claims["app_id"] = app.ID

	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return false, "", err
	}

	return true, tokenString, nil
}

func validateToken(tokenStr string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Неожиданный метод подписи: %v", token.Header["alg"])
		}
		return []byte(secretKey), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("Неверный токен")
}
