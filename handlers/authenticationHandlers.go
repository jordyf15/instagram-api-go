package handlers

import (
	"encoding/json"
	"instagram-go/models"
	"instagram-go/services"
	"io/ioutil"
	"net/http"
	"sync"
	"time"

	"github.com/golang-jwt/jwt"
)

type AuthenticationHandlers struct {
	sync.Mutex
	service services.IAuthenticationService
}

func NewAuthenticationHandler(service services.IAuthenticationService) *AuthenticationHandlers {
	return &AuthenticationHandlers{
		service: service,
	}
}

func (ah *AuthenticationHandlers) PostAuthenticationHandler(w http.ResponseWriter, r *http.Request) {
	bodyBytes, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		response := models.NewMessage("An error has occured in our server")
		responseBytes, err := json.Marshal(response)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(responseBytes)
		return
	}

	var credential credential
	err = json.Unmarshal(bodyBytes, &credential)
	if err != nil {
		response := models.NewMessage("An error has occured in our server")
		responseBytes, err := json.Marshal(response)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(responseBytes)
		return
	}
	var badInput bool
	errorMessage := ""
	if credential.Password == "" {
		badInput = true
		errorMessage += "Password must not be empty"
	}
	if credential.Username == "" {
		if badInput {
			errorMessage += ", "
		}
		badInput = true
		errorMessage += "Username must not be empty"
	}
	if badInput {
		response := models.NewMessage(errorMessage)
		responseBytes, err := json.Marshal(response)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}
		w.WriteHeader(http.StatusBadRequest)
		w.Write(responseBytes)
		return
	}

	ah.Lock()
	userId, err := ah.service.VerifyCredential(credential.Username, credential.Password)
	defer ah.Unlock()
	if err != nil {
		if err.Error() == "username not found" {
			response := models.NewMessage("User does not exist")
			responseBytes, err := json.Marshal(response)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(err.Error()))
				return
			}
			w.WriteHeader(http.StatusBadRequest)
			w.Write(responseBytes)
			return
		} else {
			response := models.NewMessage("Password is wrong")
			responseBytes, err := json.Marshal(response)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(err.Error()))
				return
			}
			w.WriteHeader(http.StatusUnauthorized)
			w.Write(responseBytes)
			return
		}
	}

	claims := jwt.MapClaims{}
	claims["authorized"] = true
	claims["user_id"] = userId
	claims["exp"] = time.Now().Add(time.Minute * 30).Unix()
	sign := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := sign.SignedString([]byte("secret"))
	if err != nil {
		response := models.NewMessage("An error has occured in our server")
		responseBytes, err := json.Marshal(response)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(responseBytes)
		return
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
