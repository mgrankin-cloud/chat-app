package jwt

import (
	"Messenger-android/messenger/auth-service/sso/internal/domain/models"
	"github.com/golang-jwt/jwt/v5"
	"time"
)

func NewToken(user models.User, app models.App, duration time.Duration) (bool, string, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	//adding info into token
	claims := token.Claims.(jwt.MapClaims)
	claims["uid"] = user.ID
	claims["email"] = user.Email
	claims["username"] = user.Username
	claims["phone"] = user.Phone
	claims["exp"] = time.Now().Add(duration).Unix()
	claims["app_id"] = app.ID

	tokenString, err := token.SignedString([]byte(app.Secret))
	if err != nil {
		return false, "", err
	}

	return true, tokenString, nil
}
