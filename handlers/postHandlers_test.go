package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"instagram-go/models"
	"io"
	"io/fs"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
)

var phInsertPostMock func(models.Post) error
var phFindAllPostMock func() ([]models.Post, error)
var phGetPostUserIdMock func(string) (string, error)
var phUpdatePostMock func(string, string) error
var phCheckIfPostExistMock func(string) (bool, error)
var phDeletePostMock func(string) error
var phCreateMock func(string) (*os.File, error)
var phCopyMock func(io.Writer, io.Reader) (int64, error)
var phGetUserIdFromTokenMock func(string) (string, error)
var phMkDirAllMock func(string, fs.FileMode) error

type postServiceMock struct {
}

func newPostServiceMock() *postServiceMock {
	return &postServiceMock{}
}

func (psm *postServiceMock) InsertPost(document models.Post) error {
	return phInsertPostMock(document)
}

func (psm *postServiceMock) FindAllPost() ([]models.Post, error) {
	return phFindAllPostMock()
}

func (psm *postServiceMock) GetPostUserId(postId string) (string, error) {
	return phGetPostUserIdMock(postId)
}

func (psm *postServiceMock) UpdatePost(postId string, caption string) error {
	return phUpdatePostMock(postId, caption)
}

func (psm *postServiceMock) CheckIfPostExist(postId string) (bool, error) {
	return phCheckIfPostExistMock(postId)
}

func (psm *postServiceMock) DeletePost(postId string) error {
	return phDeletePostMock(postId)
}

type postHandlerHeaderMock struct {
}

func newPostHandlerHeaderMock() *postHandlerHeaderMock {
	return &postHandlerHeaderMock{}
}

func (phhm *postHandlerHeaderMock) getUserIdFromToken(tokenString string) (string, error) {
	return phGetUserIdFromTokenMock(tokenString)
}

type postFileOsHandlerMock struct {
}

func newPostFileOsHandlerMock() *postFileOsHandlerMock {
	return &postFileOsHandlerMock{}
}

func (pfohm *postFileOsHandlerMock) create(name string) (*os.File, error) {
	return phCreateMock(name)
}

func (pfohm *postFileOsHandlerMock) copy(dst io.Writer, src io.Reader) (int64, error) {
	return phCopyMock(dst, src)
}

func (pfohm *postFileOsHandlerMock) mkDirAll(path string, perm fs.FileMode) error {
	return phMkDirAllMock(path, perm)
}

