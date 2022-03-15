package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"instagram-go/models"
	"net/http"
	"net/http/httptest"
	"testing"
)

var commentGetUserIdFromTokenMock func(string) (string, error)

var commentCheckIfPostExistMock func(string) (bool, error)

var commentFindAllPostCommentMock func(string) ([]models.Comment, error)
var commentInsertCommentMock func(models.Comment) error
var commentGetCommentUserIdMock func(string) (string, error)
var commentUpdateCommentMock func(string, string) error
var commentCheckIfCommentExistMock func(string) (bool, error)
var commentDeleteCommentMock func(string) error

type commentHandlerHeaderMock struct {
}

func newCommentHandlerHeaderMock() *commentHandlerHeaderMock {
	return &commentHandlerHeaderMock{}
}
func (chhm *commentHandlerHeaderMock) getUserIdFromToken(tokenString string) (string, error) {
	return commentGetUserIdFromTokenMock(tokenString)
}

type commentCommentServiceMock struct {
}

func newCommentCommentServiceMock() *commentCommentServiceMock {
	return &commentCommentServiceMock{}
}

func (ccsm *commentCommentServiceMock) GetCommentUserId(commentId string) (string, error) {
	return commentGetCommentUserIdMock(commentId)
}

func (ccsm *commentCommentServiceMock) FindAllPostComment(postId string) ([]models.Comment, error) {
	return commentFindAllPostCommentMock(postId)
}

func (ccsm *commentCommentServiceMock) InsertComment(comment models.Comment) error {
	return commentInsertCommentMock(comment)
}

func (ccsm *commentCommentServiceMock) UpdateComment(commentId string, comment string) error {
	return commentUpdateCommentMock(commentId, comment)
}

func (ccsm *commentCommentServiceMock) DeleteComment(commentId string) error {
	return commentDeleteCommentMock(commentId)
}

func (ccsm *commentCommentServiceMock) CheckIfCommentExist(commentId string) (bool, error) {
	return commentCheckIfCommentExistMock(commentId)
}

type commentPostServiceMock struct {
}

func newCommentPostServiceMock() *commentPostServiceMock {
	return &commentPostServiceMock{}
}

func (cpsm *commentPostServiceMock) InsertPost(post models.Post) error {
	return nil
}

func (cpsm *commentPostServiceMock) FindAllPost() ([]models.Post, error) {
	return nil, nil
}

func (cpsm *commentPostServiceMock) GetPostUserId(postId string) (string, error) {
	return "", nil
}

func (cpsm *commentPostServiceMock) UpdatePost(postId string, caption string) error {
	return nil
}

func (cpsm *commentPostServiceMock) CheckIfPostExist(postId string) (bool, error) {
	return commentCheckIfPostExistMock(postId)
}

