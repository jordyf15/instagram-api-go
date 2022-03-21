package http

import (
	"encoding/json"
	"instagram-go/domain"
	"io/ioutil"
	"net/http"
	"strings"
)

type UserHandler struct {
	userUsecase domain.UserUsecase
}

func NewUserHandler(userUseCase domain.UserUsecase) domain.UserHandler {
	return &UserHandler{
		userUsecase: userUseCase,
	}
}

func (uh *UserHandler) PostUser(w http.ResponseWriter, r *http.Request) {
	bodyBytes, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		response := domain.NewMessage(domain.ErrInternalServerError.Error())
		responseBytes, errMarshal := json.Marshal(response)
		if errMarshal != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(errMarshal.Error()))
			return
		}
		w.WriteHeader(userGetStatusCode(domain.ErrInternalServerError))
		w.Write(responseBytes)
		return
	}

	var user domain.User
	err = json.Unmarshal(bodyBytes, &user)
	if err != nil {
		response := domain.NewMessage(domain.ErrInternalServerError.Error())
		responseBytes, errMarshal := json.Marshal(response)
		if errMarshal != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(errMarshal.Error()))
			return
		}
		w.WriteHeader(userGetStatusCode(domain.ErrInternalServerError))
		w.Write(responseBytes)
		return
	}
	if user.Email == "" {
		response := domain.NewMessage(domain.ErrMissingEmailInput.Error())
		responseBytes, errMarshal := json.Marshal(response)
		if errMarshal != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(errMarshal.Error()))
			return
		}
		w.WriteHeader(userGetStatusCode(domain.ErrMissingEmailInput))
		w.Write(responseBytes)
		return
	}
	if user.Fullname == "" {
		response := domain.NewMessage(domain.ErrMissingFullNameInput.Error())
		responseBytes, errMarshal := json.Marshal(response)
		if errMarshal != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(errMarshal.Error()))
			return
		}
		w.WriteHeader(userGetStatusCode(domain.ErrMissingFullNameInput))
		w.Write(responseBytes)
		return
	}
	if user.Username == "" {
		response := domain.NewMessage(domain.ErrMissingUsernameInput.Error())
		responseBytes, errMarshal := json.Marshal(response)
		if errMarshal != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(errMarshal.Error()))
			return
		}
		w.WriteHeader(userGetStatusCode(domain.ErrMissingUsernameInput))
		w.Write(responseBytes)
		return
	}

	if user.Password == "" {
		response := domain.NewMessage(domain.ErrMissingPasswordInput.Error())
		responseBytes, errMarshal := json.Marshal(response)
		if errMarshal != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(errMarshal.Error()))
			return
		}
		w.WriteHeader(userGetStatusCode(domain.ErrMissingPasswordInput))
		w.Write(responseBytes)
		return
	}

	err = uh.userUsecase.InsertUser(&user)
	if err != nil {
		response := domain.NewMessage(err.Error())
		responseBytes, errMarshal := json.Marshal(response)
		if errMarshal != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(errMarshal.Error()))
			return
		}
		w.WriteHeader(userGetStatusCode(err))
		w.Write(responseBytes)
		return
	} else {
		response := domain.NewMessage("User successfully registered")
		responseBytes, errMarshal := json.Marshal(response)
		if errMarshal != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(errMarshal.Error()))
			return
		}
		w.WriteHeader(http.StatusCreated)
		w.Write(responseBytes)
		return
	}
}

