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

func TestGetCommentsSuite(t *testing.T) {
	suite.Run(t, new(GetCommentsSuite))
}

func TestPostCommentSuite(t *testing.T) {
	suite.Run(t, new(PostCommentSuite))
}

func TestPutCommentSuite(t *testing.T) {
	suite.Run(t, new(PutCommentSuite))
}

func TestDeleteCommentSuite(t *testing.T) {
	suite.Run(t, new(DeleteCommentSuite))
}

type GetCommentsSuite struct {
	suite.Suite
	commentUsecase *mocks.CommentUsecase
}

func (gcs *GetCommentsSuite) SetupTest() {
	gcs.commentUsecase = new(mocks.CommentUsecase)
}

func (gcs *GetCommentsSuite) TestFindCommentsError() {
	gcs.commentUsecase.On("FindComments", mock.Anything).Return(nil, domain.ErrInternalServerError)
	commentHandler := commentHttp.NewCommentHandler(gcs.commentUsecase)
	req, _ := http.NewRequest("GET", "/posts/postid1/comments", nil)
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(commentHandler.Comments)
	handler.ServeHTTP(rr, req)

	assert.Equalf(gcs.T(), http.StatusInternalServerError, rr.Code, "Should have responded with http status code %v but got %v", http.StatusInternalServerError, rr.Code)
	expectedBody := `{"message":"` + domain.ErrInternalServerError.Error() + `"}`
	assert.Equalf(gcs.T(), expectedBody, rr.Body.String(), "Should have responded with body %s but got %s", expectedBody, rr.Body.String())
}

func (gcs *GetCommentsSuite) TestGetCommentsSuccessful() {
	gcs.commentUsecase.On("FindComments", mock.Anything).Return(&[]domain.Comment{}, nil)
	commentHandler := commentHttp.NewCommentHandler(gcs.commentUsecase)
	req, _ := http.NewRequest("GET", "/posts/postid1/comments", nil)
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(commentHandler.Comments)
	handler.ServeHTTP(rr, req)

	assert.Equalf(gcs.T(), http.StatusOK, rr.Code, "Should have responded with http status code %s but got %s", http.StatusOK, rr.Code)
}

type PostCommentSuite struct {
	suite.Suite
	commentUsecase *mocks.CommentUsecase
}

func (pcs *PostCommentSuite) SetupTest() {
	pcs.commentUsecase = new(mocks.CommentUsecase)
}

func (pcs *PostCommentSuite) TestCommentNotProvided() {
	requestBody, _ := json.Marshal(map[string]string{
		"comment": "",
	})
	commentHandler := commentHttp.NewCommentHandler(pcs.commentUsecase)
	req, _ := http.NewRequest("POST", "/posts/postid1/comments", bytes.NewBuffer(requestBody))
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(commentHandler.Comments)
	handler.ServeHTTP(rr, req)

	assert.Equalf(pcs.T(), http.StatusBadRequest, rr.Code, "Should have responded with http status code %v but got %v", http.StatusBadRequest, rr.Code)
	expectedBody := `{"message":"` + domain.ErrMissingCommentInput.Error() + `"}`
	assert.Equalf(pcs.T(), expectedBody, rr.Body.String(), "Should have responded with body %s but got %s", expectedBody, rr.Body.String())
}

func (pcs *PostCommentSuite) TestPostCommentError() {
	requestBody, _ := json.Marshal(map[string]string{
		"comment": "a new comment",
	})
	pcs.commentUsecase.On("PostComment", mock.Anything, mock.Anything).Return(domain.ErrInternalServerError)
	commentHandler := commentHttp.NewCommentHandler(pcs.commentUsecase)
	req, _ := http.NewRequest("POST", "/posts/postid1/comments", bytes.NewBuffer(requestBody))
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(commentHandler.Comments)
	handler.ServeHTTP(rr, req)

	assert.Equalf(pcs.T(), http.StatusInternalServerError, rr.Code, "Should have responded with http status code %v but got %v", http.StatusInternalServerError, rr.Code)
	expectedBody := `{"message":"` + domain.ErrInternalServerError.Error() + `"}`
	assert.Equalf(pcs.T(), expectedBody, rr.Body.String(), "Should have responded with body %s but got %s", expectedBody, rr.Body.String())
}

func (pcs *PostCommentSuite) TestPostCommentSuccessful() {
	requestBody, _ := json.Marshal(map[string]string{
		"comment": "a new comment",
	})

	pcs.commentUsecase.On("PostComment", mock.Anything, mock.Anything).Return(nil)
	commentHandler := commentHttp.NewCommentHandler(pcs.commentUsecase)
	req, _ := http.NewRequest("POST", "/posts/postid1/comments", bytes.NewBuffer(requestBody))
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(commentHandler.Comments)
	handler.ServeHTTP(rr, req)

	assert.Equalf(pcs.T(), http.StatusCreated, rr.Code, "Should have responded with http status code %v but got %v", http.StatusCreated, rr.Code)
	expectedBody := `{"message":"Comment successfully Created"}`
	assert.Equalf(pcs.T(), expectedBody, rr.Body.String(), "Should have responded with body %s but got %s", expectedBody, rr.Body.String())
}

