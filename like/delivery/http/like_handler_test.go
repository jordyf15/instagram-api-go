package http_test

import (
	"instagram-go/domain"
	"instagram-go/domain/mocks"
	likeHttp "instagram-go/like/delivery/http"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

func TestLikeHandlerSuite(t *testing.T) {
	suite.Run(t, new(LikeHandlerSuite))
}

type LikeHandlerSuite struct {
	suite.Suite
	likeUsecase *mocks.LikeUsecase
}

func (lh *LikeHandlerSuite) SetupTest() {
	lh.likeUsecase = new(mocks.LikeUsecase)
}

func (lh *LikeHandlerSuite) TestPostLikePostInsertPostLikeError() {
	lh.likeUsecase.On("InsertPostLike", mock.AnythingOfType("string"), mock.AnythingOfType("string")).Return(domain.ErrInternalServerError)
	likeHandler := likeHttp.NewLikeHandler(lh.likeUsecase)
	req, _ := http.NewRequest("POST", "/posts/postid1/likes", nil)
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(likeHandler.PostLikePost)
	handler.ServeHTTP(rr, req)

	assert.Equalf(lh.T(), http.StatusInternalServerError, rr.Code, "Should have responded with http status code %v but got %v", http.StatusInternalServerError, rr.Code)
	expectedBody := `{"message":"` + domain.ErrInternalServerError.Error() + `"}`
	assert.Equalf(lh.T(), expectedBody, rr.Body.String(), "Should have responded with body %s but got %s", expectedBody, rr.Body.String())
}

func (lh *LikeHandlerSuite) TestPostLikePostSuccessful() {
	lh.likeUsecase.On("InsertPostLike", mock.AnythingOfType("string"), mock.AnythingOfType("string")).Return(nil)
	likeHandler := likeHttp.NewLikeHandler(lh.likeUsecase)
	req, _ := http.NewRequest("POST", "/posts/postid1/likes", nil)
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(likeHandler.PostLikePost)
	handler.ServeHTTP(rr, req)

	assert.Equalf(lh.T(), http.StatusCreated, rr.Code, "Should have responded with http status code %v but got %v", http.StatusCreated, rr.Code)
}

func (lh *LikeHandlerSuite) TestDeleteLikePostDeletePostLikeError() {
	lh.likeUsecase.On("DeletePostLike", mock.AnythingOfType("string"), mock.AnythingOfType("string")).Return(domain.ErrInternalServerError)
	likeHandler := likeHttp.NewLikeHandler(lh.likeUsecase)
	req, _ := http.NewRequest("DELETE", "/posts/postid1/likes/likeid1", nil)
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(likeHandler.DeleteLikePost)
	handler.ServeHTTP(rr, req)

	assert.Equalf(lh.T(), http.StatusInternalServerError, rr.Code, "Should have responded with http status code %v but got %v", http.StatusInternalServerError, rr.Code)
	expectedBody := `{"message":"` + domain.ErrInternalServerError.Error() + `"}`
	assert.Equalf(lh.T(), expectedBody, rr.Body.String(), "Should have responded with body %s but got %s", expectedBody, rr.Body.String())
}

func (lh *LikeHandlerSuite) TestDeleteLikePostSuccessful() {
	lh.likeUsecase.On("DeletePostLike", mock.AnythingOfType("string"), mock.AnythingOfType("string")).Return(nil)
	likeHandler := likeHttp.NewLikeHandler(lh.likeUsecase)
	req, _ := http.NewRequest("DELETE", "/posts/postid1/likes/likeid1", nil)
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(likeHandler.DeleteLikePost)
	handler.ServeHTTP(rr, req)
	assert.Equalf(lh.T(), http.StatusOK, rr.Code, "Should have responded with http status code %v but got %v", http.StatusOK, rr.Code)
}

func (lh *LikeHandlerSuite) TestPostCommentLikeInsertCommentLikeError() {
	lh.likeUsecase.On("InsertCommentLike", mock.AnythingOfType("string"), mock.AnythingOfType("string")).Return(domain.ErrInternalServerError)
	likeHandler := likeHttp.NewLikeHandler(lh.likeUsecase)
	req, _ := http.NewRequest("POST", "/posts/postid1/comments/commentid1/likes", nil)
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(likeHandler.PostCommentLike)
	handler.ServeHTTP(rr, req)

	assert.Equalf(lh.T(), http.StatusInternalServerError, rr.Code, "Should have responded with http status code %v but got %v", http.StatusInternalServerError, rr.Code)
	expectedBody := `{"message":"` + domain.ErrInternalServerError.Error() + `"}`
	assert.Equalf(lh.T(), expectedBody, rr.Body.String(), "Should have responded with body %s but got %s", expectedBody, rr.Body.String())
}

func (lh *LikeHandlerSuite) TestPostCommentLikeSuccessful() {
	lh.likeUsecase.On("InsertCommentLike", mock.AnythingOfType("string"), mock.AnythingOfType("string")).Return(nil)
	likeHandler := likeHttp.NewLikeHandler(lh.likeUsecase)
	req, _ := http.NewRequest("POST", "/posts/postid1/comments/commentid1/likes", nil)
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(likeHandler.PostCommentLike)
	handler.ServeHTTP(rr, req)

	assert.Equalf(lh.T(), http.StatusCreated, rr.Code, "Should have responded with http status code %v but got %v", http.StatusCreated, rr.Code)
}

func (lh *LikeHandlerSuite) TestDeleteCommentLikeDeleteCommentLikeError() {
	lh.likeUsecase.On("DeleteCommentLike", mock.AnythingOfType("string"), mock.AnythingOfType("string")).Return(domain.ErrInternalServerError)
	likeHandler := likeHttp.NewLikeHandler(lh.likeUsecase)
	req, _ := http.NewRequest("DELETE", "/posts/postid1/comments/commentid1/likes/deleteid1", nil)
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(likeHandler.DeleteCommentLike)
	handler.ServeHTTP(rr, req)

	assert.Equalf(lh.T(), http.StatusInternalServerError, rr.Code, "Should have responded with http status code %v but got %v", http.StatusInternalServerError, rr.Code)
	expectedBody := `{"message":"` + domain.ErrInternalServerError.Error() + `"}`
	assert.Equalf(lh.T(), expectedBody, rr.Body.String(), "Should have responded with body %s but got %s", expectedBody, rr.Body.String())
}

func (lh *LikeHandlerSuite) TestDeleteCommentLikeSuccessful() {
	lh.likeUsecase.On("DeleteCommentLike", mock.AnythingOfType("string"), mock.AnythingOfType("string")).Return(nil)
	likeHandler := likeHttp.NewLikeHandler(lh.likeUsecase)
	req, _ := http.NewRequest("DELETE", "/posts/postid1/comments/commentid1/likes/deleteid1", nil)
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(likeHandler.DeleteCommentLike)
	handler.ServeHTTP(rr, req)

	assert.Equalf(lh.T(), http.StatusOK, rr.Code, "Should have responded with http status code %v but got %v", http.StatusOK, rr.Code)
}
