package domain

import "golang.org/x/crypto/bcrypt"

type IAuthenticationHelper interface {
	CompareHashAndPassword([]byte, []byte) error
}

type AuthenticationHelper struct {
}

func NewAuthenticationHelper() *AuthenticationHelper {
	return &AuthenticationHelper{}
}

func (ah *AuthenticationHelper) CompareHashAndPassword(hashedPassword []byte, password []byte) error {
	return bcrypt.CompareHashAndPassword(hashedPassword, password)
}