func (cpsm *commentPostServiceMock) DeletePost(postId string) error {
	return nil
}
func TestGetComments(t *testing.T) {
	commentCommentServiceMock := newCommentCommentServiceMock()
	commentPostServiceMock := newCommentPostServiceMock()
	commentHandlerHeaderMock := newCommentHandlerHeaderMock()
	commentHandlers := NewCommentHandlers(commentCommentServiceMock, commentPostServiceMock, commentHandlerHeaderMock)
	method := "GET"
	postId := "post-566feabf-113f-4739-a226-9ba7cd6f4fc1"
	urlLink := "http://localhost:8000/posts/" + postId + "/comments"

	// if checkifpostexist return error
	commentCheckIfPostExistMock = func(s string) (bool, error) {
		return false, errors.New("CheckIfPostExist return error")
	}
	requestBody, _ := json.Marshal(map[string]string{})
	req, err := http.NewRequest(method, urlLink, bytes.NewBuffer(requestBody))
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(commentHandlers.getComments)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusInternalServerError {
		t.Errorf("getComments returned wrong status code: got %v instead of %v", rr.Code, status)
	}
	expected := `{"message":"An error has occured in our server"}`
	if rr.Body.String() != expected {
		t.Errorf("getComments returned wrong body: got %v instead of %v", rr.Body.String(), expected)
	}

	// if checkIfPostExist return false
	commentCheckIfPostExistMock = func(s string) (bool, error) {
		return false, nil
	}
	requestBody, _ = json.Marshal(map[string]string{})
	req, err = http.NewRequest(method, urlLink, bytes.NewBuffer(requestBody))
	if err != nil {
		t.Fatal(err)
	}
	rr = httptest.NewRecorder()
	handler = http.HandlerFunc(commentHandlers.getComments)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("getComments returned wrong status code: got %v instead of %v", rr.Code, http.StatusBadRequest)
	}
	expected = `{"message":"Post does not exist"}`
	if rr.Body.String() != expected {
		t.Errorf("getComments returned wrong body: got %v instead of %v", rr.Body.String(), expected)
	}

	// if findAllPostComment return error
	commentCheckIfPostExistMock = func(s string) (bool, error) {
		return true, nil
	}
	commentFindAllPostCommentMock = func(s string) ([]models.Comment, error) {
		return nil, errors.New("findAllPostComment return error")
	}
	requestBody, _ = json.Marshal(map[string]string{})
	req, err = http.NewRequest(method, urlLink, bytes.NewBuffer(requestBody))
	if err != nil {
		t.Fatal(err)
	}
	rr = httptest.NewRecorder()
	handler = http.HandlerFunc(commentHandlers.getComments)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusInternalServerError {
		t.Errorf("getComments returned wrong status code: got %v instead of %v", rr.Code, http.StatusInternalServerError)
	}
	expected = `{"message":"An error has occured in our server"}`
	if rr.Body.String() != expected {
		t.Errorf("getComments returned wrong body: got %v instead of %v", rr.Body.String(), expected)
	}

	// if successful
	commentFindAllPostCommentMock = func(s string) ([]models.Comment, error) {
		return []models.Comment{}, nil
	}
	requestBody, _ = json.Marshal(map[string]string{})
	req, err = http.NewRequest(method, urlLink, bytes.NewBuffer(requestBody))
	if err != nil {
		t.Fatal(err)
	}
	rr = httptest.NewRecorder()
	handler = http.HandlerFunc(commentHandlers.getComments)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("getComments returned wrong status code: got %v instead of %v", rr.Code, http.StatusOK)
	}
}