func TestPostPostHandler(t *testing.T) {
	postServiceMock := newPostServiceMock()
	postFileOsHandlerMock := newPostFileOsHandlerMock()
	postHandlerHeaderMock := newPostHandlerHeaderMock()
	postHandlers := NewPostHandlers(postServiceMock, postHandlerHeaderMock, postFileOsHandlerMock)
	userId := "user-ef0813d4-0eaf-47c7-a9b6-5d86ba1fd5ec"
	method := "POST"
	urlLink := "/posts"

	// getUserIdFromToken returns error
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	writer.Close()
	phGetUserIdFromTokenMock = func(s string) (string, error) {
		return "", errors.New("GetUserIdFromToken returns error")
	}
	req, err := http.NewRequest(method, urlLink, bytes.NewReader(body.Bytes()))
	if err != nil {
		t.Fatal(err.Error())
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(postHandlers.postPostHandler)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusInternalServerError {
		t.Errorf("PostPostHandler returned wrong status code: got %v instead of %v", rr.Code, http.StatusInternalServerError)
	}
	expected := `{"message":"An error has occured in our server"}`
	if rr.Body.String() != expected {
		t.Errorf("PostPostHandler returned unexpected body: got %v instead of %v", rr.Body.String(), expected)
	}

	// if os.mkdirall returns error
	body = &bytes.Buffer{}
	writer = multipart.NewWriter(body)
	writer.Close()
	phGetUserIdFromTokenMock = func(s string) (string, error) {
		return userId, nil
	}
	phMkDirAllMock = func(s string, fm fs.FileMode) error {
		return errors.New("MkDirAll returns error")
	}
	req, err = http.NewRequest(method, urlLink, bytes.NewReader(body.Bytes()))
	if err != nil {
		t.Fatal(err.Error())
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())
	rr = httptest.NewRecorder()
	handler = http.HandlerFunc(postHandlers.postPostHandler)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusInternalServerError {
		t.Errorf("PostPostHandler returned wrong status code: got %v instead of %v", rr.Code, http.StatusInternalServerError)
	}
	expected = `{"message":"An error has occured in our server"}`
	if rr.Body.String() != expected {
		t.Errorf("PostPostHandler returned unexpected body: got %v instead of %v", rr.Body.String(), expected)
	}

	// if post's visual media is empty}
	phMkDirAllMock = func(s string, fm fs.FileMode) error {
		return nil
	}
	req, err = http.NewRequest(method, urlLink, bytes.NewReader(body.Bytes()))
	if err != nil {
		t.Fatal(err.Error())
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())
	rr = httptest.NewRecorder()
	handler = http.HandlerFunc(postHandlers.postPostHandler)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("PostPostHandler returned wrong status code: got %v instead of %v", rr.Code, http.StatusBadRequest)
	}
	expected = `{"message":"Visual Medias must not be empty"}`
	if rr.Body.String() != expected {
		t.Errorf("PostPostHandler returned unexpected body: got %v instead of %v", rr.Body.String(), expected)
	}

	// if post's visual medias types is not supported
	body = &bytes.Buffer{}
	writer = multipart.NewWriter(body)
	fw, err := writer.CreateFormFile("visual_medias", "bmp.bmp")
	if err != nil {
		t.Fatal(err.Error())
	}
	file, err := os.Open("./test_visual_medias/bmp.bmp")
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
	handler = http.HandlerFunc(postHandlers.postPostHandler)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("PostPostHandler returned wrong status code: got %v instead of %v", rr.Code, http.StatusBadRequest)
	}
	expected = `{"message":"Uploaded Visual Medias type is not supported"}`
	if rr.Body.String() != expected {
		t.Errorf("PostPostHandler returned unexpected body: got %v instead of %v", rr.Body.String(), expected)
	}

	// if os.create returns error
	phCreateMock = func(s string) (*os.File, error) {
		return nil, errors.New("os.create returns error")
	}
	body = &bytes.Buffer{}
	writer = multipart.NewWriter(body)
	fw, err = writer.CreateFormFile("visual_medias", "gif.gif")
	if err != nil {
		t.Fatal(err.Error())
	}
	file, err = os.Open("./test_visual_medias/gif.gif")
	if err != nil {
		t.Fatal(err.Error())
	}
	_, err = io.Copy(fw, file)
	if err != nil {
		t.Fatal(err.Error())
	}

	fw, err = writer.CreateFormFile("visual_medias", "jpg.jpg")
	if err != nil {
		t.Fatal(err.Error())
	}
	file, err = os.Open("./test_visual_medias/jpg.jpg")
	if err != nil {
		t.Fatal(err.Error())
	}
	_, err = io.Copy(fw, file)
	if err != nil {
		t.Fatal(err.Error())
	}

	fw, err = writer.CreateFormFile("visual_medias", "mp4.mp4")
	if err != nil {
		t.Fatal(err.Error())
	}
	file, err = os.Open("./test_visual_medias/mp4.mp4")
	if err != nil {
		t.Fatal(err.Error())
	}
	_, err = io.Copy(fw, file)
	if err != nil {
		t.Fatal(err.Error())
	}

	fw, err = writer.CreateFormFile("visual_medias", "png.png")
	if err != nil {
		t.Fatal(err.Error())
	}
	file, err = os.Open("./test_visual_medias/png.png")
	if err != nil {
		t.Fatal(err.Error())
	}
	_, err = io.Copy(fw, file)
	if err != nil {
		t.Fatal(err.Error())
	}

	fw, err = writer.CreateFormFile("visual_medias", "tiff.tiff")
	if err != nil {
		t.Fatal(err.Error())
	}
	file, err = os.Open("./test_visual_medias/tiff.tiff")
	if err != nil {
		t.Fatal(err.Error())
	}
	_, err = io.Copy(fw, file)
	if err != nil {
		t.Fatal(err.Error())
	}

	fw, err = writer.CreateFormFile("visual_medias", "webm.webm")
	if err != nil {
		t.Fatal(err.Error())
	}
	file, err = os.Open("./test_visual_medias/webm.webm")
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
	handler = http.HandlerFunc(postHandlers.postPostHandler)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusInternalServerError {
		t.Errorf("PostPostHandler returned wrong status code: got %v instead of %v", rr.Code, http.StatusInternalServerError)
	}
	expected = `{"message":"An error has occured in our server"}`
	if rr.Body.String() != expected {
		t.Errorf("PostPostHandler returned unexpected body: got %v instead of %v", rr.Body.String(), expected)
	}

	// if io.copy returns error
	phCreateMock = func(s string) (*os.File, error) {
		return nil, nil
	}
	phCopyMock = func(w io.Writer, r io.Reader) (int64, error) {
		return 0, errors.New("io.copy returns error")
	}
	req, err = http.NewRequest(method, urlLink, bytes.NewReader(body.Bytes()))
	if err != nil {
		t.Fatal(err.Error())
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())
	rr = httptest.NewRecorder()
	handler = http.HandlerFunc(postHandlers.postPostHandler)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusInternalServerError {
		t.Errorf("PostPostHandler returned wrong status code: got %v instead of %v", rr.Code, http.StatusInternalServerError)
	}
	expected = `{"message":"An error has occured in our server"}`
	if rr.Body.String() != expected {
		t.Errorf("PostPostHandler returned unexpected body: got %v instead of %v", rr.Body.String(), expected)
	}

	// if post's caption is empty
	phCopyMock = func(w io.Writer, r io.Reader) (int64, error) {
		return 0, nil
	}
	req, err = http.NewRequest(method, urlLink, bytes.NewReader(body.Bytes()))
	if err != nil {
		t.Fatal(err.Error())
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())
	rr = httptest.NewRecorder()
	handler = http.HandlerFunc(postHandlers.postPostHandler)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("PostPostHandler returned wrong status code: got %v instead of %v", rr.Code, http.StatusBadRequest)
	}
	expected = `{"message":"Caption must not be empty"}`
	if rr.Body.String() != expected {
		t.Errorf("PostPostHandler returned unexpected body: got %v instead of %v", rr.Body.String(), expected)
	}

	// if insert post returns error
	phInsertPostMock = func(p models.Post) error {
		return errors.New("InsertPost returns error")
	}
	body = &bytes.Buffer{}
	writer = multipart.NewWriter(body)
	fw, err = writer.CreateFormFile("visual_medias", "gif.gif")
	if err != nil {
		t.Fatal(err.Error())
	}
	file, err = os.Open("./test_visual_medias/gif.gif")
	if err != nil {
		t.Fatal(err.Error())
	}
	_, err = io.Copy(fw, file)
	if err != nil {
		t.Fatal(err.Error())
	}

	fw, err = writer.CreateFormField("caption")
	if err != nil {
		t.Fatal(err.Error())
	}
	_, err = io.Copy(fw, strings.NewReader("a new caption"))
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
	handler = http.HandlerFunc(postHandlers.postPostHandler)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusInternalServerError {
		t.Errorf("PostPostHandler returned wrong status code: got %v instead of %v", rr.Code, http.StatusInternalServerError)
	}
	expected = `{"message":"An error has occured in our server"}`
	if rr.Body.String() != expected {
		t.Errorf("PostPostHandler returned unexpected body: got %v instead of %v", rr.Body.String(), expected)
	}

	// if success
	phInsertPostMock = func(p models.Post) error {
		return nil
	}
	req, err = http.NewRequest(method, urlLink, bytes.NewReader(body.Bytes()))
	if err != nil {
		t.Fatal(err.Error())
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())
	rr = httptest.NewRecorder()
	handler = http.HandlerFunc(postHandlers.postPostHandler)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("PostPostHandler returned wrong status code: got %v instead of %v", rr.Code, http.StatusCreated)
	}
	expected = `{"message":"Post successfully Created"}`
	if rr.Body.String() != expected {
		t.Errorf("PostPostHandler returned unexpected body: got %v instead of %v", rr.Body.String(), expected)
	}
}

