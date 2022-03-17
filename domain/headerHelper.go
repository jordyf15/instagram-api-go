package domain

import (
	"fmt"

	"github.com/golang-jwt/jwt"
)

type IHeaderHelper interface {
	GetUserIdFromToken(string) (string, error)
}

type HeaderHelper struct {
}

func NewHeaderHelper() *HeaderHelper {
	return &HeaderHelper{}
}

func (hh *HeaderHelper) GetUserIdFromToken(tokenString string) (string, error) {
	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		return []byte("secret"), nil
	})
	if err != nil {
		return "", err
	}
	claims := token.Claims.(jwt.MapClaims)
	userId := fmt.Sprintf("%v", claims["user_id"])
	return userId, nil
}
