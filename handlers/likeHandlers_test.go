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

var likeGetUserIdFromTokenMock func(string) (string, error)

var likeCheckIfPostExistMock func(string) (bool, error)
var likeIsLikeExistMock func(string, string, string) (bool, error)
var likeInsertLikeMock func(models.Like) error
var likeIsLikeExistByIdMock func(string) (bool, error)
var likeGetLikeUserIdMock func(string) (string, error)
var likeDeleteLikeMock func(string) error
var likeCheckIfCommentExistMock func(string) (bool, error)

type likeHandlerHeaderMock struct {
}

func newLikeHandlerHeaderMock() *likeHandlerHeaderMock {
	return &likeHandlerHeaderMock{}
}

func (lhhm *likeHandlerHeaderMock) getUserIdFromToken(tokenString string) (string, error) {
	return likeGetUserIdFromTokenMock(tokenString)
}

type likeLikeServiceMock struct {
}

func newLikeLikeServiceMock() *likeLikeServiceMock {
	return &likeLikeServiceMock{}
}

func (llsm *likeLikeServiceMock) InsertLike(like models.Like) error {
	return likeInsertLikeMock(like)
}

func (llsm *likeLikeServiceMock) DeleteLike(likeId string) error {
	return likeDeleteLikeMock(likeId)
}

func (llsm *likeLikeServiceMock) IsLikeExist(userId string, resourceId string, resourceType string) (bool, error) {
	return likeIsLikeExistMock(userId, resourceId, resourceType)
}

func (llsm *likeLikeServiceMock) GetLikeUserId(likeId string) (string, error) {
	return likeGetLikeUserIdMock(likeId)
}

func (llsm *likeLikeServiceMock) IsLikeExistById(likeId string) (bool, error) {
	return likeIsLikeExistByIdMock(likeId)
}

type likePostServiceMock struct {
}

func newLikePostServiceMock() *likePostServiceMock {
	return &likePostServiceMock{}
}

func (lpsm *likePostServiceMock) InsertPost(post models.Post) error {
	return nil
}

func (lpsm *likePostServiceMock) FindAllPost() ([]models.Post, error) {
	return nil, nil
}

func (lpsm *likePostServiceMock) GetPostUserId(postId string) (string, error) {
	return "", nil
}

func (lpsm *likePostServiceMock) UpdatePost(postId string, caption string) error {
	return nil
}

func (lpsm *likePostServiceMock) CheckIfPostExist(postId string) (bool, error) {
	return likeCheckIfPostExistMock(postId)
}

func (lpsm *likePostServiceMock) DeletePost(postId string) error {
	return nil
}

type likeCommentServiceMock struct {
}

func newLikeCommentServiceMock() *likeCommentServiceMock {
	return &likeCommentServiceMock{}
}

func (lcsm *likeCommentServiceMock) GetCommentUserId(commentId string) (string, error) {
	return "", nil
}

func (lcsm *likeCommentServiceMock) FindAllPostComment(postId string) ([]models.Comment, error) {
	return nil, nil
}

func (lcsm *likeCommentServiceMock) InsertComment(comment models.Comment) error {
	return nil
}

func (lcsm *likeCommentServiceMock) UpdateComment(commentId string, newComment string) error {
	return nil
}

func (lcsm *likeCommentServiceMock) DeleteComment(commentId string) error {
	return nil
}

