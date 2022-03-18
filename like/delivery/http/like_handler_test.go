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

func TestPostLikePostSuite(t *testing.T) {
	suite.Run(t, new(PostLikePostSuite))
}

func TestDeleteLikePostSuite(t *testing.T) {
	suite.Run(t, new(DeleteLikePostSuite))
}
func TestPostCommentLikeSuite(t *testing.T) {
	suite.Run(t, new(PostCommentLikeSuite))
}
func TestDeleteCommentLikeSuite(t *testing.T) {
	suite.Run(t, new(DeleteCommentLikeSuite))
}

type PostLikePostSuite struct {
	suite.Suite
	likeUsecase *mocks.LikeUsecase
}

func (plps *PostLikePostSuite) SetupTest() {
	plps.likeUsecase = new(mocks.LikeUsecase)
}

func (plps *PostLikePostSuite) TestInsertPostLikeError() {
	plps.likeUsecase.On("InsertPostLike", mock.Anything, mock.Anything).Return(domain.ErrInternalServerError)
	likeHandler := likeHttp.NewLikeHandler(plps.likeUsecase)
	req, _ := http.NewRequest("POST", "/posts/postid1/likes", nil)
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(likeHandler.PostLikePost)
	handler.ServeHTTP(rr, req)

	assert.Equalf(plps.T(), http.StatusInternalServerError, rr.Code, "Should have responded with http status code %v but got %v", http.StatusInternalServerError, rr.Code)
	expectedBody := `{"message":"` + domain.ErrInternalServerError.Error() + `"}`
	assert.Equalf(plps.T(), expectedBody, rr.Body.String(), "Should have responded with body %s but got %s", expectedBody, rr.Body.String())
}

func (plps *PostLikePostSuite) TestPostLikePostSuccessful() {
	plps.likeUsecase.On("InsertPostLike", mock.Anything, mock.Anything).Return(nil)
	likeHandler := likeHttp.NewLikeHandler(plps.likeUsecase)
	req, _ := http.NewRequest("POST", "/posts/postid1/likes", nil)
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(likeHandler.PostLikePost)
	handler.ServeHTTP(rr, req)

	assert.Equalf(plps.T(), http.StatusCreated, rr.Code, "Should have responded with http status code %v but got %v", http.StatusCreated, rr.Code)
}

type DeleteLikePostSuite struct {
	suite.Suite
	likeUsecase *mocks.LikeUsecase
}

func (dlps *DeleteLikePostSuite) SetupTest() {
	dlps.likeUsecase = new(mocks.LikeUsecase)
}
func (dlps *DeleteLikePostSuite) TestDeletePostLikeError() {
	dlps.likeUsecase.On("DeletePostLike", mock.Anything, mock.Anything).Return(domain.ErrInternalServerError)
	likeHandler := likeHttp.NewLikeHandler(dlps.likeUsecase)
	req, _ := http.NewRequest("DELETE", "/posts/postid1/likes/likeid1", nil)
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(likeHandler.DeleteLikePost)
	handler.ServeHTTP(rr, req)

	assert.Equalf(dlps.T(), http.StatusInternalServerError, rr.Code, "Should have responded with http status code %v but got %v", http.StatusInternalServerError, rr.Code)
	expectedBody := `{"message":"` + domain.ErrInternalServerError.Error() + `"}`
	assert.Equalf(dlps.T(), expectedBody, rr.Body.String(), "Should have responded with body %s but got %s", expectedBody, rr.Body.String())
}

func (dlps *DeleteLikePostSuite) TestDeleteLikePostSuccessful() {
	dlps.likeUsecase.On("DeletePostLike", mock.Anything, mock.Anything).Return(nil)
	likeHandler := likeHttp.NewLikeHandler(dlps.likeUsecase)
	req, _ := http.NewRequest("DELETE", "/posts/postid1/likes/likeid1", nil)
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(likeHandler.DeleteLikePost)
	handler.ServeHTTP(rr, req)
	assert.Equalf(dlps.T(), http.StatusOK, rr.Code, "Should have responded with http status code %v but got %v", http.StatusOK, rr.Code)
}