func TestPostComment(t *testing.T) {
	commentCommentServiceMock := newCommentCommentServiceMock()
	commentPostServiceMock := newCommentPostServiceMock()
	commentHandlerHeaderMock := newCommentHandlerHeaderMock()
	commentHandlers := NewCommentHandlers(commentCommentServiceMock, commentPostServiceMock, commentHandlerHeaderMock)
	method := "POST"
	postId := "post-566feabf-113f-4739-a226-9ba7cd6f4fc1"
	urlLink := "http://localhost:8000/posts/" + postId + "/comments"
	userId := "user-182ef588-bf73-424e-92d4-9da61af3c418"
	// if getUserIdFromToken return error
	commentGetUserIdFromTokenMock = func(s string) (string, error) {
		return "", errors.New("getUserIdFromToken return error")
	}
	requestBody, _ := json.Marshal(map[string]string{})
	req, err := http.NewRequest(method, urlLink, bytes.NewBuffer(requestBody))
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(commentHandlers.postComment)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusInternalServerError {
		t.Errorf("postComment returned wrong status code: got %v instead of %v", rr.Code, http.StatusInternalServerError)
	}
	expected := `{"message":"An error has occured in our server"}`
	if rr.Body.String() != expected {
		t.Errorf("postComment returned wrong body: got %v instead of %v", rr.Body.String(), expected)
	}

	// if checkifpostexist return error
	commentGetUserIdFromTokenMock = func(s string) (string, error) {
		return userId, nil
	}
	commentCheckIfPostExistMock = func(s string) (bool, error) {
		return false, errors.New("checkIfPostExist return error")
	}
	requestBody, _ = json.Marshal(map[string]string{})
	req, err = http.NewRequest(method, urlLink, bytes.NewBuffer(requestBody))
	if err != nil {
		t.Fatal(err)
	}
	rr = httptest.NewRecorder()
	handler = http.HandlerFunc(commentHandlers.postComment)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusInternalServerError {
		t.Errorf("postComment returned wrong status code: got %v instead of %v", rr.Code, http.StatusInternalServerError)
	}
	expected = `{"message":"An error has occured in our server"}`
	if rr.Body.String() != expected {
		t.Errorf("postComment returned wrong body: got %v instead of %v", rr.Body.String(), expected)
	}

	// if checkIfPostExist return false
	commentCheckIfPostExistMock = func(s string) (bool, error) {
		return false, nil
	}
	requestBody, _ = json.Marshal(map[string]string{})
	req, err = http.NewRequest(method, urlLink, bytes.NewBuffer(requestBody))
	if err != nil {
		t.Fatal(err)
	}
	rr = httptest.NewRecorder()
	handler = http.HandlerFunc(commentHandlers.postComment)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("postComment returned wrong status code: got %v instead of %v", rr.Code, http.StatusBadRequest)
	}
	expected = `{"message":"Post does not exist"}`
	if rr.Body.String() != expected {
		t.Errorf("postComment returned wrong body: got %v instead of %v", rr.Body.String(), expected)
	}

	// if comment empty
	commentCheckIfPostExistMock = func(s string) (bool, error) {
		return true, nil
	}
	requestBody, _ = json.Marshal(map[string]string{})
	req, err = http.NewRequest(method, urlLink, bytes.NewBuffer(requestBody))
	if err != nil {
		t.Fatal(err)
	}
	rr = httptest.NewRecorder()
	handler = http.HandlerFunc(commentHandlers.postComment)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("postComment returned wrong status code: got %v instead of %v", rr.Code, http.StatusBadRequest)
	}
	expected = `{"message":"Comment must not be empty"}`
	if rr.Body.String() != expected {
		t.Errorf("postComment returned wrong body: got %v instead of %v", rr.Body.String(), expected)
	}

	// if insertComment return error
	commentInsertCommentMock = func(c models.Comment) error {
		return errors.New("InsertComment return error")
	}
	requestBody, _ = json.Marshal(map[string]string{
		"comment": "a new comment",
	})
	req, err = http.NewRequest(method, urlLink, bytes.NewBuffer(requestBody))
	if err != nil {
		t.Fatal(err)
	}
	rr = httptest.NewRecorder()
	handler = http.HandlerFunc(commentHandlers.postComment)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusInternalServerError {
		t.Errorf("postComment returned wrong status code: got %v instead of %v", rr.Code, http.StatusInternalServerError)
	}
	expected = `{"message":"An error has occured in our server"}`
	if rr.Body.String() != expected {
		t.Errorf("postComment returned wrong body: got %v instead of %v", rr.Body.String(), expected)
	}

	// if successful
	commentInsertCommentMock = func(c models.Comment) error {
		return nil
	}
	requestBody, _ = json.Marshal(map[string]string{
		"comment": "a new comment",
	})
	req, err = http.NewRequest(method, urlLink, bytes.NewBuffer(requestBody))
	if err != nil {
		t.Fatal(err)
	}
	rr = httptest.NewRecorder()
	handler = http.HandlerFunc(commentHandlers.postComment)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("postComment returned wrong status code: got %v instead of %v", rr.Code, http.StatusCreated)
	}
}