func (lcsm *likeCommentServiceMock) CheckIfCommentExist(commentId string) (bool, error) {
	return likeCheckIfCommentExistMock(commentId)
}
func TestPostPostLikeHandler(t *testing.T) {
	likeLikeServiceMock := newLikeLikeServiceMock()
	likeHandlerHeaderMock := newLikeHandlerHeaderMock()
	likeCommentServiceMock := newLikeCommentServiceMock()
	likePostServiceMock := newLikePostServiceMock()
	likeHandlers := NewLikeHandlers(likeLikeServiceMock, likePostServiceMock, likeCommentServiceMock, likeHandlerHeaderMock)
	method := "POST"
	urlLink := "/posts/post-566feabf-113f-4739-a226-9ba7cd6f4fc1/likes"
	// if getUserIdToken return error
	likeGetUserIdFromTokenMock = func(s string) (string, error) {
		return "", errors.New("GetUserIdFromToken return error")
	}
	requestBody, _ := json.Marshal(map[string]string{})
	req, err := http.NewRequest(method, urlLink, bytes.NewBuffer(requestBody))
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(likeHandlers.PostPostLikeHandler)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusInternalServerError {
		t.Errorf("PostPostLikeHandler returned wrong status code: got %v instead of %v",
			rr.Code, http.StatusInternalServerError)
	}
	expected := `{"message":"An error has occured in our server"}`
	if rr.Body.String() != expected {
		t.Errorf("PostPostLikeHandler returned unexpected body: got %v instead of %v",
			rr.Body.String(), expected)
	}

	// if checkIfPostExist return error
	likeGetUserIdFromTokenMock = func(s string) (string, error) {
		return "", nil
	}
	likeCheckIfPostExistMock = func(s string) (bool, error) {
		return false, errors.New("CheckIfPostExist return error")
	}
	requestBody, _ = json.Marshal(map[string]string{})
	req, err = http.NewRequest(method, urlLink, bytes.NewBuffer(requestBody))
	if err != nil {
		t.Fatal(err)
	}
	rr = httptest.NewRecorder()
	handler = http.HandlerFunc(likeHandlers.PostPostLikeHandler)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusInternalServerError {
		t.Errorf("PostPostLikeHandler returned wrong status code: got %v instead of %v",
			rr.Code, http.StatusInternalServerError)
	}
	expected = `{"message":"An error has occured in our server"}`
	if rr.Body.String() != expected {
		t.Errorf("PostPostLikeHandler returned unexpected body: got %v instead of %v",
			rr.Body.String(), expected)
	}

	// if checkIfPostExist return false
	likeCheckIfPostExistMock = func(s string) (bool, error) {
		return false, nil
	}
	requestBody, _ = json.Marshal(map[string]string{})
	req, err = http.NewRequest(method, urlLink, bytes.NewBuffer(requestBody))
	if err != nil {
		t.Fatal(err)
	}
	rr = httptest.NewRecorder()
	handler = http.HandlerFunc(likeHandlers.PostPostLikeHandler)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("PostPostLikeHandler returned wrong status code: got %v instead of %v",
			rr.Code, http.StatusBadRequest)
	}
	expected = `{"message":"Post does not exist"}`
	if rr.Body.String() != expected {
		t.Errorf("PostPostLikeHandler returned unexpected body: got %v instead of %v",
			rr.Body.String(), expected)
	}

	// if isLikeExist return error
	likeCheckIfPostExistMock = func(s string) (bool, error) {
		return true, nil
	}
	likeIsLikeExistMock = func(s1, s2, s3 string) (bool, error) {
		return false, errors.New("IsLikeExist return error")
	}
	requestBody, _ = json.Marshal(map[string]string{})
	req, err = http.NewRequest(method, urlLink, bytes.NewBuffer(requestBody))
	if err != nil {
		t.Fatal(err)
	}
	rr = httptest.NewRecorder()
	handler = http.HandlerFunc(likeHandlers.PostPostLikeHandler)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusInternalServerError {
		t.Errorf("PostPostLikeHandler returned wrong status code: got %v instead of %v",
			rr.Code, http.StatusInternalServerError)
	}
	expected = `{"message":"An error has occured in our server"}`
	if rr.Body.String() != expected {
		t.Errorf("PostPostLikeHandler returned unexpected body: got %v instead of %v",
			rr.Body.String(), expected)
	}

	// if isLikeExist return true
	likeIsLikeExistMock = func(s1, s2, s3 string) (bool, error) {
		return true, nil
	}
	requestBody, _ = json.Marshal(map[string]string{})
	req, err = http.NewRequest(method, urlLink, bytes.NewBuffer(requestBody))
	if err != nil {
		t.Fatal(err)
	}
	rr = httptest.NewRecorder()
	handler = http.HandlerFunc(likeHandlers.PostPostLikeHandler)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("PostPostLikeHandler returned wrong status code: got %v instead of %v",
			rr.Code, http.StatusBadRequest)
	}
	expected = `{"message":"User have already liked this post"}`
	if rr.Body.String() != expected {
		t.Errorf("PostPostLikeHandler returned unexpected body: got %v instead of %v",
			rr.Body.String(), expected)
	}

	// if insertLike return error
	likeIsLikeExistMock = func(s1, s2, s3 string) (bool, error) {
		return false, nil
	}
	likeInsertLikeMock = func(l models.Like) error {
		return errors.New("InsertLike returns error")
	}
	requestBody, _ = json.Marshal(map[string]string{})
	req, err = http.NewRequest(method, urlLink, bytes.NewBuffer(requestBody))
	if err != nil {
		t.Fatal(err)
	}
	rr = httptest.NewRecorder()
	handler = http.HandlerFunc(likeHandlers.PostPostLikeHandler)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusInternalServerError {
		t.Errorf("PostPostLikeHandler returned wrong status code: got %v instead of %v",
			rr.Code, http.StatusInternalServerError)
	}
	expected = `{"message":"An error has occured in our server"}`
	if rr.Body.String() != expected {
		t.Errorf("PostPostLikeHandler returned unexpected body: got %v instead of %v",
			rr.Body.String(), expected)
	}

	// if successful
	likeInsertLikeMock = func(l models.Like) error {
		return nil
	}
	requestBody, _ = json.Marshal(map[string]string{})
	req, err = http.NewRequest(method, urlLink, bytes.NewBuffer(requestBody))
	if err != nil {
		t.Fatal(err)
	}
	rr = httptest.NewRecorder()
	handler = http.HandlerFunc(likeHandlers.PostPostLikeHandler)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("PostPostLikeHandler returned wrong status code: got %v instead of %v",
			rr.Code, http.StatusCreated)
	}
}