func TestGetPostsHandler(t *testing.T) {
	postServiceMock := newPostServiceMock()
	postFileOsHandlerMock := newPostFileOsHandlerMock()
	postHandlerHeaderMock := newPostHandlerHeaderMock()
	postHandlers := NewPostHandlers(postServiceMock, postHandlerHeaderMock, postFileOsHandlerMock)

	method := "GET"
	urlLink := "/posts"

	// if FindAllPost returns error
	phFindAllPostMock = func() ([]models.Post, error) {
		return []models.Post{}, errors.New("FindAllPost returns error")
	}
	requestBody, _ := json.Marshal(map[string]string{})
	req, err := http.NewRequest(method, urlLink, bytes.NewBuffer(requestBody))
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(postHandlers.getPostsHandler)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusInternalServerError {
		t.Errorf("GetPostHandler returned wrong status code: got %v instead of %v", rr.Code, http.StatusInternalServerError)
	}
	expected := `{"message":"An error has occured in our server"}`
	if rr.Body.String() != expected {
		t.Errorf("GetPostHandler returned unexpected body: got %v instead of %v", rr.Body.String(), expected)
	}

	// getAllPost successful
	phFindAllPostMock = func() ([]models.Post, error) {
		return []models.Post{
			*models.NewPost("post-"+uuid.NewString(),
				"user-"+uuid.NewString(),
				[]string{"./visual_medias/jpg.jpg"},
				"a caption", 0, time.Now(), time.Now()),
			*models.NewPost("post-"+uuid.NewString(),
				"user-"+uuid.NewString(),
				[]string{"./visual_medias/jpg.jpg"},
				"a caption", 0, time.Now(), time.Now()),
		}, nil
	}
	requestBody, _ = json.Marshal(map[string]string{})
	req, err = http.NewRequest(method, urlLink, bytes.NewBuffer(requestBody))
	if err != nil {
		t.Fatal(err)
	}
	rr = httptest.NewRecorder()
	handler = http.HandlerFunc(postHandlers.getPostsHandler)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("GetPostHandler returned wrong status code: got %v instead of %v", rr.Code, http.StatusOK)
	}
}

