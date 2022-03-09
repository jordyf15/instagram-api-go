package handlers

import (
	"encoding/json"
	"fmt"
	"image"
	"image/jpeg"
	"instagram-go/models"
	"instagram-go/services"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"github.com/nfnt/resize"
)

type UserHandlers struct {
	sync.Mutex
	service services.UserService
}

func NewUserHandlers(service services.UserService) *UserHandlers {
	return &UserHandlers{
		service: service,
	}
}

func (uh *UserHandlers) PostUserHandler(w http.ResponseWriter, r *http.Request) {
	bodyBytes, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	ct := r.Header.Get("content-type")
	if ct != r.Header.Get("content-type") {
		w.WriteHeader(http.StatusUnsupportedMediaType)
		w.Write([]byte(fmt.Sprintf("need content-type 'application/json', but got '%s'", ct)))
		return
	}

	var user models.User
	err = json.Unmarshal(bodyBytes, &user)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	user.Id = "user-" + uuid.NewString()
	uh.Lock()
	err = uh.service.InsertUser(user)
	defer uh.Unlock()

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	} else {
		response := models.NewMessage("User successfully registered")
		responseBytes, err := json.Marshal(response)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}
		w.WriteHeader(http.StatusCreated)
		w.Write(responseBytes)
	}
}

func (uh *UserHandlers) PutUserHandler(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(r.URL.String(), "/")
	userIdParam := parts[2]
	tokenString := r.Header.Get("Authorization")
	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		if jwt.GetSigningMethod("HS256") != t.Method {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return []byte("secret"), nil
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	claims := token.Claims.(jwt.MapClaims)
	userIdToken := claims["user_id"]
	if userIdParam != userIdToken {
		response := models.NewMessage("Not authorized")
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

	r.ParseMultipartForm(10 << 20)
	profilePictureFile, _, err := r.FormFile("profile_picture")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	username := r.FormValue("username")
	fullName := r.FormValue("full_name")
	password := r.FormValue("password")
	email := r.FormValue("email")

	var updatedUser models.User
	newpath := filepath.Join(".", "profile_pictures")
	err = os.MkdirAll(newpath, os.ModePerm)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	if err == nil {
		defer profilePictureFile.Close()
		originalProfilePicture, _, err := image.Decode(profilePictureFile)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}

		smallProfilePictureUrl, err := saveFileToLocale("small", originalProfilePicture, userIdParam)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}
		averageProfilePictureUrl, err := saveFileToLocale("average", originalProfilePicture, userIdParam)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}
		largeProfilePictureUrl, err := saveFileToLocale("large", originalProfilePicture, userIdParam)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}

		smallProfilePicture := models.NewProfilePicture("small", "150 x 150 px", smallProfilePictureUrl)
		averageProfilePicture := models.NewProfilePicture("average", "400 x 400 px", averageProfilePictureUrl)
		largeProfilePicture := models.NewProfilePicture("large", "800 x 800 px", largeProfilePictureUrl)

		updatedUser = *models.NewUser(userIdParam, username, fullName, password,
			email, []models.ProfilePicture{*smallProfilePicture, *averageProfilePicture, *largeProfilePicture})
	} else {
		updatedUser = *models.NewUser(userIdParam, username, fullName, password, email, nil)
	}
	uh.Lock()
	err = uh.service.UpdateUser(updatedUser)
	defer uh.Unlock()

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	} else {
		response := *models.NewMessage("User successfully Updated")
		responseBytes, err := json.Marshal(response)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
		}
		w.WriteHeader(http.StatusOK)
		w.Write(responseBytes)
	}
}

func saveFileToLocale(size string, originalProfilePicture image.Image, userId string) (string, error) {
	var resizedProfilePicture image.Image
	resizedProfilePictureUrl := "./profile_pictures/" + size + "-profile-picture-" + userId + ".jpeg"
	switch size {
	case "small":
		resizedProfilePicture = resize.Resize(150, 150, originalProfilePicture, resize.Lanczos3)
	case "average":
		resizedProfilePicture = resize.Resize(400, 400, originalProfilePicture, resize.Lanczos3)
	case "large":
		resizedProfilePicture = resize.Resize(800, 800, originalProfilePicture, resize.Lanczos3)
	}

	resizedProfilePictureFile, err := os.Create(resizedProfilePictureUrl)
	if err != nil {
		return resizedProfilePictureUrl, err
	}
	jpeg.Encode(resizedProfilePictureFile, resizedProfilePicture, nil)
	resizedProfilePictureFile.Close()
	return resizedProfilePictureUrl, nil
}