func TestPutComment(t *testing.T) {
	commentCommentServiceMock := newCommentCommentServiceMock()
	commentPostServiceMock := newCommentPostServiceMock()
	commentHandlerHeaderMock := newCommentHandlerHeaderMock()
	commentHandlers := NewCommentHandlers(commentCommentServiceMock, commentPostServiceMock, commentHandlerHeaderMock)
	method := "PUT"
	postId := "post-566feabf-113f-4739-a226-9ba7cd6f4fc1"
	userId := "user-182ef588-bf73-424e-92d4-9da61af3c418"
	commentId := "comment-2b69841d-3ba8-4a34-abb7-b6650c2452dd"
	urlLink := "http://localhost:8000/posts/" + postId + "/comments/" + commentId

	// if getuseridfromtoken return error
	commentGetUserIdFromTokenMock = func(s string) (string, error) {
		return "", errors.New("getUserIdFromToken return error")
	}
	requestBody, _ := json.Marshal(map[string]string{})
	req, err := http.NewRequest(method, urlLink, bytes.NewBuffer(requestBody))
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(commentHandlers.putComment)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusInternalServerError {
		t.Errorf("putComment returned wrong status code: got %v instead of %v", rr.Code, http.StatusInternalServerError)
	}
	expected := `{"message":"An error has occured in our server"}`
	if rr.Body.String() != expected {
		t.Errorf("putComment returned wrong body: got %v instead of %v", rr.Body.String(), expected)
	}

	// if checkIfCommentExist return error
	commentGetUserIdFromTokenMock = func(s string) (string, error) {
		return userId, nil
	}
	commentCheckIfCommentExistMock = func(s string) (bool, error) {
		return false, errors.New("checkIfCommentExist return error")
	}
	requestBody, _ = json.Marshal(map[string]string{})
	req, err = http.NewRequest(method, urlLink, bytes.NewBuffer(requestBody))
	if err != nil {
		t.Fatal(err)
	}
	rr = httptest.NewRecorder()
	handler = http.HandlerFunc(commentHandlers.putComment)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusInternalServerError {
		t.Errorf("putComment returned wrong status code: got %v instead of %v", rr.Code, http.StatusInternalServerError)
	}
	expected = `{"message":"An error has occured in our server"}`
	if rr.Body.String() != expected {
		t.Errorf("putComment returned wrong body: got %v instead of %v", rr.Body.String(), expected)
	}

	// if checkIfCommentExist return false
	commentCheckIfCommentExistMock = func(s string) (bool, error) {
		return false, nil
	}
	requestBody, _ = json.Marshal(map[string]string{})
	req, err = http.NewRequest(method, urlLink, bytes.NewBuffer(requestBody))
	if err != nil {
		t.Fatal(err)
	}
	rr = httptest.NewRecorder()
	handler = http.HandlerFunc(commentHandlers.putComment)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("putComment returned wrong status code: got %v instead of %v", rr.Code, http.StatusBadRequest)
	}
	expected = `{"message":"Comment does not exist"}`
	if rr.Body.String() != expected {
		t.Errorf("putComment returned wrong body: got %v instead of %v", rr.Body.String(), expected)
	}

	// if getcommentUserId return error
	commentCheckIfCommentExistMock = func(s string) (bool, error) {
		return true, nil
	}
	commentGetCommentUserIdMock = func(s string) (string, error) {
		return "", errors.New("getCommentUserId return error")
	}
	requestBody, _ = json.Marshal(map[string]string{})
	req, err = http.NewRequest(method, urlLink, bytes.NewBuffer(requestBody))
	if err != nil {
		t.Fatal(err)
	}
	rr = httptest.NewRecorder()
	handler = http.HandlerFunc(commentHandlers.putComment)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusInternalServerError {
		t.Errorf("putComment returned wrong status code: got %v instead of %v", rr.Code, http.StatusInternalServerError)
	}
	expected = `{"message":"An error has occured in our server"}`
	if rr.Body.String() != expected {
		t.Errorf("putComment returned wrong body: got %v instead of %v", rr.Body.String(), expected)
	}

	// if user id does not match getCommentUserId result
	commentGetCommentUserIdMock = func(s string) (string, error) {
		return "user-bbc4c22f-0129-4b14-af7a-86bbd2709b80", nil
	}
	requestBody, _ = json.Marshal(map[string]string{})
	req, err = http.NewRequest(method, urlLink, bytes.NewBuffer(requestBody))
	if err != nil {
		t.Fatal(err)
	}
	rr = httptest.NewRecorder()
	handler = http.HandlerFunc(commentHandlers.putComment)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusUnauthorized {
		t.Errorf("putComment returned wrong status code: got %v instead of %v", rr.Code, http.StatusUnauthorized)
	}
	expected = `{"message":"User is not authorized to update this comment"}`
	if rr.Body.String() != expected {
		t.Errorf("putComment returned wrong body: got %v instead of %v", rr.Body.String(), expected)
	}

	// if comment is empty
	commentGetCommentUserIdMock = func(s string) (string, error) {
		return userId, nil
	}
	requestBody, _ = json.Marshal(map[string]string{})
	req, err = http.NewRequest(method, urlLink, bytes.NewBuffer(requestBody))
	if err != nil {
		t.Fatal(err)
	}
	rr = httptest.NewRecorder()
	handler = http.HandlerFunc(commentHandlers.putComment)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("putComment returned wrong status code: got %v instead of %v", rr.Code, http.StatusBadRequest)
	}
	expected = `{"message":"Comment must not be empty"}`
	if rr.Body.String() != expected {
		t.Errorf("putComment returned wrong body: got %v instead of %v", rr.Body.String(), expected)
	}

	// if updateComment return error
	commentUpdateCommentMock = func(s1, s2 string) error {
		return errors.New("UpdateComment return error")
	}
	requestBody, _ = json.Marshal(map[string]string{
		"comment": "an updated comment",
	})
	req, err = http.NewRequest(method, urlLink, bytes.NewBuffer(requestBody))
	if err != nil {
		t.Fatal(err)
	}
	rr = httptest.NewRecorder()
	handler = http.HandlerFunc(commentHandlers.putComment)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusInternalServerError {
		t.Errorf("putComment returned wrong status code: got %v instead of %v", rr.Code, http.StatusInternalServerError)
	}
	expected = `{"message":"An error has occured in our server"}`
	if rr.Body.String() != expected {
		t.Errorf("putComment returned wrong body: got %v instead of %v", rr.Body.String(), expected)
	}

	// if successful
	commentUpdateCommentMock = func(s1, s2 string) error {
		return nil
	}
	requestBody, _ = json.Marshal(map[string]string{
		"comment": "an updated comment",
	})
	req, err = http.NewRequest(method, urlLink, bytes.NewBuffer(requestBody))
	if err != nil {
		t.Fatal(err)
	}
	rr = httptest.NewRecorder()
	handler = http.HandlerFunc(commentHandlers.putComment)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("putComment returned wrong status code: got %v instead of %v", rr.Code, http.StatusOK)
	}
	expected = `{"message":"Comment successfully Updated"}`
	if rr.Body.String() != expected {
		t.Errorf("putComment returned wrong body: got %v instead of %v", rr.Body.String(), expected)
	}
}

