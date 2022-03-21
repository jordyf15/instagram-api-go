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

func TestPostHandlerSuite(t *testing.T) {
	suite.Run(t, new(PostHandlerSuite))
}

type PostHandlerSuite struct {
	suite.Suite
	postUsecase *mocks.PostUsecase
}

func (ph *PostHandlerSuite) SetupTest() {
	ph.postUsecase = new(mocks.PostUsecase)
}

func (ph *PostHandlerSuite) TestPostPostVisualMediaNotProvided() {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	writer.Close()
	postHandler := postHttp.NewPostHandler(ph.postUsecase)
	req, _ := http.NewRequest("POST", "/posts", bytes.NewReader(body.Bytes()))
	req.Header.Set("Content-Type", writer.FormDataContentType())
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(postHandler.Posts)
	handler.ServeHTTP(rr, req)

	assert.Equalf(ph.T(), http.StatusBadRequest, rr.Code, "Should have responded with http status code %v but got %v", http.StatusBadRequest, rr.Code)
	expectedBody := `{"message":"` + domain.ErrMissingVisualMediasInput.Error() + `"}`
	assert.Equalf(ph.T(), expectedBody, rr.Body.String(), "Should have responded with body %s but got %s", expectedBody, rr.Body.String())
}

func (ph *PostHandlerSuite) TestPostPostInsupportedVisualMediaTypeProvided() {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	fw, _ := writer.CreateFormFile("visual_medias", "bmp.bmp")
	file, _ := os.Open("./test_visual_medias/bmp.bmp")
	_, _ = io.Copy(fw, file)
	writer.Close()
	postHandler := postHttp.NewPostHandler(ph.postUsecase)
	req, _ := http.NewRequest("POST", "/posts", bytes.NewReader(body.Bytes()))
	req.Header.Set("Content-Type", writer.FormDataContentType())
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(postHandler.Posts)
	handler.ServeHTTP(rr, req)

	assert.Equalf(ph.T(), http.StatusBadRequest, rr.Code, "Should have responded with http status code %v but got %v", http.StatusBadRequest, rr.Code)
	expectedBody := `{"message":"` + domain.ErrUnsupportedVisualMediaType.Error() + `"}`
	assert.Equalf(ph.T(), expectedBody, rr.Body.String(), "Should have responded with body %s but got %s", expectedBody, rr.Body.String())
}

func (ph *PostHandlerSuite) TestPostPostCaptionNotProvided() {
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
	postHandler := postHttp.NewPostHandler(ph.postUsecase)
	req, _ := http.NewRequest("POST", "/posts", bytes.NewReader(body.Bytes()))
	req.Header.Set("Content-Type", writer.FormDataContentType())
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(postHandler.Posts)
	handler.ServeHTTP(rr, req)

	assert.Equalf(ph.T(), http.StatusBadRequest, rr.Code, "Should have responded with http status code %v but got %v", http.StatusBadRequest, rr.Code)
	expectedBody := `{"message":"` + domain.ErrMissingCaptionInput.Error() + `"}`
	assert.Equalf(ph.T(), expectedBody, rr.Body.String(), "Should have responded with body %s but got %s", expectedBody, rr.Body.String())
}

func (ph *PostHandlerSuite) TestPostPostInsertPostError() {
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
	ph.postUsecase.On("InsertPost", mock.AnythingOfType("*domain.Post"), mock.AnythingOfType("string"), mock.AnythingOfType("[]*multipart.FileHeader")).Return(domain.ErrInternalServerError)
	postHandler := postHttp.NewPostHandler(ph.postUsecase)
	req, _ := http.NewRequest("POST", "/posts", bytes.NewReader(body.Bytes()))
	req.Header.Set("Content-Type", writer.FormDataContentType())
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(postHandler.Posts)
	handler.ServeHTTP(rr, req)

	assert.Equalf(ph.T(), http.StatusInternalServerError, rr.Code, "Should have responded with http status code %v but got %v", http.StatusInternalServerError, rr.Code)
	expectedBody := `{"message":"` + domain.ErrInternalServerError.Error() + `"}`
	assert.Equalf(ph.T(), expectedBody, rr.Body.String(), "Should have responded with body %s but got %s", expectedBody, rr.Body.String())
}

func (ph *PostHandlerSuite) TestPostPostSuccessful() {
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
	ph.postUsecase.On("InsertPost", mock.AnythingOfType("*domain.Post"), mock.AnythingOfType("string"), mock.Anything).Return(nil)
	postHandler := postHttp.NewPostHandler(ph.postUsecase)
	req, _ := http.NewRequest("POST", "/posts", bytes.NewReader(body.Bytes()))
	req.Header.Set("Content-Type", writer.FormDataContentType())
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(postHandler.Posts)
	handler.ServeHTTP(rr, req)

	assert.Equalf(ph.T(), http.StatusCreated, rr.Code, "Should have responded with http status code %v but got %v", http.StatusCreated, rr.Code)
	expectedBody := `{"message":"Post successfully Created"}`
	assert.Equalf(ph.T(), expectedBody, rr.Body.String(), "Should have responded with body %s but got %s", expectedBody, rr.Body.String())
}

func (ph *PostHandlerSuite) TestGetPostsFindPostError() {
	req, _ := http.NewRequest("GET", "/posts", nil)
	rr := httptest.NewRecorder()
	ph.postUsecase.On("FindPosts").Return(nil, domain.ErrInternalServerError)
	postHandler := postHttp.NewPostHandler(ph.postUsecase)
	handler := http.HandlerFunc(postHandler.Posts)
	handler.ServeHTTP(rr, req)

	assert.Equalf(ph.T(), http.StatusInternalServerError, rr.Code, "Should have responded with http status code %v but got %v", http.StatusInternalServerError, rr.Code)
	expectedBody := `{"message":"` + domain.ErrInternalServerError.Error() + `"}`
	assert.Equalf(ph.T(), expectedBody, rr.Body.String(), "Should have responded with body %s but got %s", expectedBody, rr.Body.String())
}

