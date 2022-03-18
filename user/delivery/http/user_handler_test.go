package http_test

import (
	"bytes"
	"encoding/json"
	"instagram-go/domain"
	"instagram-go/domain/mocks"
	userHttp "instagram-go/user/delivery/http"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

func TestPostUserSuite(t *testing.T) {
	suite.Run(t, new(PostUserSuite))
}

func TestPutUserSuite(t *testing.T) {
	suite.Run(t, new(PutUserSuite))
}
func TestAuthenticateUserSuite(t *testing.T) {
	suite.Run(t, new(AuthenticateUserSuite))
}

type PostUserSuite struct {
	suite.Suite
	userUsecase *mocks.UserUsecase
}

func (pus *PostUserSuite) SetupTest() {
	pus.userUsecase = new(mocks.UserUsecase)
}

func (pus *PostUserSuite) TestEmailNotProvided() {
	requestBody, _ := json.Marshal(map[string]string{
		"fullname": "jordy ferdian",
		"username": "jordyf15",
		"email":    "",
		"password": "jordyjordy",
	})

	req, _ := http.NewRequest("POST", "/users", bytes.NewBuffer(requestBody))
	rr := httptest.NewRecorder()
	userHandler := userHttp.NewUserHandler(pus.userUsecase)
	handler := http.HandlerFunc(userHandler.PostUser)
	handler.ServeHTTP(rr, req)

	assert.Equalf(pus.T(), http.StatusBadRequest, rr.Code, "Should have responded with http status code %s but got %s", http.StatusBadRequest, rr.Code)
	expectedBody := `{"message":"email must not be empty"}`
	assert.Equalf(pus.T(), expectedBody, rr.Body.String(), "Should have responded with body %s but got %s", expectedBody, rr.Body.String())
}

func (pus *PostUserSuite) TestFullnameIsNotProvided() {
	requestBody, _ := json.Marshal(map[string]string{
		"email":    "jordyferdian@gmail.com",
		"username": "jordyf15",
		"fullname": "",
		"password": "jordyjordy",
	})
	req, _ := http.NewRequest("POST", "/users", bytes.NewBuffer(requestBody))
	rr := httptest.NewRecorder()
	userHandler := userHttp.NewUserHandler(pus.userUsecase)
	handler := http.HandlerFunc(userHandler.PostUser)
	handler.ServeHTTP(rr, req)

	assert.Equalf(pus.T(), http.StatusBadRequest, rr.Code, "Should have responded with http status code %s but got %s", http.StatusBadRequest, rr.Code)
	expectedBody := `{"message":"full name must not be empty"}`
	assert.Equalf(pus.T(), expectedBody, rr.Body.String(), "Should have responded with body %s but got %s", expectedBody, rr.Body.String())
}

func (pus *PostUserSuite) TestUsernameIsNotProvided() {
	requestBody, _ := json.Marshal(map[string]string{
		"email":    "jordyferdian@gmail.com",
		"fullname": "jordy ferdian",
		"username": "",
		"password": "jordyjordy",
	})
	req, _ := http.NewRequest("POST", "/users", bytes.NewBuffer(requestBody))
	rr := httptest.NewRecorder()
	userHandler := userHttp.NewUserHandler(pus.userUsecase)
	handler := http.HandlerFunc(userHandler.PostUser)
	handler.ServeHTTP(rr, req)

	assert.Equalf(pus.T(), http.StatusBadRequest, rr.Code, "Should have responded with http status code %s but got %s", http.StatusBadRequest, rr.Code)
	expectedBody := `{"message":"username must not be empty"}`
	assert.Equalf(pus.T(), expectedBody, rr.Body.String(), "Should have responded with body %s but got %s", expectedBody, rr.Body.String())
}

func (pus *PostUserSuite) TestPasswordIsNotProvided() {
	requestBody, _ := json.Marshal(map[string]string{
		"email":    "jordyferdian@gmail.com",
		"fullname": "jordy ferdian",
		"username": "jordyf15",
		"password": "",
	})
	req, _ := http.NewRequest("POST", "/users", bytes.NewBuffer(requestBody))
	rr := httptest.NewRecorder()
	userHandler := userHttp.NewUserHandler(pus.userUsecase)
	handler := http.HandlerFunc(userHandler.PostUser)
	handler.ServeHTTP(rr, req)

	assert.Equalf(pus.T(), http.StatusBadRequest, rr.Code, "Should have responded with http status code %s but got %s", http.StatusBadRequest, rr.Code)
	expectedBody := `{"message":"password must not be empty"}`
	assert.Equalf(pus.T(), expectedBody, rr.Body.String(), "Should have responded with body %s but got %s", expectedBody, rr.Body.String())
}

func (pus *PostUserSuite) TestInsertUserError() {
	requestBody, _ := json.Marshal(map[string]string{
		"email":    "jordyferdian@gmail.com",
		"fullname": "jordy ferdian",
		"username": "jordyf15",
		"password": "jordyjordy",
	})
	pus.userUsecase.On("InsertUser", mock.Anything).Return(domain.ErrInternalServerError)
	req, _ := http.NewRequest("POST", "/users", bytes.NewBuffer(requestBody))
	rr := httptest.NewRecorder()
	userHandler := userHttp.NewUserHandler(pus.userUsecase)
	handler := http.HandlerFunc(userHandler.PostUser)
	handler.ServeHTTP(rr, req)

	assert.Equalf(pus.T(), http.StatusInternalServerError, rr.Code, "Should have responded with http status code %s but got %s", http.StatusInternalServerError, rr.Code)
	expectedBody := `{"message":"an error has occured in our server"}`
	assert.Equalf(pus.T(), expectedBody, rr.Body.String(), "Should have responded with body %s but got %s", expectedBody, rr.Body.String())
}

func (pus *PostUserSuite) PostUserSuccessful() {
	requestBody, _ := json.Marshal(map[string]string{
		"email":    "jordyferdian@gmail.com",
		"fullname": "jordy ferdian",
		"username": "jordyf15",
		"password": "jordyjordy",
	})
	pus.userUsecase.On("InsertUser", mock.Anything).Return(nil)
	req, _ := http.NewRequest("POST", "/users", bytes.NewBuffer(requestBody))
	rr := httptest.NewRecorder()
	userHandler := userHttp.NewUserHandler(pus.userUsecase)
	handler := http.HandlerFunc(userHandler.PostUser)
	handler.ServeHTTP(rr, req)

	assert.Equalf(pus.T(), http.StatusCreated, rr.Code, "Should have responded with http status code %s but got %s", http.StatusCreated, rr.Code)
	expectedBody := `{"message":"User successfully registered"}`
	assert.Equalf(pus.T(), expectedBody, rr.Body.String(), "Should have responded with body %s but got %s", expectedBody, rr.Body.String())
}

type PutUserSuite struct {
	suite.Suite
	userUsecase *mocks.UserUsecase
}

func (pus *PutUserSuite) SetupTest() {
	pus.userUsecase = new(mocks.UserUsecase)
}

func (pus *PutUserSuite) TestInvalidProfilePicture() {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	fw, _ := writer.CreateFormField("username")
	_, _ = io.Copy(fw, strings.NewReader("jordyf88"))
	fw, _ = writer.CreateFormField("full_name")
	_, _ = io.Copy(fw, strings.NewReader("jordy feru"))
	fw, _ = writer.CreateFormField("password")
	_, _ = io.Copy(fw, strings.NewReader("jorjor123"))
	fw, _ = writer.CreateFormField("email")
	_, _ = io.Copy(fw, strings.NewReader("jordyjordy@gmail.com"))
	fw, _ = writer.CreateFormFile("profile_picture", "bmp.bmp")
	file, _ := os.Open("./test_profile_pictures/bmp.bmp")
	_, _ = io.Copy(fw, file)
	writer.Close()
	req, _ := http.NewRequest("PUT", "/users/userid1", bytes.NewReader(body.Bytes()))
	req.Header.Set("Content-Type", writer.FormDataContentType())
	rr := httptest.NewRecorder()
	userHandler := userHttp.NewUserHandler(pus.userUsecase)
	handler := http.HandlerFunc(userHandler.PutUser)
	handler.ServeHTTP(rr, req)

	assert.Equalf(pus.T(), http.StatusBadRequest, rr.Code, "Should have responded with http status code %s but got %s", http.StatusBadRequest, rr.Code)
	expectedMessage := `{"message":"` + domain.ErrInvalidProfilePicture.Error() + `"}`
	assert.Equalf(pus.T(), expectedMessage, rr.Body.String(), "Should have responded with body %s but got %s", expectedMessage, rr.Body.String())
}

func (pus *PutUserSuite) TestUpdateUserError() {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	fw, _ := writer.CreateFormField("username")
	_, _ = io.Copy(fw, strings.NewReader("jordyf88"))
	fw, _ = writer.CreateFormField("full_name")
	_, _ = io.Copy(fw, strings.NewReader("jordy feru"))
	fw, _ = writer.CreateFormField("password")
	_, _ = io.Copy(fw, strings.NewReader("jorjor123"))
	fw, _ = writer.CreateFormField("email")
	_, _ = io.Copy(fw, strings.NewReader("jordyjordy@gmail.com"))
	writer.Close()
	pus.userUsecase.On("UpdateUser", mock.Anything, mock.Anything, mock.Anything).Return(domain.ErrInternalServerError)
	req, _ := http.NewRequest("PUT", "/users/userid1", bytes.NewReader(body.Bytes()))
	req.Header.Set("Content-Type", writer.FormDataContentType())
	rr := httptest.NewRecorder()
	userHandler := userHttp.NewUserHandler(pus.userUsecase)
	handler := http.HandlerFunc(userHandler.PutUser)
	handler.ServeHTTP(rr, req)

	assert.Equalf(pus.T(), http.StatusInternalServerError, rr.Code, "Should have responded with http status code %s but got %s", http.StatusInternalServerError, rr.Code)
	expectedMessage := `{"message":"` + domain.ErrInternalServerError.Error() + `"}`
	assert.Equalf(pus.T(), expectedMessage, rr.Body.String(), "Should have responded with body %s but got %s", expectedMessage, rr.Body.String())
}

func (pus *PutUserSuite) TestPutUserSuccessful() {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	fw, _ := writer.CreateFormField("username")
	_, _ = io.Copy(fw, strings.NewReader("jordyf88"))
	fw, _ = writer.CreateFormField("full_name")
	_, _ = io.Copy(fw, strings.NewReader("jordy feru"))
	fw, _ = writer.CreateFormField("password")
	_, _ = io.Copy(fw, strings.NewReader("jorjor123"))
	fw, _ = writer.CreateFormField("email")
	_, _ = io.Copy(fw, strings.NewReader("jordyjordy@gmail.com"))
	writer.Close()
	pus.userUsecase.On("UpdateUser", mock.Anything, mock.Anything, mock.Anything).Return(nil)
	req, _ := http.NewRequest("PUT", "/users/userid1", bytes.NewReader(body.Bytes()))
	req.Header.Set("Content-Type", writer.FormDataContentType())
	rr := httptest.NewRecorder()
	userHandler := userHttp.NewUserHandler(pus.userUsecase)
	handler := http.HandlerFunc(userHandler.PutUser)
	handler.ServeHTTP(rr, req)

	assert.Equalf(pus.T(), http.StatusOK, rr.Code, "Should have responded with http status code %s but got %s", http.StatusOK, rr.Code)
	expectedMessage := `{"message":"User successfully Updated"}`
	assert.Equalf(pus.T(), expectedMessage, rr.Body.String(), "Should have responded with body %s but got %s", expectedMessage, rr.Body.String())
}

type AuthenticateUserSuite struct {
	suite.Suite
	userUsecase *mocks.UserUsecase
}

func (aus *AuthenticateUserSuite) SetupTest() {
	aus.userUsecase = new(mocks.UserUsecase)
}

func (aus *AuthenticateUserSuite) TestUsernameNotProvided() {
	requestBody, _ := json.Marshal(map[string]string{
		"username": "",
		"password": "jordyjordy",
	})
	req, _ := http.NewRequest("POST", "/authentications", bytes.NewBuffer(requestBody))
	rr := httptest.NewRecorder()
	userHandler := userHttp.NewUserHandler(aus.userUsecase)
	handler := http.HandlerFunc(userHandler.AuthenticateUser)
	handler.ServeHTTP(rr, req)

	assert.Equalf(aus.T(), http.StatusBadRequest, rr.Code, "Should have responded with http status code %s but got %s", http.StatusBadRequest, rr.Code)
	expectedBody := `{"message":"username must not be empty"}`
	assert.Equalf(aus.T(), expectedBody, rr.Body.String(), "Should have responded with body %s but got %s", expectedBody, rr.Body.String())
}

func (aus *AuthenticateUserSuite) TestPasswordNotProvided() {
	requestBody, _ := json.Marshal(map[string]string{
		"username": "jordyf15",
		"password": "",
	})
	req, _ := http.NewRequest("POST", "/authentications", bytes.NewBuffer(requestBody))
	rr := httptest.NewRecorder()
	userHandler := userHttp.NewUserHandler(aus.userUsecase)
	handler := http.HandlerFunc(userHandler.AuthenticateUser)
	handler.ServeHTTP(rr, req)

	assert.Equalf(aus.T(), http.StatusBadRequest, rr.Code, "Should have responded with http status code %s but got %s", http.StatusBadRequest, rr.Code)
	expectedBody := `{"message":"password must not be empty"}`
	assert.Equalf(aus.T(), expectedBody, rr.Body.String(), "Should have responded with body %s but got %s", expectedBody, rr.Body.String())
}

func (aus *AuthenticateUserSuite) TestVerifyCredentialError() {
	requestBody, _ := json.Marshal(map[string]string{
		"username": "jordyf15",
		"password": "jordyjordy",
	})
	aus.userUsecase.On("VerifyCredential", mock.Anything, mock.Anything).Return("", domain.ErrInternalServerError)
	req, _ := http.NewRequest("POST", "/authentications", bytes.NewBuffer(requestBody))
	rr := httptest.NewRecorder()
	userHandler := userHttp.NewUserHandler(aus.userUsecase)
	handler := http.HandlerFunc(userHandler.AuthenticateUser)
	handler.ServeHTTP(rr, req)

	assert.Equalf(aus.T(), http.StatusInternalServerError, rr.Code, "Should have responded with http status code %s but got %s", http.StatusInternalServerError, rr.Code)
	expectedBody := `{"message":"` + domain.ErrInternalServerError.Error() + `"}`
	assert.Equalf(aus.T(), expectedBody, rr.Body.String(), "Should have responded with body %s but got %s", expectedBody, rr.Body.String())
}

func (aus *AuthenticateUserSuite) TestAuthenticateUserSuccessful() {
	requestBody, _ := json.Marshal(map[string]string{
		"username": "jordyf15",
		"password": "jordyjordy",
	})
	aus.userUsecase.On("VerifyCredential", mock.Anything, mock.Anything).Return("token", nil)
	req, _ := http.NewRequest("POST", "/authentications", bytes.NewBuffer(requestBody))
	rr := httptest.NewRecorder()
	userHandler := userHttp.NewUserHandler(aus.userUsecase)
	handler := http.HandlerFunc(userHandler.AuthenticateUser)
	handler.ServeHTTP(rr, req)

	assert.Equalf(aus.T(), http.StatusOK, rr.Code, "Should have responded with http status code %s but got %s", http.StatusOK, rr.Code)
	expectedBody := `{"message":"User successfully authenticated","data":{"access_token":"token"}}`
	assert.Equalf(aus.T(), expectedBody, rr.Body.String(), "Should have responded with body %s but got %s", expectedBody, rr.Body.String())
}