func TestDeletePostLikeHandler(t *testing.T) {
	likeLikeServiceMock := newLikeLikeServiceMock()
	likeHandlerHeaderMock := newLikeHandlerHeaderMock()
	likeCommentServiceMock := newLikeCommentServiceMock()
	likePostServiceMock := newLikePostServiceMock()
	likeHandlers := NewLikeHandlers(likeLikeServiceMock, likePostServiceMock, likeCommentServiceMock, likeHandlerHeaderMock)
	method := "DELETE"
	urlLink := "http://localhost:8000/posts/post-566feabf-113f-4739-a226-9ba7cd6f4fc1/likes/like-5e80c8ca-4744-4c5c-b771-b0a79b9d07b7"
	userId := "user-566feabf-113f-4739-a226-9ba7cd6f4fc1"
	// if getUserIdToken return error
	likeGetUserIdFromTokenMock = func(s string) (string, error) {
		return "", errors.New("GetUserIdFromToken return error")
	}
	requestBody, _ := json.Marshal(map[string]string{})
	req, err := http.NewRequest(method, urlLink, bytes.NewBuffer(requestBody))
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(likeHandlers.DeletePostLikeHandler)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusInternalServerError {
		t.Errorf("DeletePostLikeHandler returned wrong status code: got %v instead of %v",
			rr.Code, http.StatusInternalServerError)
	}
	expected := `{"message":"An error has occured in our server"}`
	if rr.Body.String() != expected {
		t.Errorf("DeletePostLikeHandler returned unexpected body: got %v instead of %v",
			rr.Body.String(), expected)
	}

	// if isLikeExistById return error
	likeGetUserIdFromTokenMock = func(s string) (string, error) {
		return userId, nil
	}
	likeIsLikeExistByIdMock = func(s string) (bool, error) {
		return false, errors.New("IsLikeExistById return error")
	}
	requestBody, _ = json.Marshal(map[string]string{})
	req, err = http.NewRequest(method, urlLink, bytes.NewBuffer(requestBody))
	if err != nil {
		t.Fatal(err)
	}
	rr = httptest.NewRecorder()
	handler = http.HandlerFunc(likeHandlers.DeletePostLikeHandler)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusInternalServerError {
		t.Errorf("DeletePostLikeHandler returned wrong status code: got %v instead of %v",
			rr.Code, http.StatusInternalServerError)
	}
	expected = `{"message":"An error has occured in our server"}`
	if rr.Body.String() != expected {
		t.Errorf("DeletePostLikeHandler returned unexpected body: got %v instead of %v",
			rr.Body.String(), expected)
	}

	// if isLikeExistById return false
	likeIsLikeExistByIdMock = func(s string) (bool, error) {
		return false, nil
	}
	requestBody, _ = json.Marshal(map[string]string{})
	req, err = http.NewRequest(method, urlLink, bytes.NewBuffer(requestBody))
	if err != nil {
		t.Fatal(err)
	}
	rr = httptest.NewRecorder()
	handler = http.HandlerFunc(likeHandlers.DeletePostLikeHandler)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("DeletePostLikeHandler returned wrong status code: got %v instead of %v",
			rr.Code, http.StatusBadRequest)
	}
	expected = `{"message":"Like does not exist"}`
	if rr.Body.String() != expected {
		t.Errorf("DeletePostLikeHandler returned unexpected body: got %v instead of %v",
			rr.Body.String(), expected)
	}

	// if getLikeUserId return error
	likeIsLikeExistByIdMock = func(s string) (bool, error) {
		return true, nil
	}
	likeGetLikeUserIdMock = func(s string) (string, error) {
		return "", errors.New("GetLikeUserId return error")
	}
	requestBody, _ = json.Marshal(map[string]string{})
	req, err = http.NewRequest(method, urlLink, bytes.NewBuffer(requestBody))
	if err != nil {
		t.Fatal(err)
	}
	rr = httptest.NewRecorder()
	handler = http.HandlerFunc(likeHandlers.DeletePostLikeHandler)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusInternalServerError {
		t.Errorf("DeletePostLikeHandler returned wrong status code: got %v instead of %v",
			rr.Code, http.StatusInternalServerError)
	}
	expected = `{"message":"An error has occured in our server"}`
	if rr.Body.String() != expected {
		t.Errorf("DeletePostLikeHandler returned unexpected body: got %v instead of %v",
			rr.Body.String(), expected)
	}

	// if user id is not the same with getLikeUserId result
	likeGetLikeUserIdMock = func(s string) (string, error) {
		return "user-182ef588-bf73-424e-92d4-9da61af3c418", nil
	}
	requestBody, _ = json.Marshal(map[string]string{})
	req, err = http.NewRequest(method, urlLink, bytes.NewBuffer(requestBody))
	if err != nil {
		t.Fatal(err)
	}
	rr = httptest.NewRecorder()
	handler = http.HandlerFunc(likeHandlers.DeletePostLikeHandler)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusUnauthorized {
		t.Errorf("DeletePostLikeHandler returned wrong status code: got %v instead of %v",
			rr.Code, http.StatusUnauthorized)
	}
	expected = `{"message":"User is not authorized"}`
	if rr.Body.String() != expected {
		t.Errorf("DeletePostLikeHandler returned unexpected body: got %v instead of %v",
			rr.Body.String(), expected)
	}

	// if deleteLike return error
	likeGetLikeUserIdMock = func(s string) (string, error) {
		return userId, nil
	}
	likeDeleteLikeMock = func(s string) error {
		return errors.New("DeleteLike return error")
	}
	requestBody, _ = json.Marshal(map[string]string{})
	req, err = http.NewRequest(method, urlLink, bytes.NewBuffer(requestBody))
	if err != nil {
		t.Fatal(err)
	}
	rr = httptest.NewRecorder()
	handler = http.HandlerFunc(likeHandlers.DeletePostLikeHandler)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusInternalServerError {
		t.Errorf("DeletePostLikeHandler returned wrong status code: got %v instead of %v",
			rr.Code, http.StatusInternalServerError)
	}
	expected = `{"message":"An error has occured in our server"}`
	if rr.Body.String() != expected {
		t.Errorf("DeletePostLikeHandler returned unexpected body: got %v instead of %v",
			rr.Body.String(), expected)
	}

	// if successful
	likeDeleteLikeMock = func(s string) error {
		return nil
	}
	requestBody, _ = json.Marshal(map[string]string{})
	req, err = http.NewRequest(method, urlLink, bytes.NewBuffer(requestBody))
	if err != nil {
		t.Fatal(err)
	}
	rr = httptest.NewRecorder()
	handler = http.HandlerFunc(likeHandlers.DeletePostLikeHandler)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("DeletePostLikeHandler returned wrong status code: got %v instead of %v",
			rr.Code, http.StatusOK)
	}
}