func (ph *PostHandlerSuite) TestGetPostsSuccessful() {
	req, _ := http.NewRequest("GET", "/posts", nil)
	rr := httptest.NewRecorder()
	ph.postUsecase.On("FindPosts").Return(&[]domain.Post{}, nil)
	postHandler := postHttp.NewPostHandler(ph.postUsecase)
	handler := http.HandlerFunc(postHandler.Posts)
	handler.ServeHTTP(rr, req)

	assert.Equalf(ph.T(), http.StatusOK, rr.Code, "Should have responded with http status code %v but got %v", http.StatusOK, rr.Code)
}
func (ph *PostHandlerSuite) TestPutPostMissingCaption() {
	requestBody, _ := json.Marshal(map[string]string{
		"caption": "",
	})
	req, _ := http.NewRequest("PUT", "/posts/postid1", bytes.NewBuffer(requestBody))
	rr := httptest.NewRecorder()

	postHandler := postHttp.NewPostHandler(ph.postUsecase)
	handler := http.HandlerFunc(postHandler.Post)
	handler.ServeHTTP(rr, req)

	assert.Equalf(ph.T(), http.StatusBadRequest, rr.Code, "Should have responded with http status code %v but got %v", http.StatusBadRequest, rr.Code)
	expectedBody := `{"message":"` + domain.ErrMissingCaptionInput.Error() + `"}`
	assert.Equalf(ph.T(), expectedBody, rr.Body.String(), "Should have responded with body %s but got %s", expectedBody, rr.Body.String())
}

func (ph *PostHandlerSuite) TestPutPostUpdatePostError() {
	requestBody, _ := json.Marshal(map[string]string{
		"caption": "caption1",
	})
	req, _ := http.NewRequest("PUT", "/posts/postid1", bytes.NewBuffer(requestBody))
	rr := httptest.NewRecorder()
	ph.postUsecase.On("UpdatePost", mock.AnythingOfType("string"), mock.AnythingOfType("string"), mock.AnythingOfType("string")).Return(domain.ErrInternalServerError)
	postHandler := postHttp.NewPostHandler(ph.postUsecase)
	handler := http.HandlerFunc(postHandler.Post)
	handler.ServeHTTP(rr, req)

	assert.Equalf(ph.T(), http.StatusInternalServerError, rr.Code, "Should have responded with http status code %v but got %v", http.StatusInternalServerError, rr.Code)
	expectedBody := `{"message":"` + domain.ErrInternalServerError.Error() + `"}`
	assert.Equalf(ph.T(), expectedBody, rr.Body.String(), "Should have responded with body %s but got %s", expectedBody, rr.Body.String())
}

func (ph *PostHandlerSuite) TestPutPostSuccessful() {
	requestBody, _ := json.Marshal(map[string]string{
		"caption": "caption1",
	})
	req, _ := http.NewRequest("PUT", "/posts/postid1", bytes.NewBuffer(requestBody))
	rr := httptest.NewRecorder()
	ph.postUsecase.On("UpdatePost", mock.AnythingOfType("string"), mock.AnythingOfType("string"), mock.AnythingOfType("string")).Return(nil)
	postHandler := postHttp.NewPostHandler(ph.postUsecase)
	handler := http.HandlerFunc(postHandler.Post)
	handler.ServeHTTP(rr, req)

	assert.Equalf(ph.T(), http.StatusOK, rr.Code, "Should have responded with http status code %v but got %v", http.StatusOK, rr.Code)
	expectedBody := `{"message":"Post successfully Updated"}`
	assert.Equalf(ph.T(), expectedBody, rr.Body.String(), "Should have responded with body %s but got %s", expectedBody, rr.Body.String())
}

func (ph *PostHandlerSuite) TestDeletePostDeletePostError() {
	req, _ := http.NewRequest("DELETE", "/posts/postid1", nil)
	rr := httptest.NewRecorder()

	ph.postUsecase.On("DeletePost", mock.AnythingOfType("string"), mock.AnythingOfType("string")).Return(domain.ErrInternalServerError)
	postHandler := postHttp.NewPostHandler(ph.postUsecase)
	handler := http.HandlerFunc(postHandler.Post)
	handler.ServeHTTP(rr, req)

	assert.Equalf(ph.T(), http.StatusInternalServerError, rr.Code, "Should have responded with http status code %v but got %v", http.StatusInternalServerError, rr.Code)
	expectedBody := `{"message":"` + domain.ErrInternalServerError.Error() + `"}`
	assert.Equalf(ph.T(), expectedBody, rr.Body.String(), "Should have responded with body %s but got %s", expectedBody, rr.Body.String())
}

func (ph *PostHandlerSuite) TestDeletePostSuccessful() {
	req, _ := http.NewRequest("DELETE", "/posts/postid1", nil)
	rr := httptest.NewRecorder()

	ph.postUsecase.On("DeletePost", mock.AnythingOfType("string"), mock.AnythingOfType("string")).Return(nil)
	postHandler := postHttp.NewPostHandler(ph.postUsecase)
	handler := http.HandlerFunc(postHandler.Post)
	handler.ServeHTTP(rr, req)

	assert.Equalf(ph.T(), http.StatusOK, rr.Code, "Should have responded with http status code %v but got %v", http.StatusOK, rr.Code)
	expectedBody := `{"message":"Post successfully Deleted"}`
	assert.Equalf(ph.T(), expectedBody, rr.Body.String(), "Should have responded with body %s but got %s", expectedBody, rr.Body.String())
}
