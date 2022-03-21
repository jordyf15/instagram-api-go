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

func TestUserHandlerSuite(t *testing.T) {
	suite.Run(t, new(UserHandlerSuite))
}

type UserHandlerSuite struct {
	suite.Suite
	userUsecase *mocks.UserUsecase
}

func (uh *UserHandlerSuite) SetupTest() {
	uh.userUsecase = new(mocks.UserUsecase)
}

func (uh *UserHandlerSuite) TestPostUserEmailNotProvided() {
	requestBody, _ := json.Marshal(map[string]string{
		"fullname": "jordy ferdian",
		"username": "jordyf15",
		"email":    "",
		"password": "jordyjordy",
	})

	req, _ := http.NewRequest("POST", "/users", bytes.NewBuffer(requestBody))
	rr := httptest.NewRecorder()
	userHandler := userHttp.NewUserHandler(uh.userUsecase)
	handler := http.HandlerFunc(userHandler.PostUser)
	handler.ServeHTTP(rr, req)

	assert.Equalf(uh.T(), http.StatusBadRequest, rr.Code, "Should have responded with http status code %v but got %v", http.StatusBadRequest, rr.Code)
	expectedBody := `{"message":"email must not be empty"}`
	assert.Equalf(uh.T(), expectedBody, rr.Body.String(), "Should have responded with body %s but got %s", expectedBody, rr.Body.String())
}

func (uh *UserHandlerSuite) TestPostUserFullnameIsNotProvided() {
	requestBody, _ := json.Marshal(map[string]string{
		"email":    "jordyferdian@gmail.com",
		"username": "jordyf15",
		"fullname": "",
		"password": "jordyjordy",
	})
	req, _ := http.NewRequest("POST", "/users", bytes.NewBuffer(requestBody))
	rr := httptest.NewRecorder()
	userHandler := userHttp.NewUserHandler(uh.userUsecase)
	handler := http.HandlerFunc(userHandler.PostUser)
	handler.ServeHTTP(rr, req)

	assert.Equalf(uh.T(), http.StatusBadRequest, rr.Code, "Should have responded with http status code %v but got %v", http.StatusBadRequest, rr.Code)
	expectedBody := `{"message":"full name must not be empty"}`
	assert.Equalf(uh.T(), expectedBody, rr.Body.String(), "Should have responded with body %s but got %s", expectedBody, rr.Body.String())
}

func (uh *UserHandlerSuite) TestPostUserUsernameIsNotProvided() {
	requestBody, _ := json.Marshal(map[string]string{
		"email":    "jordyferdian@gmail.com",
		"fullname": "jordy ferdian",
		"username": "",
		"password": "jordyjordy",
	})
	req, _ := http.NewRequest("POST", "/users", bytes.NewBuffer(requestBody))
	rr := httptest.NewRecorder()
	userHandler := userHttp.NewUserHandler(uh.userUsecase)
	handler := http.HandlerFunc(userHandler.PostUser)
	handler.ServeHTTP(rr, req)

	assert.Equalf(uh.T(), http.StatusBadRequest, rr.Code, "Should have responded with http status code %v but got %v", http.StatusBadRequest, rr.Code)
	expectedBody := `{"message":"username must not be empty"}`
	assert.Equalf(uh.T(), expectedBody, rr.Body.String(), "Should have responded with body %s but got %s", expectedBody, rr.Body.String())
}

func (uh *UserHandlerSuite) TestPostUserPasswordIsNotProvided() {
	requestBody, _ := json.Marshal(map[string]string{
		"email":    "jordyferdian@gmail.com",
		"fullname": "jordy ferdian",
		"username": "jordyf15",
		"password": "",
	})
	req, _ := http.NewRequest("POST", "/users", bytes.NewBuffer(requestBody))
	rr := httptest.NewRecorder()
	userHandler := userHttp.NewUserHandler(uh.userUsecase)
	handler := http.HandlerFunc(userHandler.PostUser)
	handler.ServeHTTP(rr, req)

	assert.Equalf(uh.T(), http.StatusBadRequest, rr.Code, "Should have responded with http status code %v but got %v", http.StatusBadRequest, rr.Code)
	expectedBody := `{"message":"password must not be empty"}`
	assert.Equalf(uh.T(), expectedBody, rr.Body.String(), "Should have responded with body %s but got %s", expectedBody, rr.Body.String())
}