func TestPostCommentLikeHandler(t *testing.T) {
	likeLikeServiceMock := newLikeLikeServiceMock()
	likeHandlerHeaderMock := newLikeHandlerHeaderMock()
	likeCommentServiceMock := newLikeCommentServiceMock()
	likePostServiceMock := newLikePostServiceMock()
	likeHandlers := NewLikeHandlers(likeLikeServiceMock, likePostServiceMock, likeCommentServiceMock, likeHandlerHeaderMock)
	method := "POST"
	urlLink := "http://localhost:8000/posts/post-edc48b68-cbdf-45aa-b063-c565fe7b6cac/comments/comment-9968bfd0-30de-4917-b57c-c2d30e765832/likes"
	userId := "user-566feabf-113f-4739-a226-9ba7cd6f4fc1"

	// if getUserIdToken return error
	likeGetUserIdFromTokenMock = func(s string) (string, error) {
		return "", errors.New("GetUserIdFromToken returns error")
	}
	requestBody, _ := json.Marshal(map[string]string{})
	req, err := http.NewRequest(method, urlLink, bytes.NewBuffer(requestBody))
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(likeHandlers.PostCommentLikeHandler)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusInternalServerError {
		t.Errorf("PostCommentLikeHandler returned wrong status code: got %v instead of %v",
			rr.Code, http.StatusInternalServerError)
	}
	expected := `{"message":"An error has occured in our server"}`
	if rr.Body.String() != expected {
		t.Errorf("PostCommentLikeHandler returned unexpected body: got %v instead of %v",
			rr.Body.String(), expected)
	}

	// if checkIfCommentExist return error
	likeGetUserIdFromTokenMock = func(s string) (string, error) {
		return userId, nil
	}
	likeCheckIfCommentExistMock = func(s string) (bool, error) {
		return false, errors.New("CheckIfCommentExist return error")
	}
	requestBody, _ = json.Marshal(map[string]string{})
	req, err = http.NewRequest(method, urlLink, bytes.NewBuffer(requestBody))
	if err != nil {
		t.Fatal(err)
	}
	rr = httptest.NewRecorder()
	handler = http.HandlerFunc(likeHandlers.PostCommentLikeHandler)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusInternalServerError {
		t.Errorf("PostCommentLikeHandler returned wrong status code: got %v instead of %v",
			rr.Code, http.StatusInternalServerError)
	}
	expected = `{"message":"An error has occured in our server"}`
	if rr.Body.String() != expected {
		t.Errorf("PostCommentLikeHandler returned unexpected body: got %v instead of %v",
			rr.Body.String(), expected)
	}

	// if checkIfCommentExist return false
	likeCheckIfCommentExistMock = func(s string) (bool, error) {
		return false, nil
	}
	requestBody, _ = json.Marshal(map[string]string{})
	req, err = http.NewRequest(method, urlLink, bytes.NewBuffer(requestBody))
	if err != nil {
		t.Fatal(err)
	}
	rr = httptest.NewRecorder()
	handler = http.HandlerFunc(likeHandlers.PostCommentLikeHandler)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("PostCommentLikeHandler returned wrong status code: got %v instead of %v",
			rr.Code, http.StatusBadRequest)
	}
	expected = `{"message":"Comment does not exist"}`
	if rr.Body.String() != expected {
		t.Errorf("PostCommentLikeHandler returned unexpected body: got %v instead of %v",
			rr.Body.String(), expected)
	}

	// if isLikeExist return error
	likeCheckIfCommentExistMock = func(s string) (bool, error) {
		return true, nil
	}
	likeIsLikeExistMock = func(s1, s2, s3 string) (bool, error) {
		return false, errors.New("IsLikeExist return error")
	}

	requestBody, _ = json.Marshal(map[string]string{})
	req, err = http.NewRequest(method, urlLink, bytes.NewBuffer(requestBody))
	if err != nil {
		t.Fatal(err)
	}
	rr = httptest.NewRecorder()
	handler = http.HandlerFunc(likeHandlers.PostCommentLikeHandler)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusInternalServerError {
		t.Errorf("PostCommentLikeHandler returned wrong status code: got %v instead of %v",
			rr.Code, http.StatusInternalServerError)
	}
	expected = `{"message":"An error has occured in our server"}`
	if rr.Body.String() != expected {
		t.Errorf("PostCommentLikeHandler returned unexpected body: got %v instead of %v",
			rr.Body.String(), expected)
	}

	// if isLikeExist return true
	likeIsLikeExistMock = func(s1, s2, s3 string) (bool, error) {
		return true, nil
	}
	requestBody, _ = json.Marshal(map[string]string{})
	req, err = http.NewRequest(method, urlLink, bytes.NewBuffer(requestBody))
	if err != nil {
		t.Fatal(err)
	}
	rr = httptest.NewRecorder()
	handler = http.HandlerFunc(likeHandlers.PostCommentLikeHandler)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("PostCommentLikeHandler returned wrong status code: got %v instead of %v",
			rr.Code, http.StatusBadRequest)
	}
	expected = `{"message":"User have already liked this comment"}`
	if rr.Body.String() != expected {
		t.Errorf("PostCommentLikeHandler returned unexpected body: got %v instead of %v",
			rr.Body.String(), expected)
	}

	// if insertLike return error
	likeIsLikeExistMock = func(s1, s2, s3 string) (bool, error) {
		return false, nil
	}
	likeInsertLikeMock = func(l models.Like) error {
		return errors.New("InsertLike return error")
	}
	requestBody, _ = json.Marshal(map[string]string{})
	req, err = http.NewRequest(method, urlLink, bytes.NewBuffer(requestBody))
	if err != nil {
		t.Fatal(err)
	}
	rr = httptest.NewRecorder()
	handler = http.HandlerFunc(likeHandlers.PostCommentLikeHandler)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusInternalServerError {
		t.Errorf("PostCommentLikeHandler returned wrong status code: got %v instead of %v",
			rr.Code, http.StatusInternalServerError)
	}
	expected = `{"message":"An error has occured in our server"}`
	if rr.Body.String() != expected {
		t.Errorf("PostCommentLikeHandler returned unexpected body: got %v instead of %v",
			rr.Body.String(), expected)
	}

	// if successful
	likeInsertLikeMock = func(l models.Like) error {
		return nil
	}
	requestBody, _ = json.Marshal(map[string]string{})
	req, err = http.NewRequest(method, urlLink, bytes.NewBuffer(requestBody))
	if err != nil {
		t.Fatal(err)
	}
	rr = httptest.NewRecorder()
	handler = http.HandlerFunc(likeHandlers.PostCommentLikeHandler)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("PostCommentLikeHandler returned wrong status code: got %v instead of %v",
			rr.Code, http.StatusCreated)
	}
}

