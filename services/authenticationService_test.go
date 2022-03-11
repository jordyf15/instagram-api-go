package services

import (
	"context"
	"errors"
	"instagram-go/models"
	"testing"
)

var findOneUserMock func(context.Context, interface{}) (*models.User, error)
var compareHashAndPasswordMock func([]byte, []byte) error

type authenticationQueryMock struct {
}

func newAuthenticationQueryMock() *authenticationQueryMock {
	return &authenticationQueryMock{}
}

func (aqm *authenticationQueryMock) findOneUser(context context.Context, filter interface{}) (*models.User, error) {
	return findOneUserMock(context, filter)
}

type authenticationVerificationMock struct {
}

func newAuthenticationVerificationMock() *authenticationVerificationMock {
	return &authenticationVerificationMock{}
}

func (avm *authenticationVerificationMock) compareHashAndPassword(hashedPassword []byte, password []byte) error {
	return compareHashAndPasswordMock(hashedPassword, password)
}

func TestVerifyCredential(t *testing.T) {
	authenticationQueryMock := newAuthenticationQueryMock()
	authenticationVerificationMock := newAuthenticationVerificationMock()
	authenticationService := NewAuthenticationService(authenticationQueryMock, authenticationVerificationMock)
	loginUserId := "user-2af4bb67-e06e-4c85-917e-80f95c140afe"
	findOneUserMock = func(ctx context.Context, i interface{}) (*models.User, error) {
		return nil, errors.New("FindOne on db returns error")
	}
	if _, err := authenticationService.VerifyCredential("jordyf15", "jordyjordy"); err == nil {
		t.Error("if FindOne on db returns error than VerifyCredential should also return error")
	}

	findOneUserMock = func(ctx context.Context, i interface{}) (*models.User, error) {
		return models.NewUser(loginUserId, "jordyf15", "jordy fer", "$2a$10$n9jX9dsnkUG3G0mjeuVHjeJI/oQ/i4YDbVCEWNHFX9TYlBku4vXoG", "jordyferdian88@gmail.com", nil), nil
	}
	compareHashAndPasswordMock = func(b1, b2 []byte) error {
		return errors.New("password is wrong")
	}
	if _, err := authenticationService.VerifyCredential("jordyf15", "jordyjordy"); err == nil {
		t.Error("If compareHashAndPassword returns error than VerifyCredential should also return error")
	}

	compareHashAndPasswordMock = func(b1, b2 []byte) error {
		return nil
	}
	authenticatedUserId, err := authenticationService.VerifyCredential("jordyf15", "jordyjordy")
	if err != nil {
		t.Error("if compareHashAndPassword and FindOne on db does not return error than VerifyCredential should also not return error")
	}
	if authenticatedUserId != loginUserId {
		t.Error("VerifyCredential should return authenticated user id")
	}

}