func (uh *UserHandlerSuite) TestPostUserInsertUserError() {
	requestBody, _ := json.Marshal(map[string]string{
		"email":    "jordyferdian@gmail.com",
		"fullname": "jordy ferdian",
		"username": "jordyf15",
		"password": "jordyjordy",
	})
	uh.userUsecase.On("InsertUser", mock.AnythingOfType("*domain.User")).Return(domain.ErrInternalServerError)
	req, _ := http.NewRequest("POST", "/users", bytes.NewBuffer(requestBody))
	rr := httptest.NewRecorder()
	userHandler := userHttp.NewUserHandler(uh.userUsecase)
	handler := http.HandlerFunc(userHandler.PostUser)
	handler.ServeHTTP(rr, req)

	assert.Equalf(uh.T(), http.StatusInternalServerError, rr.Code, "Should have responded with http status code %v but got %v", http.StatusInternalServerError, rr.Code)
	expectedBody := `{"message":"an error has occured in our server"}`
	assert.Equalf(uh.T(), expectedBody, rr.Body.String(), "Should have responded with body %s but got %s", expectedBody, rr.Body.String())
}

func (uh *UserHandlerSuite) TestPostUserSuccessful() {
	requestBody, _ := json.Marshal(map[string]string{
		"email":    "jordyferdian@gmail.com",
		"fullname": "jordy ferdian",
		"username": "jordyf15",
		"password": "jordyjordy",
	})
	uh.userUsecase.On("InsertUser", mock.AnythingOfType("*domain.User")).Return(nil)
	req, _ := http.NewRequest("POST", "/users", bytes.NewBuffer(requestBody))
	rr := httptest.NewRecorder()
	userHandler := userHttp.NewUserHandler(uh.userUsecase)
	handler := http.HandlerFunc(userHandler.PostUser)
	handler.ServeHTTP(rr, req)

	assert.Equalf(uh.T(), http.StatusCreated, rr.Code, "Should have responded with http status code %v but got %v", http.StatusCreated, rr.Code)
	expectedBody := `{"message":"User successfully registered"}`
	assert.Equalf(uh.T(), expectedBody, rr.Body.String(), "Should have responded with body %s but got %s", expectedBody, rr.Body.String())
}

func (uh *UserHandlerSuite) TestPutUserInvalidProfilePicture() {
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
	userHandler := userHttp.NewUserHandler(uh.userUsecase)
	handler := http.HandlerFunc(userHandler.PutUser)
	handler.ServeHTTP(rr, req)

	assert.Equalf(uh.T(), http.StatusBadRequest, rr.Code, "Should have responded with http status code %v but got %v", http.StatusBadRequest, rr.Code)
	expectedMessage := `{"message":"` + domain.ErrInvalidProfilePicture.Error() + `"}`
	assert.Equalf(uh.T(), expectedMessage, rr.Body.String(), "Should have responded with body %s but got %s", expectedMessage, rr.Body.String())
}

func (uh *UserHandlerSuite) TestUpdateUserError() {
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
	uh.userUsecase.On("UpdateUser", mock.AnythingOfType("*domain.User"), mock.AnythingOfType("string"), mock.Anything).Return(domain.ErrInternalServerError)
	req, _ := http.NewRequest("PUT", "/users/userid1", bytes.NewReader(body.Bytes()))
	req.Header.Set("Content-Type", writer.FormDataContentType())
	rr := httptest.NewRecorder()
	userHandler := userHttp.NewUserHandler(uh.userUsecase)
	handler := http.HandlerFunc(userHandler.PutUser)
	handler.ServeHTTP(rr, req)

	assert.Equalf(uh.T(), http.StatusInternalServerError, rr.Code, "Should have responded with http status code %v but got %v", http.StatusInternalServerError, rr.Code)
	expectedMessage := `{"message":"` + domain.ErrInternalServerError.Error() + `"}`
	assert.Equalf(uh.T(), expectedMessage, rr.Body.String(), "Should have responded with body %s but got %s", expectedMessage, rr.Body.String())
}

func (uh *UserHandlerSuite) TestPutUserSuccessful() {
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
	uh.userUsecase.On("UpdateUser", mock.AnythingOfType("*domain.User"), mock.AnythingOfType("string"), mock.Anything).Return(nil)
	req, _ := http.NewRequest("PUT", "/users/userid1", bytes.NewReader(body.Bytes()))
	req.Header.Set("Content-Type", writer.FormDataContentType())
	rr := httptest.NewRecorder()
	userHandler := userHttp.NewUserHandler(uh.userUsecase)
	handler := http.HandlerFunc(userHandler.PutUser)
	handler.ServeHTTP(rr, req)

	assert.Equalf(uh.T(), http.StatusOK, rr.Code, "Should have responded with http status code %v but got %v", http.StatusOK, rr.Code)
	expectedMessage := `{"message":"User successfully Updated"}`
	assert.Equalf(uh.T(), expectedMessage, rr.Body.String(), "Should have responded with body %s but got %s", expectedMessage, rr.Body.String())
}