func TestPutPostHandler(t *testing.T) {
	postServiceMock := newPostServiceMock()
	postFileOsHandlerMock := newPostFileOsHandlerMock()
	postHandlerHeaderMock := newPostHandlerHeaderMock()
	postHandlers := NewPostHandlers(postServiceMock, postHandlerHeaderMock, postFileOsHandlerMock)
	method := "PUT"
	postId := "post-e35a9c94-4d43-4dd6-b619-a3bdd1e5324a"
	userId := "user-ef0813d4-0eaf-47c7-a9b6-5d86ba1fd5ec"
	urlLink := "/posts/" + postId

	// if getUserIdFromToken return error
	phGetUserIdFromTokenMock = func(s string) (string, error) {
		return "", errors.New("getUserIdFromToken return error")
	}
	requestBody, _ := json.Marshal(map[string]string{})
	req, err := http.NewRequest(method, urlLink, bytes.NewBuffer(requestBody))
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(postHandlers.putPostHandler)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusInternalServerError {
		t.Errorf("PutPostHandler returned wrong status code: got %v instead of %v", rr.Code, http.StatusInternalServerError)
	}
	expected := `{"message":"An error has occured in our server"}`
	if rr.Body.String() != expected {
		t.Errorf("PutPostHandler returned wrong body: got %v instead of %v", rr.Body.String(), expected)
	}

	// if CheckIfPostExist return error
	phGetUserIdFromTokenMock = func(s string) (string, error) {
		return userId, nil
	}
	phCheckIfPostExistMock = func(s string) (bool, error) {
		return false, errors.New("CheckIfPostExist returns error")
	}
	requestBody, _ = json.Marshal(map[string]string{})
	req, err = http.NewRequest(method, urlLink, bytes.NewBuffer(requestBody))
	if err != nil {
		t.Fatal(err)
	}
	rr = httptest.NewRecorder()
	handler = http.HandlerFunc(postHandlers.putPostHandler)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusInternalServerError {
		t.Errorf("PutPostHandler returned wrong status code: got %v instead of %v", rr.Code, http.StatusInternalServerError)
	}
	expected = `{"message":"An error has occured in our server"}`
	if rr.Body.String() != expected {
		t.Errorf("PutPostHandler returned wrong body: got %v instead of %v", rr.Body.String(), expected)
	}

	// if CheckIfPostExist return false
	phCheckIfPostExistMock = func(s string) (bool, error) {
		return false, nil
	}
	requestBody, _ = json.Marshal(map[string]string{})
	req, err = http.NewRequest(method, urlLink, bytes.NewBuffer(requestBody))
	if err != nil {
		t.Fatal(err)
	}
	rr = httptest.NewRecorder()
	handler = http.HandlerFunc(postHandlers.putPostHandler)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("PutPostHandler returned wrong status code: got %v instead of %v", rr.Code, http.StatusBadRequest)
	}
	expected = `{"message":"Post does not exist"}`
	if rr.Body.String() != expected {
		t.Errorf("PutPostHandler returned wrong body: got %v instead of %v", rr.Body.String(), expected)
	}

	// if GetPostUserId return error
	phCheckIfPostExistMock = func(s string) (bool, error) {
		return true, nil
	}
	phGetPostUserIdMock = func(s string) (string, error) {
		return "", errors.New("getPostUserId returns error")
	}
	requestBody, _ = json.Marshal(map[string]string{})
	req, err = http.NewRequest(method, urlLink, bytes.NewBuffer(requestBody))
	if err != nil {
		t.Fatal(err)
	}
	rr = httptest.NewRecorder()
	handler = http.HandlerFunc(postHandlers.putPostHandler)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusInternalServerError {
		t.Errorf("PutPostHandler returned wrong status code: got %v instead of %v", rr.Code, http.StatusInternalServerError)
	}
	expected = `{"message":"An error has occured in our server"}`
	if rr.Body.String() != expected {
		t.Errorf("PutPostHandler returned wrong body: got %v instead of %v", rr.Body.String(), expected)
	}

	// if user id is not the same with id returned from GetPostUserId
	phGetPostUserIdMock = func(s string) (string, error) {
		return "user-ebd67c36-a9b7-48fa-a0f9-842baf8eeb83", nil
	}
	requestBody, _ = json.Marshal(map[string]string{})
	req, err = http.NewRequest(method, urlLink, bytes.NewBuffer(requestBody))
	if err != nil {
		t.Fatal(err)
	}
	rr = httptest.NewRecorder()
	handler = http.HandlerFunc(postHandlers.putPostHandler)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusUnauthorized {
		t.Errorf("PutPostHandler returned wrong status code: got %v instead of %v", rr.Code, http.StatusUnauthorized)
	}
	expected = `{"message":"User is not authorized to update this post"}`
	if rr.Body.String() != expected {
		t.Errorf("PutPostHandler returned wrong body: got %v instead of %v", rr.Body.String(), expected)
	}

	// if caption is empty
	phGetPostUserIdMock = func(s string) (string, error) {
		return userId, nil
	}
	requestBody, _ = json.Marshal(map[string]string{})
	req, err = http.NewRequest(method, urlLink, bytes.NewBuffer(requestBody))
	if err != nil {
		t.Fatal(err)
	}
	rr = httptest.NewRecorder()
	handler = http.HandlerFunc(postHandlers.putPostHandler)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("PutPostHandler returned wrong status code: got %v instead of %v", rr.Code, http.StatusBadRequest)
	}
	expected = `{"message":"Caption must not be empty"}`
	if rr.Body.String() != expected {
		t.Errorf("PutPostHandler returned wrong body: got %v instead of %v", rr.Body.String(), expected)
	}

	// if update post returns error
	phUpdatePostMock = func(s1, s2 string) error {
		return errors.New("UpdatePost returns error")
	}
	requestBody, _ = json.Marshal(map[string]string{
		"caption": "an updated caption",
	})
	req, err = http.NewRequest(method, urlLink, bytes.NewBuffer(requestBody))
	if err != nil {
		t.Fatal(err)
	}
	rr = httptest.NewRecorder()
	handler = http.HandlerFunc(postHandlers.putPostHandler)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusInternalServerError {
		t.Errorf("PutPostHandler returned wrong status code: got %v instead of %v", rr.Code, http.StatusInternalServerError)
	}
	expected = `{"message":"An error has occured in our server"}`
	if rr.Body.String() != expected {
		t.Errorf("PutPostHandler returned wrong body: got %v instead of %v", rr.Body.String(), expected)
	}

	// if successful
	phUpdatePostMock = func(s1, s2 string) error {
		return nil
	}
	requestBody, _ = json.Marshal(map[string]string{
		"caption": "an updated caption",
	})
	req, err = http.NewRequest(method, urlLink, bytes.NewBuffer(requestBody))
	if err != nil {
		t.Fatal(err)
	}
	rr = httptest.NewRecorder()
	handler = http.HandlerFunc(postHandlers.putPostHandler)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("PutPostHandler returned wrong status code: got %v instead of %v", rr.Code, http.StatusOK)
	}
	expected = `{"message":"Post successfully Updated"}`
	if rr.Body.String() != expected {
		t.Errorf("PutPostHandler returned wrong body: got %v instead of %v", rr.Body.String(), expected)
	}
}