type PutCommentSuite struct {
	suite.Suite
	commentUsecase *mocks.CommentUsecase
}

func (pcs *PutCommentSuite) SetupTest() {
	pcs.commentUsecase = new(mocks.CommentUsecase)
}

func (pcs *PutCommentSuite) TestCommentNotProvided() {
	requestBody, _ := json.Marshal(map[string]string{
		"comment": "",
	})

	commentHandler := commentHttp.NewCommentHandler(pcs.commentUsecase)
	req, _ := http.NewRequest("PUT", "/posts/postid1/comments/commentid1", bytes.NewBuffer(requestBody))
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(commentHandler.Comment)
	handler.ServeHTTP(rr, req)

	assert.Equalf(pcs.T(), http.StatusBadRequest, rr.Code, "Should have responded with http status code %v but got %v", http.StatusBadRequest, rr.Code)
	expectedBody := `{"message":"` + domain.ErrMissingCommentInput.Error() + `"}`
	assert.Equalf(pcs.T(), expectedBody, rr.Body.String(), "Should have responded with body %s but got %s", expectedBody, rr.Body.String())
}

func (pcs *PutCommentSuite) TestPutCommentError() {
	requestBody, _ := json.Marshal(map[string]string{
		"comment": "a new comment",
	})
	pcs.commentUsecase.On("PutComment", mock.Anything, mock.Anything).Return(domain.ErrInternalServerError)
	commentHandler := commentHttp.NewCommentHandler(pcs.commentUsecase)
	req, _ := http.NewRequest("PUT", "/posts/postid1/comments/commentid1", bytes.NewBuffer(requestBody))
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(commentHandler.Comment)
	handler.ServeHTTP(rr, req)

	assert.Equalf(pcs.T(), http.StatusInternalServerError, rr.Code, "Should have responded with http status code %v but got %v", http.StatusInternalServerError, rr.Code)
	expectedBody := `{"message":"` + domain.ErrInternalServerError.Error() + `"}`
	assert.Equalf(pcs.T(), expectedBody, rr.Body.String(), "Should have responded with body %s but got %s", expectedBody, rr.Body.String())
}

func (pcs *PutCommentSuite) TestPutCommentSuccessful() {
	requestBody, _ := json.Marshal(map[string]string{
		"comment": "a new comment",
	})
	pcs.commentUsecase.On("PutComment", mock.Anything, mock.Anything).Return(nil)
	commentHandler := commentHttp.NewCommentHandler(pcs.commentUsecase)
	req, _ := http.NewRequest("PUT", "/posts/postid1/comments/commentid1", bytes.NewBuffer(requestBody))
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(commentHandler.Comment)
	handler.ServeHTTP(rr, req)

	assert.Equalf(pcs.T(), http.StatusOK, rr.Code, "Should have responded with http status code %v but got %v", http.StatusOK, rr.Code)
	expectedBody := `{"message":"Comment successfully Updated"}`
	assert.Equalf(pcs.T(), expectedBody, rr.Body.String(), "Should have responded with body %s but got %s", expectedBody, rr.Body.String())

}

type DeleteCommentSuite struct {
	suite.Suite
	commentUsecase *mocks.CommentUsecase
}

func (dcs *DeleteCommentSuite) SetupTest() {
	dcs.commentUsecase = new(mocks.CommentUsecase)
}

func (dcs *DeleteCommentSuite) TestDeleteCommentError() {
	dcs.commentUsecase.On("DeleteComment", mock.Anything, mock.Anything).Return(domain.ErrInternalServerError)
	commentHandler := commentHttp.NewCommentHandler(dcs.commentUsecase)
	req, _ := http.NewRequest("DELETE", "/posts/postid1/comments/commentid1", nil)
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(commentHandler.Comment)
	handler.ServeHTTP(rr, req)

	assert.Equalf(dcs.T(), http.StatusInternalServerError, rr.Code, "Should have responded with http status code %v but got %v", http.StatusInternalServerError, rr.Code)
	expectedBody := `{"message":"` + domain.ErrInternalServerError.Error() + `"}`
	assert.Equalf(dcs.T(), expectedBody, rr.Body.String(), "Should have responded with body %s but got %s", expectedBody, rr.Body.String())
}

func (dcs *DeleteCommentSuite) TestDeleteCommentSuccessful() {
	dcs.commentUsecase.On("DeleteComment", mock.Anything, mock.Anything).Return(nil)
	commentHandler := commentHttp.NewCommentHandler(dcs.commentUsecase)
	req, _ := http.NewRequest("DELETE", "/posts/postid1/comments/commentid1", nil)
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(commentHandler.Comment)
	handler.ServeHTTP(rr, req)

	assert.Equalf(dcs.T(), http.StatusOK, rr.Code, "Should have responded with http status code %v but got %v", http.StatusOK, rr.Code)
	expectedBody := `{"message":"Comment successfully Deleted"}`
	assert.Equalf(dcs.T(), expectedBody, rr.Body.String(), "Should have responded with body %s but got %s", expectedBody, rr.Body.String())
}
