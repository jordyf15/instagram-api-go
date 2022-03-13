package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"image"
	_ "image/gif"
	_ "image/png"
	"instagram-go/models"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/google/uuid"
)

var insertUserMock func(models.User) error
var updateUserMock func(models.User) error
var checkIfUsernameExist func(string) (bool, error)
var checkIfUserExist func(string) (bool, error)
var decodeImageMock func(io.Reader) (image.Image, string, error)
var resizeAndSaveFileToLocaleMock func(string, image.Image, string, string) (string, error)
var getUserIdFromTokenMock func(string) (string, error)

type userServiceMock struct {
}

func newUserServiceMock() *userServiceMock {
	return &userServiceMock{}
}

func (usm *userServiceMock) InsertUser(newUser models.User) error {
	return insertUserMock(newUser)
}

func (usm *userServiceMock) UpdateUser(updatedUser models.User) error {
	return updateUserMock(updatedUser)
}

func (usm *userServiceMock) CheckIfUsernameExist(username string) (bool, error) {
	return checkIfUsernameExist(username)
}

func (usm *userServiceMock) CheckIfUserExist(id string) (bool, error) {
	return checkIfUserExist(id)
}

type userHandlerHeaderMock struct {
}

func newUserHandlerHeaderMock() *userHandlerHeaderMock {
	return &userHandlerHeaderMock{}
}

func (uhhm *userHandlerHeaderMock) getUserIdFromToken(tokenString string) (string, error) {
	return getUserIdFromTokenMock(tokenString)
}

type fileOsHandlerMock struct {
}

func newFileOsHandlerMock() *fileOsHandlerMock {
	return &fileOsHandlerMock{}
}

func (fohm *fileOsHandlerMock) decodeImage(r io.Reader) (image.Image, string, error) {
	return decodeImageMock(r)
}

func (fohm *fileOsHandlerMock) resizeAndSaveFileToLocale(size string, originalProfilePicture image.Image, userId string, fileType string) (string, error) {
	return resizeAndSaveFileToLocaleMock(size, originalProfilePicture, userId, fileType)
}