func TestDeletePostHandler(t *testing.T) {
	postServiceMock := newPostServiceMock()
	postFileOsHandlerMock := newPostFileOsHandlerMock()
	postHandlerHeaderMock := newPostHandlerHeaderMock()
	postHandlers := NewPostHandlers(postServiceMock, postHandlerHeaderMock, postFileOsHandlerMock)
	method := "DELETE"
	postId := "post-e35a9c94-4d43-4dd6-b619-a3bdd1e5324a"
	userId := "user-ef0813d4-0eaf-47c7-a9b6-5d86ba1fd5ec"
	urlLink := "/posts/" + postId

	// if getUserIdFromToken return error
	phGetUserIdFromTokenMock = func(s string) (string, error) {
		return "", errors.New("GetUserIdFromToken return error")
	}
	requestBody, _ := json.Marshal(map[string]string{})
	req, err := http.NewRequest(method, urlLink, bytes.NewBuffer(requestBody))
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(postHandlers.deletePostHandler)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusInternalServerError {
		t.Errorf("DeletePostHandler returned wrong status code: got %v instead of %v", rr.Code, http.StatusInternalServerError)
	}
	expected := `{"message":"An error has occured in our server"}`
	if rr.Body.String() != expected {
		t.Errorf("DeletePostHandler returned wrong body: got %v instead of %v", rr.Body.String(), expected)
	}

	// if checkIfPostExist return error
	phGetUserIdFromTokenMock = func(s string) (string, error) {
		return userId, nil
	}
	phCheckIfPostExistMock = func(s string) (bool, error) {
		return false, errors.New("checkIfPostExist return error")
	}
	requestBody, _ = json.Marshal(map[string]string{})
	req, err = http.NewRequest(method, urlLink, bytes.NewBuffer(requestBody))
	if err != nil {
		t.Fatal(err)
	}
	rr = httptest.NewRecorder()
	handler = http.HandlerFunc(postHandlers.deletePostHandler)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusInternalServerError {
		t.Errorf("DeletePostHandler returned wrong status code: got %v instead of %v", rr.Code, http.StatusInternalServerError)
	}
	expected = `{"message":"An error has occured in our server"}`
	if rr.Body.String() != expected {
		t.Errorf("DeletePostHandler returned wrong body: got %v instead of %v", rr.Body.String(), expected)
	}

	// if checkIfPostExist return false
	phCheckIfPostExistMock = func(s string) (bool, error) {
		return false, nil
	}
	requestBody, _ = json.Marshal(map[string]string{})
	req, err = http.NewRequest(method, urlLink, bytes.NewBuffer(requestBody))
	if err != nil {
		t.Fatal(err)
	}
	rr = httptest.NewRecorder()
	handler = http.HandlerFunc(postHandlers.deletePostHandler)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("DeletePostHandler returned wrong status code: got %v instead of %v", rr.Code, http.StatusBadRequest)
	}
	expected = `{"message":"Post does not exist"}`
	if rr.Body.String() != expected {
		t.Errorf("DeletePostHandler returned wrong body: got %v instead of %v", rr.Body.String(), expected)
	}

	// if getPostUserId return error
	phCheckIfPostExistMock = func(s string) (bool, error) {
		return true, nil
	}
	phGetPostUserIdMock = func(s string) (string, error) {
		return "", errors.New("getPostUserId return error")
	}
	requestBody, _ = json.Marshal(map[string]string{})
	req, err = http.NewRequest(method, urlLink, bytes.NewBuffer(requestBody))
	if err != nil {
		t.Fatal(err)
	}
	rr = httptest.NewRecorder()
	handler = http.HandlerFunc(postHandlers.deletePostHandler)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusInternalServerError {
		t.Errorf("DeletePostHandler returned wrong status code: got %v instead of %v", rr.Code, http.StatusInternalServerError)
	}
	expected = `{"message":"An error has occured in our server"}`
	if rr.Body.String() != expected {
		t.Errorf("DeletePostHandler returned wrong body: got %v instead of %v", rr.Body.String(), expected)
	}

	// if userid is not the same with getPostUserIdResult
	phGetPostUserIdMock = func(s string) (string, error) {
		return "user-82fec5d8-4682-459b-9c37-69c3bd71f6a3", nil
	}
	requestBody, _ = json.Marshal(map[string]string{})
	req, err = http.NewRequest(method, urlLink, bytes.NewBuffer(requestBody))
	if err != nil {
		t.Fatal(err)
	}
	rr = httptest.NewRecorder()
	handler = http.HandlerFunc(postHandlers.deletePostHandler)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusUnauthorized {
		t.Errorf("DeletePostHandler returned wrong status code: got %v instead of %v", rr.Code, http.StatusUnauthorized)
	}
	expected = `{"message":"User is not authorized to delete this post"}`
	if rr.Body.String() != expected {
		t.Errorf("DeletePostHandler returned wrong body: got %v instead of %v", rr.Body.String(), expected)
	}

	// if deletePost return error
	phGetPostUserIdMock = func(s string) (string, error) {
		return userId, nil
	}
	phDeletePostMock = func(s string) error {
		return errors.New("DeletePost return error")
	}
	requestBody, _ = json.Marshal(map[string]string{})
	req, err = http.NewRequest(method, urlLink, bytes.NewBuffer(requestBody))
	if err != nil {
		t.Fatal(err)
	}
	rr = httptest.NewRecorder()
	handler = http.HandlerFunc(postHandlers.deletePostHandler)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusInternalServerError {
		t.Errorf("DeletePostHandler returned wrong status code: got %v instead of %v", rr.Code, http.StatusInternalServerError)
	}
	expected = `{"message":"An error has occured in our server"}`
	if rr.Body.String() != expected {
		t.Errorf("DeletePostHandler returned wrong body: got %v instead of %v", rr.Body.String(), expected)
	}

	// if delete successfull
	phDeletePostMock = func(s string) error {
		return nil
	}
	requestBody, _ = json.Marshal(map[string]string{})
	req, err = http.NewRequest(method, urlLink, bytes.NewBuffer(requestBody))
	if err != nil {
		t.Fatal(err)
	}
	rr = httptest.NewRecorder()
	handler = http.HandlerFunc(postHandlers.deletePostHandler)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("DeletePostHandler returned wrong status code: got %v instead of %v", rr.Code, http.StatusOK)
	}
	expected = `{"message":"Post successfully Deleted"}`
	if rr.Body.String() != expected {
		t.Errorf("DeletePostHandler returned wrong body: got %v instead of %v", rr.Body.String(), expected)
	}
}
