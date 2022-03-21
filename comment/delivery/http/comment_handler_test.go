package http_test

import (
	"bytes"
	"encoding/json"
	commentHttp "instagram-go/comment/delivery/http"
	"instagram-go/domain"
	"instagram-go/domain/mocks"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

func TestCommentHandlerSuite(t *testing.T) {
	suite.Run(t, new(CommentHandlerSuite))
}

type CommentHandlerSuite struct {
	suite.Suite
	commentUsecase *mocks.CommentUsecase
}

func (ch *CommentHandlerSuite) SetupTest() {
	ch.commentUsecase = new(mocks.CommentUsecase)
}

func (ch *CommentHandlerSuite) TestGetCommentsFindCommentsError() {
	ch.commentUsecase.On("FindComments", mock.AnythingOfType("string")).Return(nil, domain.ErrInternalServerError)
	commentHandler := commentHttp.NewCommentHandler(ch.commentUsecase)
	req, _ := http.NewRequest("GET", "/posts/postid1/comments", nil)
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(commentHandler.Comments)
	handler.ServeHTTP(rr, req)

	assert.Equalf(ch.T(), http.StatusInternalServerError, rr.Code, "Should have responded with http status code %v but got %v", http.StatusInternalServerError, rr.Code)
	expectedBody := `{"message":"` + domain.ErrInternalServerError.Error() + `"}`
	assert.Equalf(ch.T(), expectedBody, rr.Body.String(), "Should have responded with body %s but got %s", expectedBody, rr.Body.String())
}

func (ch *CommentHandlerSuite) TestGetCommentsSuccessful() {
	ch.commentUsecase.On("FindComments", mock.AnythingOfType("string")).Return(&[]domain.Comment{}, nil)
	commentHandler := commentHttp.NewCommentHandler(ch.commentUsecase)
	req, _ := http.NewRequest("GET", "/posts/postid1/comments", nil)
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(commentHandler.Comments)
	handler.ServeHTTP(rr, req)

	assert.Equalf(ch.T(), http.StatusOK, rr.Code, "Should have responded with http status code %s but got %s", http.StatusOK, rr.Code)
}

func (ch *CommentHandlerSuite) TestPostCommentCommentNotProvided() {
	requestBody, _ := json.Marshal(map[string]string{
		"comment": "",
	})
	commentHandler := commentHttp.NewCommentHandler(ch.commentUsecase)
	req, _ := http.NewRequest("POST", "/posts/postid1/comments", bytes.NewBuffer(requestBody))
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(commentHandler.Comments)
	handler.ServeHTTP(rr, req)

	assert.Equalf(ch.T(), http.StatusBadRequest, rr.Code, "Should have responded with http status code %v but got %v", http.StatusBadRequest, rr.Code)
	expectedBody := `{"message":"` + domain.ErrMissingCommentInput.Error() + `"}`
	assert.Equalf(ch.T(), expectedBody, rr.Body.String(), "Should have responded with body %s but got %s", expectedBody, rr.Body.String())
}

func (ch *CommentHandlerSuite) TestPostCommentPostCommentError() {
	requestBody, _ := json.Marshal(map[string]string{
		"comment": "a new comment",
	})
	ch.commentUsecase.On("PostComment", mock.AnythingOfType("*domain.Comment"), mock.AnythingOfType("string")).Return(domain.ErrInternalServerError)
	commentHandler := commentHttp.NewCommentHandler(ch.commentUsecase)
	req, _ := http.NewRequest("POST", "/posts/postid1/comments", bytes.NewBuffer(requestBody))
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(commentHandler.Comments)
	handler.ServeHTTP(rr, req)

	assert.Equalf(ch.T(), http.StatusInternalServerError, rr.Code, "Should have responded with http status code %v but got %v", http.StatusInternalServerError, rr.Code)
	expectedBody := `{"message":"` + domain.ErrInternalServerError.Error() + `"}`
	assert.Equalf(ch.T(), expectedBody, rr.Body.String(), "Should have responded with body %s but got %s", expectedBody, rr.Body.String())
}

func (ch *CommentHandlerSuite) TestPostCommentSuccessful() {
	requestBody, _ := json.Marshal(map[string]string{
		"comment": "a new comment",
	})

	ch.commentUsecase.On("PostComment", mock.AnythingOfType("*domain.Comment"), mock.AnythingOfType("string")).Return(nil)
	commentHandler := commentHttp.NewCommentHandler(ch.commentUsecase)
	req, _ := http.NewRequest("POST", "/posts/postid1/comments", bytes.NewBuffer(requestBody))
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(commentHandler.Comments)
	handler.ServeHTTP(rr, req)

	assert.Equalf(ch.T(), http.StatusCreated, rr.Code, "Should have responded with http status code %v but got %v", http.StatusCreated, rr.Code)
	expectedBody := `{"message":"Comment successfully Created"}`
	assert.Equalf(ch.T(), expectedBody, rr.Body.String(), "Should have responded with body %s but got %s", expectedBody, rr.Body.String())
}

func (ch *CommentHandlerSuite) TestPutCommentCommentNotProvided() {
	requestBody, _ := json.Marshal(map[string]string{
		"comment": "",
	})

	commentHandler := commentHttp.NewCommentHandler(ch.commentUsecase)
	req, _ := http.NewRequest("PUT", "/posts/postid1/comments/commentid1", bytes.NewBuffer(requestBody))
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(commentHandler.Comment)
	handler.ServeHTTP(rr, req)

	assert.Equalf(ch.T(), http.StatusBadRequest, rr.Code, "Should have responded with http status code %v but got %v", http.StatusBadRequest, rr.Code)
	expectedBody := `{"message":"` + domain.ErrMissingCommentInput.Error() + `"}`
	assert.Equalf(ch.T(), expectedBody, rr.Body.String(), "Should have responded with body %s but got %s", expectedBody, rr.Body.String())
}

func (ch *CommentHandlerSuite) TestPutCommentPutCommentError() {
	requestBody, _ := json.Marshal(map[string]string{
		"comment": "a new comment",
	})
	ch.commentUsecase.On("PutComment", mock.AnythingOfType("*domain.Comment"), mock.AnythingOfType("string")).Return(domain.ErrInternalServerError)
	commentHandler := commentHttp.NewCommentHandler(ch.commentUsecase)
	req, _ := http.NewRequest("PUT", "/posts/postid1/comments/commentid1", bytes.NewBuffer(requestBody))
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(commentHandler.Comment)
	handler.ServeHTTP(rr, req)

	assert.Equalf(ch.T(), http.StatusInternalServerError, rr.Code, "Should have responded with http status code %v but got %v", http.StatusInternalServerError, rr.Code)
	expectedBody := `{"message":"` + domain.ErrInternalServerError.Error() + `"}`
	assert.Equalf(ch.T(), expectedBody, rr.Body.String(), "Should have responded with body %s but got %s", expectedBody, rr.Body.String())
}

func (ch *CommentHandlerSuite) TestPutCommentSuccessful() {
	requestBody, _ := json.Marshal(map[string]string{
		"comment": "a new comment",
	})
	ch.commentUsecase.On("PutComment", mock.AnythingOfType("*domain.Comment"), mock.AnythingOfType("string")).Return(nil)
	commentHandler := commentHttp.NewCommentHandler(ch.commentUsecase)
	req, _ := http.NewRequest("PUT", "/posts/postid1/comments/commentid1", bytes.NewBuffer(requestBody))
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(commentHandler.Comment)
	handler.ServeHTTP(rr, req)

	assert.Equalf(ch.T(), http.StatusOK, rr.Code, "Should have responded with http status code %v but got %v", http.StatusOK, rr.Code)
	expectedBody := `{"message":"Comment successfully Updated"}`
	assert.Equalf(ch.T(), expectedBody, rr.Body.String(), "Should have responded with body %s but got %s", expectedBody, rr.Body.String())

}

type DeleteCommentSuite struct {
	suite.Suite
	commentUsecase *mocks.CommentUsecase
}

func (dcs *DeleteCommentSuite) SetupTest() {
	dcs.commentUsecase = new(mocks.CommentUsecase)
}

func (ch *CommentHandlerSuite) TestDeleteCommentDeleteCommentError() {
	ch.commentUsecase.On("DeleteComment", mock.AnythingOfType("string"), mock.AnythingOfType("string")).Return(domain.ErrInternalServerError)
	commentHandler := commentHttp.NewCommentHandler(ch.commentUsecase)
	req, _ := http.NewRequest("DELETE", "/posts/postid1/comments/commentid1", nil)
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(commentHandler.Comment)
	handler.ServeHTTP(rr, req)

	assert.Equalf(ch.T(), http.StatusInternalServerError, rr.Code, "Should have responded with http status code %v but got %v", http.StatusInternalServerError, rr.Code)
	expectedBody := `{"message":"` + domain.ErrInternalServerError.Error() + `"}`
	assert.Equalf(ch.T(), expectedBody, rr.Body.String(), "Should have responded with body %s but got %s", expectedBody, rr.Body.String())
}

func (ch *DeleteCommentSuite) TestDeleteCommentSuccessful() {
	ch.commentUsecase.On("DeleteComment", mock.AnythingOfType("string"), mock.AnythingOfType("string")).Return(nil)
	commentHandler := commentHttp.NewCommentHandler(ch.commentUsecase)
	req, _ := http.NewRequest("DELETE", "/posts/postid1/comments/commentid1", nil)
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(commentHandler.Comment)
	handler.ServeHTTP(rr, req)

	assert.Equalf(ch.T(), http.StatusOK, rr.Code, "Should have responded with http status code %v but got %v", http.StatusOK, rr.Code)
	expectedBody := `{"message":"Comment successfully Deleted"}`
	assert.Equalf(ch.T(), expectedBody, rr.Body.String(), "Should have responded with body %s but got %s", expectedBody, rr.Body.String())
}
