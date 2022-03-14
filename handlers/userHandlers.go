package handlers

import (
	"encoding/json"
	"fmt"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"instagram-go/models"
	"instagram-go/services"
	"io"
	"io/fs"
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
	service           services.IUserService
	userHandlerHeader IUserHandlerHeader
	fileOsHandler     IFileOsHandler
}

func NewUserHandlers(service services.IUserService, userHandlerHeader IUserHandlerHeader, fileOsHandler IFileOsHandler) *UserHandlers {
	if userHandlerHeader == nil {
		userHandlerHeader = newUserHandlerHeader()
	}
	if fileOsHandler == nil {
		fileOsHandler = newFileOsHandler()
	}
	return &UserHandlers{
		service:           service,
		userHandlerHeader: userHandlerHeader,
		fileOsHandler:     fileOsHandler,
	}
}

type IUserHandlerHeader interface {
	getUserIdFromToken(string) (string, error)
}

type userHandlerHeader struct {
}

func newUserHandlerHeader() *userHandlerHeader {
	return &userHandlerHeader{}
}

func (uhh *userHandlerHeader) getUserIdFromToken(tokenString string) (string, error) {
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

type fileOsHandler struct {
}

type IFileOsHandler interface {
	decodeImage(io.Reader) (image.Image, string, error)
	resizeAndSaveFileToLocale(string, image.Image, string, string) (string, error)
	mkDirAll(string, fs.FileMode) error
}

func newFileOsHandler() *fileOsHandler {
	return &fileOsHandler{}
}

func (foh *fileOsHandler) decodeImage(r io.Reader) (image.Image, string, error) {
	return image.Decode(r)
}

func (foh *fileOsHandler) mkDirAll(path string, perm fs.FileMode) error {
	return os.MkdirAll(path, perm)
}

func (foh *fileOsHandler) resizeAndSaveFileToLocale(size string, originalProfilePicture image.Image, userId string, fileType string) (string, error) {
	var resizedProfilePicture image.Image
	fileExtension := strings.Split(fileType, "/")[1]
	resizedProfilePictureUrl := "./profile_pictures/" + size + "-profile-picture-" + userId + "." + fileExtension
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
	if fileExtension == "jpeg" {
		jpeg.Encode(resizedProfilePictureFile, resizedProfilePicture, nil)
	}
	if fileExtension == "png" {
		png.Encode(resizedProfilePictureFile, resizedProfilePicture)
	}
	if fileExtension == "gif" {
		gif.Encode(resizedProfilePictureFile, resizedProfilePicture, nil)
	}

	resizedProfilePictureFile.Close()
	return resizedProfilePictureUrl, nil
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

	var badInput bool
	errorMessage := ""
	if user.Email == "" {
		badInput = true
		errorMessage += "Email is not provided"
	}
	if user.Fullname == "" {
		if badInput {
			errorMessage += ", "
		}
		badInput = true
		errorMessage += "Full Name is not provided"
	}
	if user.Username == "" {
		if badInput {
			errorMessage += ", "
		}
		badInput = true
		errorMessage += "Username is not provided"
	}
	if user.Password == "" {
		if badInput {
			errorMessage += ", "
		}
		badInput = true
		errorMessage += "Password is not provided"
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
	isUserNameExist, err := uh.service.CheckIfUsernameExist(user.Username)
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
	if isUserNameExist {
		response := models.NewMessage("Username already exist")
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
	user.Id = "user-" + uuid.NewString()
	uh.Lock()
	err = uh.service.InsertUser(user)
	defer uh.Unlock()
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
	newpath := filepath.Join(".", "profile_pictures")
	err := uh.fileOsHandler.mkDirAll(newpath, os.ModePerm)
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
	err = nil
	parts := strings.Split(r.URL.String(), "/")
	userIdParam := parts[2]
	tokenString := r.Header.Get("Authorization")
	userIdToken, err := uh.userHandlerHeader.getUserIdFromToken(tokenString)
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
	isUserExist, err := uh.service.CheckIfUserExist(userIdParam)
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
	if !isUserExist {
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
	}
	if userIdParam != userIdToken {
		response := models.NewMessage("User is not authorized")
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
	username := r.FormValue("username")
	fullName := r.FormValue("full_name")
	password := r.FormValue("password")
	email := r.FormValue("email")
	var updatedUser models.User

	if err == nil {
		fileHeader := make([]byte, 512)
		if _, err := profilePictureFile.Read(fileHeader); err != nil {
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
		if _, err := profilePictureFile.Seek(0, 0); err != nil {
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
		if fileType := http.DetectContentType(fileHeader); fileType != "image/jpeg" && fileType != "image/png" && fileType != "image/gif" {
			response := models.NewMessage("Invalid profile picture file type")
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
		fileType := http.DetectContentType(fileHeader)

		defer profilePictureFile.Close()
		originalProfilePicture, _, err := uh.fileOsHandler.decodeImage(profilePictureFile)
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

		smallProfilePictureUrl, err := uh.fileOsHandler.resizeAndSaveFileToLocale("small", originalProfilePicture, userIdParam, fileType)
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

		averageProfilePictureUrl, err := uh.fileOsHandler.resizeAndSaveFileToLocale("average", originalProfilePicture, userIdParam, fileType)
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

		largeProfilePictureUrl, err := uh.fileOsHandler.resizeAndSaveFileToLocale("large", originalProfilePicture, userIdParam, fileType)
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