func (uh *UserHandlerSuite) TestAuthenticateUserUsernameNotProvided() {
	requestBody, _ := json.Marshal(map[string]string{
		"username": "",
		"password": "jordyjordy",
	})
	req, _ := http.NewRequest("POST", "/authentications", bytes.NewBuffer(requestBody))
	rr := httptest.NewRecorder()
	userHandler := userHttp.NewUserHandler(uh.userUsecase)
	handler := http.HandlerFunc(userHandler.AuthenticateUser)
	handler.ServeHTTP(rr, req)

	assert.Equalf(uh.T(), http.StatusBadRequest, rr.Code, "Should have responded with http status code %v but got %v", http.StatusBadRequest, rr.Code)
	expectedBody := `{"message":"username must not be empty"}`
	assert.Equalf(uh.T(), expectedBody, rr.Body.String(), "Should have responded with body %s but got %s", expectedBody, rr.Body.String())
}

func (uh *UserHandlerSuite) TestAuthenticateUserPasswordNotProvided() {
	requestBody, _ := json.Marshal(map[string]string{
		"username": "jordyf15",
		"password": "",
	})
	req, _ := http.NewRequest("POST", "/authentications", bytes.NewBuffer(requestBody))
	rr := httptest.NewRecorder()
	userHandler := userHttp.NewUserHandler(uh.userUsecase)
	handler := http.HandlerFunc(userHandler.AuthenticateUser)
	handler.ServeHTTP(rr, req)

	assert.Equalf(uh.T(), http.StatusBadRequest, rr.Code, "Should have responded with http status code %v but got %v", http.StatusBadRequest, rr.Code)
	expectedBody := `{"message":"password must not be empty"}`
	assert.Equalf(uh.T(), expectedBody, rr.Body.String(), "Should have responded with body %s but got %s", expectedBody, rr.Body.String())
}

func (uh *UserHandlerSuite) TestAuthenticateUserVerifyCredentialError() {
	requestBody, _ := json.Marshal(map[string]string{
		"username": "jordyf15",
		"password": "jordyjordy",
	})
	uh.userUsecase.On("VerifyCredential", mock.AnythingOfType("string"), mock.AnythingOfType("string")).Return("", domain.ErrInternalServerError)
	req, _ := http.NewRequest("POST", "/authentications", bytes.NewBuffer(requestBody))
	rr := httptest.NewRecorder()
	userHandler := userHttp.NewUserHandler(uh.userUsecase)
	handler := http.HandlerFunc(userHandler.AuthenticateUser)
	handler.ServeHTTP(rr, req)

	assert.Equalf(uh.T(), http.StatusInternalServerError, rr.Code, "Should have responded with http status code %v but got %v", http.StatusInternalServerError, rr.Code)
	expectedBody := `{"message":"` + domain.ErrInternalServerError.Error() + `"}`
	assert.Equalf(uh.T(), expectedBody, rr.Body.String(), "Should have responded with body %s but got %s", expectedBody, rr.Body.String())
}

func (uh *UserHandlerSuite) TestAuthenticateUserSuccessful() {
	requestBody, _ := json.Marshal(map[string]string{
		"username": "jordyf15",
		"password": "jordyjordy",
	})
	uh.userUsecase.On("VerifyCredential", mock.AnythingOfType("string"), mock.AnythingOfType("string")).Return("token", nil)
	req, _ := http.NewRequest("POST", "/authentications", bytes.NewBuffer(requestBody))
	rr := httptest.NewRecorder()
	userHandler := userHttp.NewUserHandler(uh.userUsecase)
	handler := http.HandlerFunc(userHandler.AuthenticateUser)
	handler.ServeHTTP(rr, req)

	assert.Equalf(uh.T(), http.StatusOK, rr.Code, "Should have responded with http status code %v but got %v", http.StatusOK, rr.Code)
	expectedBody := `{"message":"User successfully authenticated","data":{"access_token":"token"}}`
	assert.Equalf(uh.T(), expectedBody, rr.Body.String(), "Should have responded with body %s but got %s", expectedBody, rr.Body.String())
}
