package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
)

var verifyCredentialMock func(string, string) (string, error)

type authenticationServiceMock struct {
}

func newAuthenticationServiceMock() *authenticationServiceMock {
	return &authenticationServiceMock{}
}

func (asm *authenticationServiceMock) VerifyCredential(username string, password string) (string, error) {
	return verifyCredentialMock(username, password)
}

func TestPostAuthenticationHandler(t *testing.T) {
	authenticationServiceMock := newAuthenticationServiceMock()
	authenticationHandlers := NewAuthenticationHandler(authenticationServiceMock)
	method := "POST"
	urlLink := "http://localhost:8000/authentications"
	userId := "user-bbc4c22f-0129-4b14-af7a-86bbd2709b80"

	// if username is empty
	requestBody, _ := json.Marshal(map[string]string{
		"password": "jordyjordy",
	})
	req, err := http.NewRequest(method, urlLink, bytes.NewBuffer(requestBody))
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(authenticationHandlers.PostAuthenticationHandler)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("PostAuthenticationHandler returned wrong status code: got %v instead of %v", rr.Code, http.StatusBadRequest)
	}
	expected := `{"message":"Username must not be empty"}`
	if rr.Body.String() != expected {
		t.Errorf("PostAuthenticationHandler returned wrong body: got %v instead of %v", rr.Body.String(), expected)
	}

	// if password is empty
	requestBody, _ = json.Marshal(map[string]string{
		"username": "jordyf15",
	})
	req, err = http.NewRequest(method, urlLink, bytes.NewBuffer(requestBody))
	if err != nil {
		t.Fatal(err)
	}
	rr = httptest.NewRecorder()
	handler = http.HandlerFunc(authenticationHandlers.PostAuthenticationHandler)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("PostAuthenticationHandler returned wrong status code: got %v instead of %v", rr.Code, http.StatusBadRequest)
	}
	expected = `{"message":"Password must not be empty"}`
	if rr.Body.String() != expected {
		t.Errorf("PostAuthenticationHandler returned wrong body: got %v instead of %v", rr.Body.String(), expected)
	}

	// if verifycredential not found user
	verifyCredentialMock = func(s1, s2 string) (string, error) {
		return "", errors.New("username not found")
	}
	requestBody, _ = json.Marshal(map[string]string{
		"username": "jordyf15",
		"password": "jordyjordy",
	})
	req, err = http.NewRequest(method, urlLink, bytes.NewBuffer(requestBody))
	if err != nil {
		t.Fatal(err)
	}
	rr = httptest.NewRecorder()
	handler = http.HandlerFunc(authenticationHandlers.PostAuthenticationHandler)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("PostAuthenticationHandler returned wrong status code: got %v instead of %v", rr.Code, http.StatusBadRequest)
	}
	expected = `{"message":"User does not exist"}`
	if rr.Body.String() != expected {
		t.Errorf("PostAuthenticationHandler returned wrong body: got %v instead of %v", rr.Body.String(), expected)
	}

	// if verifycredential return password not match
	verifyCredentialMock = func(s1, s2 string) (string, error) {
		return "", errors.New("password not match")
	}
	requestBody, _ = json.Marshal(map[string]string{
		"username": "jordyf15",
		"password": "jordyjordy",
	})
	req, err = http.NewRequest(method, urlLink, bytes.NewBuffer(requestBody))
	if err != nil {
		t.Fatal(err)
	}
	rr = httptest.NewRecorder()
	handler = http.HandlerFunc(authenticationHandlers.PostAuthenticationHandler)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusUnauthorized {
		t.Errorf("PostAuthenticationHandler returned wrong status code: got %v instead of %v", rr.Code, http.StatusUnauthorized)
	}
	expected = `{"message":"Password is wrong"}`
	if rr.Body.String() != expected {
		t.Errorf("PostAuthenticationHandler returned wrong body: got %v instead of %v", rr.Body.String(), expected)
	}

	// if successful
	verifyCredentialMock = func(s1, s2 string) (string, error) {
		return userId, nil
	}
	requestBody, _ = json.Marshal(map[string]string{
		"username": "jordyf15",
		"password": "jordyjordy",
	})
	req, err = http.NewRequest(method, urlLink, bytes.NewBuffer(requestBody))
	if err != nil {
		t.Fatal(err)
	}
	rr = httptest.NewRecorder()
	handler = http.HandlerFunc(authenticationHandlers.PostAuthenticationHandler)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("PostAuthenticationHandler returned wrong status code: got %v instead of %v", rr.Code, http.StatusOK)
	}
}