func TestDeleteComment(t *testing.T) {
	commentCommentServiceMock := newCommentCommentServiceMock()
	commentPostServiceMock := newCommentPostServiceMock()
	commentHandlerHeaderMock := newCommentHandlerHeaderMock()
	commentHandlers := NewCommentHandlers(commentCommentServiceMock, commentPostServiceMock, commentHandlerHeaderMock)
	method := "DELETE"
	postId := "post-566feabf-113f-4739-a226-9ba7cd6f4fc1"
	userId := "user-182ef588-bf73-424e-92d4-9da61af3c418"
	commentId := "comment-2b69841d-3ba8-4a34-abb7-b6650c2452dd"
	urlLink := "http://localhost:8000/posts/" + postId + "/comments/" + commentId

	// if getuseridfromtoken return error
	commentGetUserIdFromTokenMock = func(s string) (string, error) {
		return "", errors.New("getUserIdFromToken return error")
	}
	requestBody, _ := json.Marshal(map[string]string{})
	req, err := http.NewRequest(method, urlLink, bytes.NewBuffer(requestBody))
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(commentHandlers.deleteComment)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusInternalServerError {
		t.Errorf("deleteComment returned wrong status code: got %v instead of %v", rr.Code, http.StatusInternalServerError)
	}
	expected := `{"message":"An error has occured in our server"}`
	if rr.Body.String() != expected {
		t.Errorf("deleteComment returned wrong body: got %v instead of %v", rr.Body.String(), expected)
	}

	// if checkIfCommentExist return error
	commentGetUserIdFromTokenMock = func(s string) (string, error) {
		return userId, nil
	}
	commentCheckIfCommentExistMock = func(s string) (bool, error) {
		return false, errors.New("CheckIfCommentExist return error")
	}
	requestBody, _ = json.Marshal(map[string]string{})
	req, err = http.NewRequest(method, urlLink, bytes.NewBuffer(requestBody))
	if err != nil {
		t.Fatal(err)
	}
	rr = httptest.NewRecorder()
	handler = http.HandlerFunc(commentHandlers.deleteComment)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusInternalServerError {
		t.Errorf("deleteComment returned wrong status code: got %v instead of %v", rr.Code, http.StatusInternalServerError)
	}
	expected = `{"message":"An error has occured in our server"}`
	if rr.Body.String() != expected {
		t.Errorf("deleteComment returned wrong body: got %v instead of %v", rr.Body.String(), expected)
	}

	// if checkIfCommentExist return false
	commentCheckIfCommentExistMock = func(s string) (bool, error) {
		return false, nil
	}
	requestBody, _ = json.Marshal(map[string]string{})
	req, err = http.NewRequest(method, urlLink, bytes.NewBuffer(requestBody))
	if err != nil {
		t.Fatal(err)
	}
	rr = httptest.NewRecorder()
	handler = http.HandlerFunc(commentHandlers.deleteComment)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("deleteComment returned wrong status code: got %v instead of %v", rr.Code, http.StatusBadRequest)
	}
	expected = `{"message":"Comment does not exist"}`
	if rr.Body.String() != expected {
		t.Errorf("deleteComment returned wrong body: got %v instead of %v", rr.Body.String(), expected)
	}

	// if getcommentUserId return error
	commentCheckIfCommentExistMock = func(s string) (bool, error) {
		return true, nil
	}
	commentGetCommentUserIdMock = func(s string) (string, error) {
		return "", errors.New("GetCommentUserId return error")
	}
	requestBody, _ = json.Marshal(map[string]string{})
	req, err = http.NewRequest(method, urlLink, bytes.NewBuffer(requestBody))
	if err != nil {
		t.Fatal(err)
	}
	rr = httptest.NewRecorder()
	handler = http.HandlerFunc(commentHandlers.deleteComment)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusInternalServerError {
		t.Errorf("deleteComment returned wrong status code: got %v instead of %v", rr.Code, http.StatusInternalServerError)
	}
	expected = `{"message":"An error has occured in our server"}`
	if rr.Body.String() != expected {
		t.Errorf("deleteComment returned wrong body: got %v instead of %v", rr.Body.String(), expected)
	}

	// if user id does not match getCommentUserId result
	commentGetCommentUserIdMock = func(s string) (string, error) {
		return "user-2b69841d-3ba8-4a34-abb7-b6650c2452dd", nil
	}
	requestBody, _ = json.Marshal(map[string]string{})
	req, err = http.NewRequest(method, urlLink, bytes.NewBuffer(requestBody))
	if err != nil {
		t.Fatal(err)
	}
	rr = httptest.NewRecorder()
	handler = http.HandlerFunc(commentHandlers.deleteComment)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusUnauthorized {
		t.Errorf("deleteComment returned wrong status code: got %v instead of %v", rr.Code, http.StatusUnauthorized)
	}
	expected = `{"message":"User is not authorized to delete this comment"}`
	if rr.Body.String() != expected {
		t.Errorf("deleteComment returned wrong body: got %v instead of %v", rr.Body.String(), expected)
	}

	// if deleteCOmment return error
	commentGetCommentUserIdMock = func(s string) (string, error) {
		return userId, nil
	}
	commentDeleteCommentMock = func(s string) error {
		return errors.New("deleteComment return error")
	}
	requestBody, _ = json.Marshal(map[string]string{})
	req, err = http.NewRequest(method, urlLink, bytes.NewBuffer(requestBody))
	if err != nil {
		t.Fatal(err)
	}
	rr = httptest.NewRecorder()
	handler = http.HandlerFunc(commentHandlers.deleteComment)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusInternalServerError {
		t.Errorf("deleteComment returned wrong status code: got %v instead of %v", rr.Code, http.StatusInternalServerError)
	}
	expected = `{"message":"An error has occured in our server"}`
	if rr.Body.String() != expected {
		t.Errorf("deleteComment returned wrong body: got %v instead of %v", rr.Body.String(), expected)
	}

	// if successful
	commentDeleteCommentMock = func(s string) error {
		return nil
	}
	requestBody, _ = json.Marshal(map[string]string{})
	req, err = http.NewRequest(method, urlLink, bytes.NewBuffer(requestBody))
	if err != nil {
		t.Fatal(err)
	}
	rr = httptest.NewRecorder()
	handler = http.HandlerFunc(commentHandlers.deleteComment)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("deleteComment returned wrong status code: got %v instead of %v", rr.Code, http.StatusOK)
	}
	expected = `{"message":"Comment successfully Deleted"}`
	if rr.Body.String() != expected {
		t.Errorf("deleteComment returned wrong body: got %v instead of %v", rr.Body.String(), expected)
	}
}