func TestPostUserHandler(t *testing.T) {
	userServiceMock := newUserServiceMock()
	fileOsHandlerMock := newFileOsHandlerMock()
	userHandlerHeaderMock := newUserHandlerHeaderMock()
	userHandlers := NewUserHandlers(userServiceMock, userHandlerHeaderMock, fileOsHandlerMock)
	insertUserMock = func(u models.User) error {
		return nil
	}
	checkIfUsernameExist = func(s string) (bool, error) {
		return false, nil
	}
	method := "POST"
	url := "/users"

	// if email is not provided
	requestBody, _ := json.Marshal(map[string]string{
		"fullname": "jordy ferdian",
		"username": "jordyf15",
		"password": "jordyjordy",
	})
	req, err := http.NewRequest(method, url, bytes.NewBuffer(requestBody))
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(userHandlers.PostUserHandler)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("PostUserHandler returned wrong status code: got %v instead of %v",
			status, http.StatusBadRequest)
	}
	expected := `{"message":"Email is not provided"}`
	if rr.Body.String() != expected {
		t.Errorf("PostUserHandler returned unexpected body: got %v instead of %v",
			rr.Body.String(), expected)
	}

	// if fullname is not provided
	requestBody, _ = json.Marshal(map[string]string{
		"email":    "jordyferdian@gmail.com",
		"username": "jordyf15",
		"password": "jordyjordy",
	})
	req, err = http.NewRequest(method, url, bytes.NewBuffer(requestBody))
	if err != nil {
		t.Fatal(err)
	}
	rr = httptest.NewRecorder()
	handler = http.HandlerFunc(userHandlers.PostUserHandler)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("PostUserHandler returned wrong status code: got %v instead of %v",
			status, http.StatusBadRequest)
	}
	expected = `{"message":"Full Name is not provided"}`
	if rr.Body.String() != expected {
		t.Errorf("PostUserHandler returned unexpected body: got %v instead of %v",
			rr.Body.String(), expected)
	}

	// if username is not provided
	requestBody, _ = json.Marshal(map[string]string{
		"email":    "jordyferdian@gmail.com",
		"fullname": "jordyferdian",
		"password": "jordyjordy",
	})
	req, err = http.NewRequest(method, url, bytes.NewBuffer(requestBody))
	if err != nil {
		t.Fatal(err)
	}
	rr = httptest.NewRecorder()
	handler = http.HandlerFunc(userHandlers.PostUserHandler)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("PostUserHandler returned wrong status code: got %v instead of %v",
			status, http.StatusBadRequest)
	}
	expected = `{"message":"Username is not provided"}`
	if rr.Body.String() != expected {
		t.Errorf("PostUserHandler returned unexpected body: got %v instead of %v",
			rr.Body.String(), expected)
	}

	// if password is not provided
	requestBody, _ = json.Marshal(map[string]string{
		"email":    "jordyferdian@gmail.com",
		"fullname": "jordyferdian",
		"username": "jordyf15",
	})
	req, err = http.NewRequest(method, url, bytes.NewBuffer(requestBody))
	if err != nil {
		t.Fatal(err)
	}
	rr = httptest.NewRecorder()
	handler = http.HandlerFunc(userHandlers.PostUserHandler)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("PostUserHandler returned wrong status code: got %v instead of %v",
			status, http.StatusBadRequest)
	}
	expected = `{"message":"Password is not provided"}`
	if rr.Body.String() != expected {
		t.Errorf("PostUserHandler returned unexpected body: got %v instead of %v",
			rr.Body.String(), expected)
	}

	// if none is provided
	requestBody, _ = json.Marshal(map[string]string{})
	req, err = http.NewRequest(method, url, bytes.NewBuffer(requestBody))
	if err != nil {
		t.Fatal(err)
	}
	rr = httptest.NewRecorder()
	handler = http.HandlerFunc(userHandlers.PostUserHandler)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("PostUserHandler returned wrong status code: got %v instead of %v",
			status, http.StatusBadRequest)
	}
	expected = `{"message":"Email is not provided, Full Name is not provided, Username is not provided, Password is not provided"}`
	if rr.Body.String() != expected {
		t.Errorf("PostUserHandler returned unexpected body: got %v instead of %v",
			rr.Body.String(), expected)
	}

	// if success
	requestBody, _ = json.Marshal(map[string]string{
		"email":    "jordyferdian@gmail.com",
		"fullname": "jordyferdian",
		"password": "jordyjordy",
		"username": "jordyf15",
	})
	req, err = http.NewRequest(method, url, bytes.NewBuffer(requestBody))
	if err != nil {
		t.Fatal(err)
	}
	rr = httptest.NewRecorder()
	handler = http.HandlerFunc(userHandlers.PostUserHandler)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("PostUserHandler returned wrong status code: got %v instead of %v",
			status, http.StatusBadRequest)
	}
	expected = `{"message":"User successfully registered"}`
	if rr.Body.String() != expected {
		t.Errorf("PostUserHandler returned unexpected body: got %v instead of %v",
			rr.Body.String(), expected)
	}

	checkIfUsernameExist = func(s string) (bool, error) {
		return true, nil
	}
	// if username already exist
	requestBody, _ = json.Marshal(map[string]string{
		"email":    "jordyferdian@gmail.com",
		"fullname": "jordyferdian",
		"password": "jordyjordy",
		"username": "jordyf15",
	})
	req, err = http.NewRequest(method, url, bytes.NewBuffer(requestBody))
	if err != nil {
		t.Fatal(err)
	}
	rr = httptest.NewRecorder()
	handler = http.HandlerFunc(userHandlers.PostUserHandler)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("PostUserHandler returned wrong status code: got %v instead of %v",
			status, http.StatusBadRequest)
	}
	expected = `{"message":"Username already exist"}`
	if rr.Body.String() != expected {
		t.Errorf("PostUserHandler returned unexpected body: got %v instead of %v",
			rr.Body.String(), expected)
	}

	// if InsertOne query return an error
	insertUserMock = func(u models.User) error {
		return errors.New("InsertOne on db returns error")
	}
	checkIfUsernameExist = func(s string) (bool, error) {
		return false, nil
	}
	requestBody, _ = json.Marshal(map[string]string{
		"email":    "jordyferdian@gmail.com",
		"fullname": "jordyferdian",
		"password": "jordyjordy",
		"username": "jordyf15",
	})
	req, err = http.NewRequest(method, url, bytes.NewBuffer(requestBody))
	if err != nil {
		t.Fatal(err)
	}
	rr = httptest.NewRecorder()
	handler = http.HandlerFunc(userHandlers.PostUserHandler)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusInternalServerError {
		t.Errorf("PostUserHandler returned wrong status code: got %v instead of %v",
			status, http.StatusBadRequest)
	}
	expected = `{"message":"An error has occured in our server"}`
	if rr.Body.String() != expected {
		t.Errorf("PostUserHandler returned unexpected body: got %v instead of %v",
			rr.Body.String(), expected)
	}

	// if Find query return an error
	insertUserMock = func(u models.User) error {
		return nil
	}
	checkIfUsernameExist = func(s string) (bool, error) {
		return false, errors.New("Find on db returns error")
	}
	requestBody, _ = json.Marshal(map[string]string{
		"email":    "jordyferdian@gmail.com",
		"fullname": "jordyferdian",
		"password": "jordyjordy",
		"username": "jordyf15",
	})
	req, err = http.NewRequest(method, url, bytes.NewBuffer(requestBody))
	if err != nil {
		t.Fatal(err)
	}
	rr = httptest.NewRecorder()
	handler = http.HandlerFunc(userHandlers.PostUserHandler)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusInternalServerError {
		t.Errorf("PostUserHandler returned wrong status code: got %v instead of %v",
			status, http.StatusBadRequest)
	}
	expected = `{"message":"An error has occured in our server"}`
	if rr.Body.String() != expected {
		t.Errorf("PostUserHandler returned unexpected body: got %v instead of %v",
			rr.Body.String(), expected)
	}
}

