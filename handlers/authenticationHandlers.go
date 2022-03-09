package handlers

import (
	"encoding/json"
	"fmt"
	"instagram-go/services"
	"io/ioutil"
	"net/http"
	"sync"
	"time"

	"github.com/golang-jwt/jwt"
)

type AuthenticationHandlers struct {
	sync.Mutex
	service *services.AuthenticationService
}

func NewAuthenticationHandler(service *services.AuthenticationService) *AuthenticationHandlers {
	return &AuthenticationHandlers{
		service: service,
	}
}

func (ah *AuthenticationHandlers) PostAuthenticationHandler(w http.ResponseWriter, r *http.Request) {
	bodyBytes, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
	}

	ct := r.Header.Get("content-type")
	if ct != "application/json" {
		w.WriteHeader(http.StatusUnsupportedMediaType)
		w.Write([]byte(fmt.Sprintf("need content-type 'application/json', but got '%s'", ct)))
		return
	}

	var credential credential
	err = json.Unmarshal(bodyBytes, &credential)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	ah.Lock()
	userId, err := ah.service.VerifyCredential(credential.Username, credential.Password)
	defer ah.Unlock()
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(err.Error()))
		return
	}

	claims := jwt.MapClaims{}
	claims["authorized"] = true
	claims["user_id"] = userId
	claims["exp"] = time.Now().Add(time.Minute * 30).Unix()
	sign := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := sign.SignedString([]byte("secret"))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
	}

	data := data{token}
	response := response{"User successfully authenticated", data}
	responseBytes, err := json.Marshal(response)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
	}
	w.WriteHeader(http.StatusOK)
	w.Write(responseBytes)
}

type response struct {
	Message string `json:"message"`
	Data    data   `json:"data"`
}

type data struct {
	AccessToken string `json:"access_token"`
}

type credential struct {
	Username string
	Password string
}
