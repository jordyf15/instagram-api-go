package http_test

import (
	"bytes"
	"encoding/json"
	"instagram-go/domain"
	"instagram-go/domain/mocks"
	postHttp "instagram-go/post/delivery/http"
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

func TestPostPostSuite(t *testing.T) {
	suite.Run(t, new(PostPostSuite))
}

func TestGetPostsSuite(t *testing.T) {
	suite.Run(t, new(GetPostsSuite))
}
func TestPutPostSuite(t *testing.T) {
	suite.Run(t, new(PutPostSuite))
}
func TestDeletePostSuite(t *testing.T) {
	suite.Run(t, new(DeletePostSuite))
}

type PostPostSuite struct {
	suite.Suite
	postUsecase *mocks.PostUsecase
}

func (pps *PostPostSuite) SetupTest() {
	pps.postUsecase = new(mocks.PostUsecase)
}

func (pps *PostPostSuite) TestVisualMediaNotProvided() {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	writer.Close()
	postHandler := postHttp.NewPostHandler(pps.postUsecase)
	req, _ := http.NewRequest("POST", "/posts", bytes.NewReader(body.Bytes()))
	req.Header.Set("Content-Type", writer.FormDataContentType())
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(postHandler.Posts)
	handler.ServeHTTP(rr, req)

	assert.Equalf(pps.T(), http.StatusBadRequest, rr.Code, "Should have responded with http status code %v but got %v", http.StatusBadRequest, rr.Code)
	expectedBody := `{"message":"` + domain.ErrMissingVisualMediasInput.Error() + `"}`
	assert.Equalf(pps.T(), expectedBody, rr.Body.String(), "Should have responded with body %s but got %s", expectedBody, rr.Body.String())
}

func (pps *PostPostSuite) TestInsupportedVisualMediaTypeProvided() {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	fw, _ := writer.CreateFormFile("visual_medias", "bmp.bmp")
	file, _ := os.Open("./test_visual_medias/bmp.bmp")
	_, _ = io.Copy(fw, file)
	writer.Close()
	postHandler := postHttp.NewPostHandler(pps.postUsecase)
	req, _ := http.NewRequest("POST", "/posts", bytes.NewReader(body.Bytes()))
	req.Header.Set("Content-Type", writer.FormDataContentType())
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(postHandler.Posts)
	handler.ServeHTTP(rr, req)

	assert.Equalf(pps.T(), http.StatusBadRequest, rr.Code, "Should have responded with http status code %v but got %v", http.StatusBadRequest, rr.Code)
	expectedBody := `{"message":"` + domain.ErrUnsupportedVisualMediaType.Error() + `"}`
	assert.Equalf(pps.T(), expectedBody, rr.Body.String(), "Should have responded with body %s but got %s", expectedBody, rr.Body.String())
}

func (pps *PostPostSuite) TestCaptionNotProvided() {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	fw, _ := writer.CreateFormFile("visual_medias", "gif.gif")
	file, _ := os.Open("./test_visual_medias/gif.gif")
	_, _ = io.Copy(fw, file)
	fw, _ = writer.CreateFormFile("visual_medias", "jpg.jpg")
	file, _ = os.Open("./test_visual_medias/jpg.jpg")
	_, _ = io.Copy(fw, file)
	fw, _ = writer.CreateFormFile("visual_medias", "png.png")
	file, _ = os.Open("./test_visual_medias/png.png")
	_, _ = io.Copy(fw, file)
	fw, _ = writer.CreateFormFile("visual_medias", "tiff.tiff")
	file, _ = os.Open("./test_visual_medias/tiff.tiff")
	_, _ = io.Copy(fw, file)
	writer.Close()
	postHandler := postHttp.NewPostHandler(pps.postUsecase)
	req, _ := http.NewRequest("POST", "/posts", bytes.NewReader(body.Bytes()))
	req.Header.Set("Content-Type", writer.FormDataContentType())
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(postHandler.Posts)
	handler.ServeHTTP(rr, req)

	assert.Equalf(pps.T(), http.StatusBadRequest, rr.Code, "Should have responded with http status code %v but got %v", http.StatusBadRequest, rr.Code)
	expectedBody := `{"message":"` + domain.ErrMissingCaptionInput.Error() + `"}`
	assert.Equalf(pps.T(), expectedBody, rr.Body.String(), "Should have responded with body %s but got %s", expectedBody, rr.Body.String())
}

func (pps *PostPostSuite) TestInsertPostError() {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	fw, _ := writer.CreateFormFile("visual_medias", "gif.gif")
	file, _ := os.Open("./test_visual_medias/gif.gif")
	_, _ = io.Copy(fw, file)
	fw, _ = writer.CreateFormFile("visual_medias", "jpg.jpg")
	file, _ = os.Open("./test_visual_medias/jpg.jpg")
	_, _ = io.Copy(fw, file)
	fw, _ = writer.CreateFormFile("visual_medias", "png.png")
	file, _ = os.Open("./test_visual_medias/png.png")
	_, _ = io.Copy(fw, file)
	fw, _ = writer.CreateFormFile("visual_medias", "tiff.tiff")
	file, _ = os.Open("./test_visual_medias/tiff.tiff")
	_, _ = io.Copy(fw, file)
	fw, _ = writer.CreateFormField("caption")
	_, _ = io.Copy(fw, strings.NewReader("a new caption"))
	writer.Close()
	pps.postUsecase.On("InsertPost", mock.Anything, mock.Anything, mock.Anything).Return(domain.ErrInternalServerError)
	postHandler := postHttp.NewPostHandler(pps.postUsecase)
	req, _ := http.NewRequest("POST", "/posts", bytes.NewReader(body.Bytes()))
	req.Header.Set("Content-Type", writer.FormDataContentType())
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(postHandler.Posts)
	handler.ServeHTTP(rr, req)

	assert.Equalf(pps.T(), http.StatusInternalServerError, rr.Code, "Should have responded with http status code %v but got %v", http.StatusInternalServerError, rr.Code)
	expectedBody := `{"message":"` + domain.ErrInternalServerError.Error() + `"}`
	assert.Equalf(pps.T(), expectedBody, rr.Body.String(), "Should have responded with body %s but got %s", expectedBody, rr.Body.String())
}

func (pps *PostPostSuite) TestPostPostSuccessful() {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	fw, _ := writer.CreateFormFile("visual_medias", "gif.gif")
	file, _ := os.Open("./test_visual_medias/gif.gif")
	_, _ = io.Copy(fw, file)
	fw, _ = writer.CreateFormFile("visual_medias", "jpg.jpg")
	file, _ = os.Open("./test_visual_medias/jpg.jpg")
	_, _ = io.Copy(fw, file)
	fw, _ = writer.CreateFormFile("visual_medias", "png.png")
	file, _ = os.Open("./test_visual_medias/png.png")
	_, _ = io.Copy(fw, file)
	fw, _ = writer.CreateFormFile("visual_medias", "tiff.tiff")
	file, _ = os.Open("./test_visual_medias/tiff.tiff")
	_, _ = io.Copy(fw, file)
	fw, _ = writer.CreateFormField("caption")
	_, _ = io.Copy(fw, strings.NewReader("a new caption"))
	writer.Close()
	pps.postUsecase.On("InsertPost", mock.Anything, mock.Anything, mock.Anything).Return(nil)
	postHandler := postHttp.NewPostHandler(pps.postUsecase)
	req, _ := http.NewRequest("POST", "/posts", bytes.NewReader(body.Bytes()))
	req.Header.Set("Content-Type", writer.FormDataContentType())
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(postHandler.Posts)
	handler.ServeHTTP(rr, req)

	assert.Equalf(pps.T(), http.StatusCreated, rr.Code, "Should have responded with http status code %v but got %v", http.StatusCreated, rr.Code)
	expectedBody := `{"message":"Post successfully Created"}`
	assert.Equalf(pps.T(), expectedBody, rr.Body.String(), "Should have responded with body %s but got %s", expectedBody, rr.Body.String())
}

type GetPostsSuite struct {
	suite.Suite
	postUsecase *mocks.PostUsecase
}

func (gps *GetPostsSuite) SetupTest() {
	gps.postUsecase = new(mocks.PostUsecase)
}

func (gps *GetPostsSuite) TestFindPostError() {
	req, _ := http.NewRequest("GET", "/posts", nil)
	rr := httptest.NewRecorder()
	gps.postUsecase.On("FindPosts").Return(nil, domain.ErrInternalServerError)
	postHandler := postHttp.NewPostHandler(gps.postUsecase)
	handler := http.HandlerFunc(postHandler.Posts)
	handler.ServeHTTP(rr, req)

	assert.Equalf(gps.T(), http.StatusInternalServerError, rr.Code, "Should have responded with http status code %v but got %v", http.StatusInternalServerError, rr.Code)
	expectedBody := `{"message":"` + domain.ErrInternalServerError.Error() + `"}`
	assert.Equalf(gps.T(), expectedBody, rr.Body.String(), "Should have responded with body %s but got %s", expectedBody, rr.Body.String())
}

func (gps *GetPostsSuite) TestGetPostsSuccessful() {
	req, _ := http.NewRequest("GET", "/posts", nil)
	rr := httptest.NewRecorder()
	gps.postUsecase.On("FindPosts").Return(&[]domain.Post{}, nil)
	postHandler := postHttp.NewPostHandler(gps.postUsecase)
	handler := http.HandlerFunc(postHandler.Posts)
	handler.ServeHTTP(rr, req)

	assert.Equalf(gps.T(), http.StatusOK, rr.Code, "Should have responded with http status code %v but got %v", http.StatusOK, rr.Code)
}

type PutPostSuite struct {
	suite.Suite
	postUsecase *mocks.PostUsecase
}

func (pps *PutPostSuite) SetupTest() {
	pps.postUsecase = new(mocks.PostUsecase)
}

func (pps *PutPostSuite) TestMissingCaption() {
	requestBody, _ := json.Marshal(map[string]string{
		"caption": "",
	})
	req, _ := http.NewRequest("PUT", "/posts/postid1", bytes.NewBuffer(requestBody))
	rr := httptest.NewRecorder()

	postHandler := postHttp.NewPostHandler(pps.postUsecase)
	handler := http.HandlerFunc(postHandler.Post)
	handler.ServeHTTP(rr, req)

	assert.Equalf(pps.T(), http.StatusBadRequest, rr.Code, "Should have responded with http status code %v but got %v", http.StatusBadRequest, rr.Code)
	expectedBody := `{"message":"` + domain.ErrMissingCaptionInput.Error() + `"}`
	assert.Equalf(pps.T(), expectedBody, rr.Body.String(), "Should have responded with body %s but got %s", expectedBody, rr.Body.String())
}

func (pps *PutPostSuite) TestUpdatePostError() {
	requestBody, _ := json.Marshal(map[string]string{
		"caption": "caption1",
	})
	req, _ := http.NewRequest("PUT", "/posts/postid1", bytes.NewBuffer(requestBody))
	rr := httptest.NewRecorder()
	pps.postUsecase.On("UpdatePost", mock.Anything, mock.Anything, mock.Anything).Return(domain.ErrInternalServerError)
	postHandler := postHttp.NewPostHandler(pps.postUsecase)
	handler := http.HandlerFunc(postHandler.Post)
	handler.ServeHTTP(rr, req)

	assert.Equalf(pps.T(), http.StatusInternalServerError, rr.Code, "Should have responded with http status code %v but got %v", http.StatusInternalServerError, rr.Code)
	expectedBody := `{"message":"` + domain.ErrInternalServerError.Error() + `"}`
	assert.Equalf(pps.T(), expectedBody, rr.Body.String(), "Should have responded with body %s but got %s", expectedBody, rr.Body.String())
}

func (pps *PutPostSuite) TestPutPostSuccessful() {
	requestBody, _ := json.Marshal(map[string]string{
		"caption": "caption1",
	})
	req, _ := http.NewRequest("PUT", "/posts/postid1", bytes.NewBuffer(requestBody))
	rr := httptest.NewRecorder()
	pps.postUsecase.On("UpdatePost", mock.Anything, mock.Anything, mock.Anything).Return(nil)
	postHandler := postHttp.NewPostHandler(pps.postUsecase)
	handler := http.HandlerFunc(postHandler.Post)
	handler.ServeHTTP(rr, req)

	assert.Equalf(pps.T(), http.StatusOK, rr.Code, "Should have responded with http status code %v but got %v", http.StatusOK, rr.Code)
	expectedBody := `{"message":"Post successfully Updated"}`
	assert.Equalf(pps.T(), expectedBody, rr.Body.String(), "Should have responded with body %s but got %s", expectedBody, rr.Body.String())
}

type DeletePostSuite struct {
	suite.Suite
	postUsecase *mocks.PostUsecase
}

func (dps *DeletePostSuite) SetupTest() {
	dps.postUsecase = new(mocks.PostUsecase)
}

func (dps *DeletePostSuite) TestDeletePostError() {
	req, _ := http.NewRequest("DELETE", "/posts/postid1", nil)
	rr := httptest.NewRecorder()

	dps.postUsecase.On("DeletePost", mock.Anything, mock.Anything).Return(domain.ErrInternalServerError)
	postHandler := postHttp.NewPostHandler(dps.postUsecase)
	handler := http.HandlerFunc(postHandler.Post)
	handler.ServeHTTP(rr, req)

	assert.Equalf(dps.T(), http.StatusInternalServerError, rr.Code, "Should have responded with http status code %v but got %v", http.StatusInternalServerError, rr.Code)
	expectedBody := `{"message":"` + domain.ErrInternalServerError.Error() + `"}`
	assert.Equalf(dps.T(), expectedBody, rr.Body.String(), "Should have responded with body %s but got %s", expectedBody, rr.Body.String())
}

func (dps *DeletePostSuite) TestDeletePostSuccessful() {
	req, _ := http.NewRequest("DELETE", "/posts/postid1", nil)
	rr := httptest.NewRecorder()

	dps.postUsecase.On("DeletePost", mock.Anything, mock.Anything).Return(nil)
	postHandler := postHttp.NewPostHandler(dps.postUsecase)
	handler := http.HandlerFunc(postHandler.Post)
	handler.ServeHTTP(rr, req)

	assert.Equalf(dps.T(), http.StatusOK, rr.Code, "Should have responded with http status code %v but got %v", http.StatusOK, rr.Code)
	expectedBody := `{"message":"Post successfully Deleted"}`
	assert.Equalf(dps.T(), expectedBody, rr.Body.String(), "Should have responded with body %s but got %s", expectedBody, rr.Body.String())
}