func TestDeleteCommentLikeHandler(t *testing.T) {
	likeLikeServiceMock := newLikeLikeServiceMock()
	likeHandlerHeaderMock := newLikeHandlerHeaderMock()
	likeCommentServiceMock := newLikeCommentServiceMock()
	likePostServiceMock := newLikePostServiceMock()
	likeHandlers := NewLikeHandlers(likeLikeServiceMock, likePostServiceMock, likeCommentServiceMock, likeHandlerHeaderMock)
	method := "POST"
	urlLink := "http://localhost:8000/posts/post-edc48b68-cbdf-45aa-b063-c565fe7b6cac/comments/comment-9968bfd0-30de-4917-b57c-c2d30e765832/likes/like-81b1a77b-4ddd-43a7-a1be-3df506661676"
	userId := "user-566feabf-113f-4739-a226-9ba7cd6f4fc1"

	// if getUserIdToken return error
	likeGetUserIdFromTokenMock = func(s string) (string, error) {
		return "", errors.New("GetUserIdFromToken returns error")
	}
	requestBody, _ := json.Marshal(map[string]string{})
	req, err := http.NewRequest(method, urlLink, bytes.NewBuffer(requestBody))
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(likeHandlers.DeleteCommentLikeHandler)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusInternalServerError {
		t.Errorf("DeleteCommentLikeHandler returned wrong status code: got %v instead of %v",
			rr.Code, http.StatusInternalServerError)
	}
	expected := `{"message":"An error has occured in our server"}`
	if rr.Body.String() != expected {
		t.Errorf("DeleteCommentLikeHandler returned unexpected body: got %v instead of %v",
			rr.Body.String(), expected)
	}

	// if isLikeExistById return error
	likeGetUserIdFromTokenMock = func(s string) (string, error) {
		return userId, nil
	}
	likeIsLikeExistByIdMock = func(s string) (bool, error) {
		return false, errors.New("IsLikeExistById returns error")
	}
	requestBody, _ = json.Marshal(map[string]string{})
	req, err = http.NewRequest(method, urlLink, bytes.NewBuffer(requestBody))
	if err != nil {
		t.Fatal(err)
	}
	rr = httptest.NewRecorder()
	handler = http.HandlerFunc(likeHandlers.DeleteCommentLikeHandler)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusInternalServerError {
		t.Errorf("DeleteCommentLikeHandler returned wrong status code: got %v instead of %v",
			rr.Code, http.StatusInternalServerError)
	}
	expected = `{"message":"An error has occured in our server"}`
	if rr.Body.String() != expected {
		t.Errorf("DeleteCommentLikeHandler returned unexpected body: got %v instead of %v",
			rr.Body.String(), expected)
	}

	// if isLikeExistById return false
	likeIsLikeExistByIdMock = func(s string) (bool, error) {
		return false, nil
	}
	requestBody, _ = json.Marshal(map[string]string{})
	req, err = http.NewRequest(method, urlLink, bytes.NewBuffer(requestBody))
	if err != nil {
		t.Fatal(err)
	}
	rr = httptest.NewRecorder()
	handler = http.HandlerFunc(likeHandlers.DeleteCommentLikeHandler)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("DeleteCommentLikeHandler returned wrong status code: got %v instead of %v",
			rr.Code, http.StatusBadRequest)
	}
	expected = `{"message":"Like does not exist"}`
	if rr.Body.String() != expected {
		t.Errorf("DeleteCommentLikeHandler returned unexpected body: got %v instead of %v",
			rr.Body.String(), expected)
	}

	// if getLikeUserId return error
	likeIsLikeExistByIdMock = func(s string) (bool, error) {
		return true, nil
	}
	likeGetLikeUserIdMock = func(s string) (string, error) {
		return "", errors.New("GetLikeUserId return error")
	}
	requestBody, _ = json.Marshal(map[string]string{})
	req, err = http.NewRequest(method, urlLink, bytes.NewBuffer(requestBody))
	if err != nil {
		t.Fatal(err)
	}
	rr = httptest.NewRecorder()
	handler = http.HandlerFunc(likeHandlers.DeleteCommentLikeHandler)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusInternalServerError {
		t.Errorf("DeleteCommentLikeHandler returned wrong status code: got %v instead of %v",
			rr.Code, http.StatusInternalServerError)
	}
	expected = `{"message":"An error has occured in our server"}`
	if rr.Body.String() != expected {
		t.Errorf("DeleteCommentLikeHandler returned unexpected body: got %v instead of %v",
			rr.Body.String(), expected)
	}
	// if user id is not the same with getLikeUserId result
	likeGetLikeUserIdMock = func(s string) (string, error) {
		return "user-5e80c8ca-4744-4c5c-b771-b0a79b9d07b7", nil
	}
	requestBody, _ = json.Marshal(map[string]string{})
	req, err = http.NewRequest(method, urlLink, bytes.NewBuffer(requestBody))
	if err != nil {
		t.Fatal(err)
	}
	rr = httptest.NewRecorder()
	handler = http.HandlerFunc(likeHandlers.DeleteCommentLikeHandler)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusUnauthorized {
		t.Errorf("DeleteCommentLikeHandler returned wrong status code: got %v instead of %v",
			rr.Code, http.StatusUnauthorized)
	}
	expected = `{"message":"User is not authorized"}`
	if rr.Body.String() != expected {
		t.Errorf("DeleteCommentLikeHandler returned unexpected body: got %v instead of %v",
			rr.Body.String(), expected)
	}

	// if deleteLike return error
	likeGetLikeUserIdMock = func(s string) (string, error) {
		return userId, nil
	}
	likeDeleteLikeMock = func(s string) error {
		return errors.New("DeleteLike returns error")
	}
	requestBody, _ = json.Marshal(map[string]string{})
	req, err = http.NewRequest(method, urlLink, bytes.NewBuffer(requestBody))
	if err != nil {
		t.Fatal(err)
	}
	rr = httptest.NewRecorder()
	handler = http.HandlerFunc(likeHandlers.DeleteCommentLikeHandler)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusInternalServerError {
		t.Errorf("DeleteCommentLikeHandler returned wrong status code: got %v instead of %v",
			rr.Code, http.StatusInternalServerError)
	}
	expected = `{"message":"An error has occured in our server"}`
	if rr.Body.String() != expected {
		t.Errorf("DeleteCommentLikeHandler returned unexpected body: got %v instead of %v",
			rr.Body.String(), expected)
	}

	// if successful
	likeDeleteLikeMock = func(s string) error {
		return nil
	}
	requestBody, _ = json.Marshal(map[string]string{})
	req, err = http.NewRequest(method, urlLink, bytes.NewBuffer(requestBody))
	if err != nil {
		t.Fatal(err)
	}
	rr = httptest.NewRecorder()
	handler = http.HandlerFunc(likeHandlers.DeleteCommentLikeHandler)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("DeleteCommentLikeHandler returned wrong status code: got %v instead of %v",
			rr.Code, http.StatusOK)
	}
}