func (uh *UserHandler) PutUser(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(r.URL.String(), "/")
	userIdParam := parts[2]
	tokenString := r.Header.Get("Authorization")
	r.ParseMultipartForm(10 << 20)
	profilePictureFile, _, _ := r.FormFile("profile_picture")
	username := r.FormValue("username")
	fullName := r.FormValue("full_name")
	password := r.FormValue("password")
	email := r.FormValue("email")
	var updatedUser domain.User
	updatedUser.Id = userIdParam
	updatedUser.Username = username
	updatedUser.Fullname = fullName
	updatedUser.Password = password
	updatedUser.Email = email
	if profilePictureFile != nil {
		fileHeader := make([]byte, 512)
		if _, err := profilePictureFile.Read(fileHeader); err != nil {
			response := domain.NewMessage(domain.ErrInternalServerError.Error())
			responseBytes, errMarshal := json.Marshal(response)
			if errMarshal != nil {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(errMarshal.Error()))
				return
			}
			w.WriteHeader(userGetStatusCode(domain.ErrInternalServerError))
			w.Write(responseBytes)
			return
		}
		if _, err := profilePictureFile.Seek(0, 0); err != nil {
			response := domain.NewMessage(domain.ErrInternalServerError.Error())
			responseBytes, errMarshal := json.Marshal(response)
			if errMarshal != nil {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(errMarshal.Error()))
				return
			}
			w.WriteHeader(userGetStatusCode(domain.ErrInternalServerError))
			w.Write(responseBytes)
			return
		}
		if fileType := http.DetectContentType(fileHeader); fileType != "image/jpeg" && fileType != "image/png" && fileType != "image/gif" {
			response := domain.NewMessage(domain.ErrInvalidProfilePicture.Error())
			responseBytes, errMarshal := json.Marshal(response)
			if errMarshal != nil {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(errMarshal.Error()))
				return
			}
			w.WriteHeader(userGetStatusCode(domain.ErrInvalidProfilePicture))
			w.Write(responseBytes)
			return
		}
	}

	err := uh.userUsecase.UpdateUser(&updatedUser, tokenString, profilePictureFile)
	if err != nil {
		response := domain.NewMessage(err.Error())
		responseBytes, errMarshal := json.Marshal(response)
		if errMarshal != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(errMarshal.Error()))
			return
		}
		w.WriteHeader(userGetStatusCode(err))
		w.Write(responseBytes)
		return
	}
	response := domain.NewMessage("User successfully Updated")
	responseBytes, errMarshal := json.Marshal(response)
	if errMarshal != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(errMarshal.Error()))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(responseBytes)
}

func (uh *UserHandler) AuthenticateUser(w http.ResponseWriter, r *http.Request) {
	bodyBytes, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		response := domain.NewMessage(domain.ErrInternalServerError.Error())
		responseBytes, errMarshal := json.Marshal(response)
		if errMarshal != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(errMarshal.Error()))
			return
		}
		w.WriteHeader(userGetStatusCode(domain.ErrInternalServerError))
		w.Write(responseBytes)
		return
	}
	var user domain.User
	err = json.Unmarshal(bodyBytes, &user)
	if err != nil {
		response := domain.NewMessage(domain.ErrInternalServerError.Error())
		responseBytes, errMarshal := json.Marshal(response)
		if errMarshal != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(errMarshal.Error()))
			return
		}
		w.WriteHeader(userGetStatusCode(domain.ErrInternalServerError))
		w.Write(responseBytes)
		return
	}
	if user.Username == "" {
		response := domain.NewMessage(domain.ErrMissingUsernameInput.Error())
		responseBytes, errMarshal := json.Marshal(response)
		if errMarshal != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(errMarshal.Error()))
			return
		}
		w.WriteHeader(userGetStatusCode(domain.ErrMissingUsernameInput))
		w.Write(responseBytes)
		return
	}
	if user.Password == "" {
		response := domain.NewMessage(domain.ErrMissingPasswordInput.Error())
		responseBytes, errMarshal := json.Marshal(response)
		if errMarshal != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(errMarshal.Error()))
			return
		}
		w.WriteHeader(userGetStatusCode(domain.ErrMissingPasswordInput))
		w.Write(responseBytes)
		return
	}

	accessToken, err := uh.userUsecase.VerifyCredential(user.Username, user.Password)
	if err != nil {
		response := domain.NewMessage(err.Error())
		responseBytes, errMarshal := json.Marshal(response)
		if errMarshal != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(errMarshal.Error()))
			return
		}
		w.WriteHeader(userGetStatusCode(err))
		w.Write(responseBytes)
		return
	}
	data := domain.NewDataAuthentication(accessToken)
	response := domain.NewDataResponseAuthentication("User successfully authenticated", *data)
	responseBytes, errMarshal := json.Marshal(response)
	if errMarshal != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(errMarshal.Error()))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(responseBytes)
}

func userGetStatusCode(err error) int {
	switch err {
	case domain.ErrInternalServerError:
		return http.StatusInternalServerError
	case domain.ErrMissingEmailInput, domain.ErrMissingFullNameInput, domain.ErrMissingUsernameInput, domain.ErrMissingPasswordInput, domain.ErrInvalidProfilePicture:
		return http.StatusBadRequest
	case domain.ErrUsernameConflict:
		return http.StatusConflict
	case domain.ErrUserNotFound:
		return http.StatusNotFound
	case domain.ErrUnauthorizedUserUpdate:
		return http.StatusUnauthorized
	case domain.ErrPasswordWrong:
		return http.StatusForbidden
	}
	return http.StatusOK
}
