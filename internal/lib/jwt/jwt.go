package jwt

import (
	"time"

	"github.com/golang-jwt/jwt/v5"

	"Tasks/internal/model"
)

func NewToken(user model.User, secret string, duration time.Duration) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["id"] = user.ID
	claims["login"] = user.Login
	claims["exp"] = time.Now().Add(duration).Unix()

	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