func TestPutUserHandler(t *testing.T) {
	userServiceMock := newUserServiceMock()
	fileOsHandlerMock := newFileOsHandlerMock()
	userHandlerHeaderMock := newUserHandlerHeaderMock()
	userHandlers := NewUserHandlers(userServiceMock, userHandlerHeaderMock, fileOsHandlerMock)
	method := "PUT"
	userId := "user-45d6cd8e-795a-4710-9e81-5332d57e819b"
	urlLink := "/users/" + userId

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	fw, err := writer.CreateFormField("username")
	if err != nil {
		t.Fatal(err.Error())
	}
	_, err = io.Copy(fw, strings.NewReader("jordyf88"))
	if err != nil {
		t.Fatal(err.Error())
	}
	fw, err = writer.CreateFormField("full_name")
	if err != nil {
		t.Fatal(err.Error())
	}
	_, err = io.Copy(fw, strings.NewReader("jordy feru"))
	if err != nil {
		t.Fatal(err.Error())
	}
	fw, err = writer.CreateFormField("password")
	if err != nil {
		t.Fatal(err.Error())
	}
	_, err = io.Copy(fw, strings.NewReader("jorjor123"))
	if err != nil {
		t.Fatal(err.Error())
	}
	fw, err = writer.CreateFormField("email")
	if err != nil {
		t.Fatal(err.Error())
	}
	_, err = io.Copy(fw, strings.NewReader("jordyjordy@gmail.com"))
	if err != nil {
		t.Fatal(err.Error())
	}
	writer.Close()

	// // if user does not exist
	updateUserMock = func(u models.User) error {
		return nil
	}
	checkIfUserExist = func(s string) (bool, error) {
		return false, nil
	}
	decodeImageMock = func(r io.Reader) (image.Image, string, error) {
		return nil, "", nil
	}
	resizeAndSaveFileToLocaleMock = func(s1 string, i image.Image, s2 string, s3 string) (string, error) {
		return "", nil
	}
	getUserIdFromTokenMock = func(s string) (string, error) {
		return userId, nil
	}
	req, err := http.NewRequest(method, urlLink, bytes.NewReader(body.Bytes()))
	if err != nil {
		t.Fatal(err.Error())
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(userHandlers.PutUserHandler)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("PutUserHandler returned wrong status code: got %v instead of %v", status, http.StatusBadRequest)
	}
	expected := `{"message":"User does not exist"}`
	if rr.Body.String() != expected {
		t.Errorf("PutUserHandler returned unexpected body: got %v instead of %v", rr.Body.String(), expected)
	}

	// if the user is not authorized
	checkIfUserExist = func(s string) (bool, error) {
		return true, nil
	}
	getUserIdFromTokenMock = func(s string) (string, error) {
		return "user-" + uuid.NewString(), nil
	}
	req, err = http.NewRequest(method, urlLink, bytes.NewReader(body.Bytes()))
	if err != nil {
		t.Fatal(err.Error())
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())
	rr = httptest.NewRecorder()
	handler = http.HandlerFunc(userHandlers.PutUserHandler)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusUnauthorized {
		t.Errorf("PutUserHandler returned wrong status code: got %v instead of %v", status, http.StatusUnauthorized)
	}
	expected = `{"message":"User is not authorized"}`
	if rr.Body.String() != expected {
		t.Errorf("PutUserHandler returned unexpected body: got %v instead of %v", rr.Body.String(), expected)
	}

	// if uploaded profile_picture is not a picture
	getUserIdFromTokenMock = func(s string) (string, error) {
		return userId, nil
	}
	body = &bytes.Buffer{}
	writer = multipart.NewWriter(body)
	fw, err = writer.CreateFormField("username")
	if err != nil {
		t.Fatal(err.Error())
	}
	_, err = io.Copy(fw, strings.NewReader("jordyf88"))
	if err != nil {
		t.Fatal(err.Error())
	}
	fw, err = writer.CreateFormField("full_name")
	if err != nil {
		t.Fatal(err.Error())
	}
	_, err = io.Copy(fw, strings.NewReader("jordy feru"))
	if err != nil {
		t.Fatal(err.Error())
	}
	fw, err = writer.CreateFormField("password")
	if err != nil {
		t.Fatal(err.Error())
	}
	_, err = io.Copy(fw, strings.NewReader("jorjor123"))
	if err != nil {
		t.Fatal(err.Error())
	}
	fw, err = writer.CreateFormField("email")
	if err != nil {
		t.Fatal(err.Error())
	}
	_, err = io.Copy(fw, strings.NewReader("jordyjordy@gmail.com"))
	if err != nil {
		t.Fatal(err.Error())
	}
	fw, err = writer.CreateFormFile("profile_picture", "tiff.tiff")
	if err != nil {
		t.Fatal(err.Error())
	}
	file, err := os.Open("./test_profile_pictures/tiff.tiff")
	if err != nil {
		t.Fatal(err.Error())
	}
	_, err = io.Copy(fw, file)
	if err != nil {
		t.Fatal(err.Error())
	}
	writer.Close()
	req, err = http.NewRequest(method, urlLink, bytes.NewReader(body.Bytes()))
	if err != nil {
		t.Fatal(err.Error())
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())
	rr = httptest.NewRecorder()
	handler = http.HandlerFunc(userHandlers.PutUserHandler)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("PutUserHandler returned wrong status code: got %v instead of %v", status, http.StatusBadRequest)
	}
	expected = `{"message":"Invalid profile picture file type"}`
	if rr.Body.String() != expected {
		t.Errorf("PutUserHandler returned unexpected body: got %v instead of %v", rr.Body.String(), expected)
	}

	// if image decode returns error
	decodeImageMock = func(r io.Reader) (image.Image, string, error) {
		return nil, "", errors.New("decode image returns error")
	}
	body = &bytes.Buffer{}
	writer = multipart.NewWriter(body)
	fw, err = writer.CreateFormField("username")
	if err != nil {
		t.Fatal(err.Error())
	}
	_, err = io.Copy(fw, strings.NewReader("jordyf88"))
	if err != nil {
		t.Fatal(err.Error())
	}
	fw, err = writer.CreateFormField("full_name")
	if err != nil {
		t.Fatal(err.Error())
	}
	_, err = io.Copy(fw, strings.NewReader("jordy feru"))
	if err != nil {
		t.Fatal(err.Error())
	}
	fw, err = writer.CreateFormField("password")
	if err != nil {
		t.Fatal(err.Error())
	}
	_, err = io.Copy(fw, strings.NewReader("jorjor123"))
	if err != nil {
		t.Fatal(err.Error())
	}
	fw, err = writer.CreateFormField("email")
	if err != nil {
		t.Fatal(err.Error())
	}
	_, err = io.Copy(fw, strings.NewReader("jordyjordy@gmail.com"))
	if err != nil {
		t.Fatal(err.Error())
	}
	fw, err = writer.CreateFormFile("profile_picture", "png.png")
	if err != nil {
		t.Fatal(err.Error())
	}
	file, err = os.Open("./test_profile_pictures/png.png")
	if err != nil {
		t.Fatal(err.Error())
	}
	_, err = io.Copy(fw, file)
	if err != nil {
		t.Fatal(err.Error())
	}
	writer.Close()
	req, err = http.NewRequest(method, urlLink, bytes.NewReader(body.Bytes()))
	if err != nil {
		t.Fatal(err.Error())
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())
	rr = httptest.NewRecorder()
	handler = http.HandlerFunc(userHandlers.PutUserHandler)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusInternalServerError {
		t.Errorf("PutUserHandler returned wrong status code: got %v instead of %v", status, http.StatusBadRequest)
	}
	expected = `{"message":"An error has occured in our server"}`
	if rr.Body.String() != expected {
		t.Errorf("PutUserHandler returned unexpected body: got %v instead of %v", rr.Body.String(), expected)
	}

	// if savefiletolocal returns error
	decodeImageMock = func(r io.Reader) (image.Image, string, error) {
		return nil, "", nil
	}
	resizeAndSaveFileToLocaleMock = func(s1 string, i image.Image, s2 string, s3 string) (string, error) {
		return "", errors.New("resizeAndSaveFileToLocal returns error")
	}
	req, err = http.NewRequest(method, urlLink, bytes.NewReader(body.Bytes()))
	if err != nil {
		t.Fatal(err.Error())
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())
	rr = httptest.NewRecorder()
	handler = http.HandlerFunc(userHandlers.PutUserHandler)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusInternalServerError {
		t.Errorf("PutUserHandler returned wrong status code: got %v instead of %v", status, http.StatusBadRequest)
	}
	expected = `{"message":"An error has occured in our server"}`
	if rr.Body.String() != expected {
		t.Errorf("PutUserHandler returned unexpected body: got %v instead of %v", rr.Body.String(), expected)
	}

	// if update user throws error
	resizeAndSaveFileToLocaleMock = func(s1 string, i image.Image, s2 string, s3 string) (string, error) {
		return "", nil
	}
	updateUserMock = func(u models.User) error {
		return errors.New("updateUser to db returns error")
	}
	req, err = http.NewRequest(method, urlLink, bytes.NewReader(body.Bytes()))
	if err != nil {
		t.Fatal(err.Error())
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())
	rr = httptest.NewRecorder()
	handler = http.HandlerFunc(userHandlers.PutUserHandler)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusInternalServerError {
		t.Errorf("PutUserHandler returned wrong status code: got %v instead of %v", status, http.StatusBadRequest)
	}
	expected = `{"message":"An error has occured in our server"}`
	if rr.Body.String() != expected {
		t.Errorf("PutUserHandler returned unexpected body: got %v instead of %v", rr.Body.String(), expected)
	}

	// if successful
	updateUserMock = func(u models.User) error {
		return nil
	}
	req, err = http.NewRequest(method, urlLink, bytes.NewReader(body.Bytes()))
	if err != nil {
		t.Fatal(err.Error())
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())
	rr = httptest.NewRecorder()
	handler = http.HandlerFunc(userHandlers.PutUserHandler)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("PutUserHandler returned wrong status code: got %v instead of %v", status, http.StatusBadRequest)
	}
	expected = `{"message":"User successfully Updated"}`
	if rr.Body.String() != expected {
		t.Errorf("PutUserHandler returned unexpected body: got %v instead of %v", rr.Body.String(), expected)
	}
}