type PostCommentLikeSuite struct {
	suite.Suite
	likeUsecase *mocks.LikeUsecase
}

func (pcls *PostCommentLikeSuite) SetupTest() {
	pcls.likeUsecase = new(mocks.LikeUsecase)
}

func (pcls *PostCommentLikeSuite) TestInsertCommentLikeError() {
	pcls.likeUsecase.On("InsertCommentLike", mock.Anything, mock.Anything).Return(domain.ErrInternalServerError)
	likeHandler := likeHttp.NewLikeHandler(pcls.likeUsecase)
	req, _ := http.NewRequest("POST", "/posts/postid1/comments/commentid1/likes", nil)
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(likeHandler.PostCommentLike)
	handler.ServeHTTP(rr, req)

	assert.Equalf(pcls.T(), http.StatusInternalServerError, rr.Code, "Should have responded with http status code %v but got %v", http.StatusInternalServerError, rr.Code)
	expectedBody := `{"message":"` + domain.ErrInternalServerError.Error() + `"}`
	assert.Equalf(pcls.T(), expectedBody, rr.Body.String(), "Should have responded with body %s but got %s", expectedBody, rr.Body.String())
}

func (pcls *PostCommentLikeSuite) TestPostCommentLikeSuccessful() {
	pcls.likeUsecase.On("InsertCommentLike", mock.Anything, mock.Anything).Return(nil)
	likeHandler := likeHttp.NewLikeHandler(pcls.likeUsecase)
	req, _ := http.NewRequest("POST", "/posts/postid1/comments/commentid1/likes", nil)
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(likeHandler.PostCommentLike)
	handler.ServeHTTP(rr, req)

	assert.Equalf(pcls.T(), http.StatusCreated, rr.Code, "Should have responded with http status code %v but got %v", http.StatusCreated, rr.Code)
}

type DeleteCommentLikeSuite struct {
	suite.Suite
	likeUsecase *mocks.LikeUsecase
}

func (dcls *DeleteCommentLikeSuite) SetupTest() {
	dcls.likeUsecase = new(mocks.LikeUsecase)
}

func (dcls *DeleteCommentLikeSuite) TestDeleteCommentLikeError() {
	dcls.likeUsecase.On("DeleteCommentLike", mock.Anything, mock.Anything).Return(domain.ErrInternalServerError)
	likeHandler := likeHttp.NewLikeHandler(dcls.likeUsecase)
	req, _ := http.NewRequest("DELETE", "/posts/postid1/comments/commentid1/likes/deleteid1", nil)
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(likeHandler.DeleteCommentLike)
	handler.ServeHTTP(rr, req)

	assert.Equalf(dcls.T(), http.StatusInternalServerError, rr.Code, "Should have responded with http status code %v but got %v", http.StatusInternalServerError, rr.Code)
	expectedBody := `{"message":"` + domain.ErrInternalServerError.Error() + `"}`
	assert.Equalf(dcls.T(), expectedBody, rr.Body.String(), "Should have responded with body %s but got %s", expectedBody, rr.Body.String())
}

func (dcls *DeleteCommentLikeSuite) TestDeleteCommentLikeSuccessful() {
	dcls.likeUsecase.On("DeleteCommentLike", mock.Anything, mock.Anything).Return(nil)
	likeHandler := likeHttp.NewLikeHandler(dcls.likeUsecase)
	req, _ := http.NewRequest("DELETE", "/posts/postid1/comments/commentid1/likes/deleteid1", nil)
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(likeHandler.DeleteCommentLike)
	handler.ServeHTTP(rr, req)

	assert.Equalf(dcls.T(), http.StatusOK, rr.Code, "Should have responded with http status code %v but got %v", http.StatusOK, rr.Code)
}
